// This package contains utilities used in the configuration
// and interaction with OSR Scanner remote servers.
package remote

import (
	"bytes"
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"net"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

// Err when you don't execute connect before using a server.
var ServerNotInitialized = fmt.Errorf("server not initialized. use Connect() before using this command")

// Server defines one remote server.
// This server can execute one command at a time
// And has not direct access to its Standard Input
// (The idea is to execute unattended commands)
type Server struct {
	Name       string          // name of the remote server
	Address    string          // Address of the remote server
	Username   string          // Username of the remote server
	sshClient  *ssh.Client     // SSH connection
	sftpClient *sftp.Client    // A lazy loaded sshClient for SFTP file transmission
	Sessions   []*ssh.Session  // Opened SSH Sessions
	Output     *logs.OSROutput // Special Output with stdout + stderr of the server.
}

func (s *Server) GetAttachments() (attachments []string) {
	if s.Output != nil {
		attachments = []string{s.Output.Path}
	}
	return
}

func (s *Server) GetAddr() (net.Addr, error) {
	if s.sshClient == nil {
		return nil, ServerNotInitialized
	}
	return s.sshClient.Conn.RemoteAddr(), nil
}

// Connect connects to a server using Public Keys or password in config file
// and returns nil if the connection was successful.
func (s *Server) Connect() error {
	if s.sshClient != nil {
		// Our work here is done. Or should we return an error?
		return nil
	}
	privKeySigner, err := s.getPrivKeySigner()
	if err != nil {
		return err
	}
	hostKey, err := s.getHostKey()
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privKeySigner),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	logs.Log.WithFields(logrus.Fields{
		"server_name": s.Name,
		"address":     s.Address,
		"username":    s.Username,
	}).Info("Connecting to server...")
	connectServer := fmt.Sprintf("%s:%d", s.Address, DefaultSSHPort)
	connection, err := ssh.Dial("tcp", connectServer, config)
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Connection failed")
		return err
	}
	logs.Log.WithFields(logrus.Fields{
		"server_name": s.Name,
		"address":     s.Address,
		"username":    s.Username,
	}).Info("Connected!")
	newOutput, err := logs.NewOutput("remote", s.Name)
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("could not create output")
		_ = connection.Close()
		return err
	}
	s.Output = newOutput
	s.Sessions = make([]*ssh.Session, 0)
	s.sshClient = connection
	return nil
}

func (s *Server) GetSSHClient() (*ssh.Client, error) {
	if s.sshClient != nil {
		return s.sshClient, nil
	}
	return nil, ServerNotInitialized
}

func (s *Server) GetSFTPClient() (*sftp.Client, error) {
	if s.sshClient == nil {
		return nil, ServerNotInitialized
	} else if s.sftpClient == nil {
		sftpClient, err := sftp.NewClient(s.sshClient)
		if err != nil {
			return nil, err
		}
		s.sftpClient = sftpClient
	}
	return s.sftpClient, nil
}

func (s *Server) NewSession() (*ssh.Session, error) {
	sess, err := s.sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	s.Sessions = append(s.Sessions, sess)
	return sess, nil
}

// Connect to a server using a interactively-placed password.
// If it fails, it returns an error.
// TODO: make it non-interactive (password in config file, but deleting it after the key exchange.
func (s *Server) ConnectWithPassword() error {
	if !terminal.IsTerminal(int(syscall.Stdin)) {
		return fmt.Errorf("cannot init in a non interactive shell: interactive password input needed")
	}
	fmt.Printf("Insert Password from %s (%s): ", s.Name, s.Address)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(bytePassword)),
		},
		HostKeyCallback: s.SaveHostKey(),
	}

	connectServer := fmt.Sprintf("%s:%d", s.Address, DefaultSSHPort)
	client, err := ssh.Dial("tcp", connectServer, config)
	if err != nil {
		return ConnectionError{
			Name:       s.Name,
			Address:    s.Address,
			Username:   s.Username,
			AuthMethod: "password",
			Err:        err,
		}
	}
	newOutput, err := logs.NewOutput("remote", s.Name)
	if err != nil {
		return err
	}
	s.Output = newOutput
	s.Sessions = make([]*ssh.Session, 0)
	s.sshClient = client
	return nil
}

// Get health information associated to a server. This method creates
// a connection for obtaining this values.
func (s *Server) GetInfo() (*ServerInfo, error) {
	if s.sshClient == nil {
		return nil, ServerNotInitialized
	}
	session, err := s.NewSession()
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run("POSIXLY_CORRECT=1 df") // This returns 512-blocks and works equally in all systems (linux+openbsd)
	if err != nil {
		session.Close()
		return nil, err
	}
	disks, err := ParseDf(b)
	if err != nil {
		session.Close()
		return nil, err
	}
	return &ServerInfo{
		Name:    s.Name,
		Date:    time.Now(),
		Address: s.Address,
		Disks:   disks,
	}, err
}

// Exec executes a remote command non interactively.
// The output is appended in the Output file related to the server.
// Returns only when the command has been executed
// Returns an error if something fails.
func (s *Server) RemoteExecute(command string) error {
	if s.sshClient == nil {
		return ServerNotInitialized
	}
	sess, err := s.NewSession()
	if err != nil {
		return err
	}
	// Set out
	sess.Stdout = s.Output.Writer
	sess.Stderr = s.Output.Writer
	err = sess.Run(command)
	if err != nil {
		return err
	}
	return nil
}

// Exec executes a local command non interactively.
// That means, the local file is uploaded to "destination" directory on the server, and then executed.
// The output is appended in the Output file related to the server.
// If the local file is not executable, this will fail.
// The local file must be located on 'scripts' config folder.
// The params struct, if not nil, allows to format the script using the parameters defined on it.
// Returns an error if something fails.
func (s *Server) ExecuteLocalScript(destination, source string, params utils.Params) error {
	// TODO: sanitize source so you couldn't update files outside the scripts path?
	if s.sshClient == nil {
		return ServerNotInitialized
	}
	scriptRoot, err := GetScriptsPath()
	if err != nil {
		return err
	}
	paths, err := s.SendFiles(params, destination, path.Join(scriptRoot, source))
	if err != nil {
		return err
	}
	// I asked for only one path.
	path := paths[0]
	return s.RemoteExecute(path)
}

// Sends an arbitrary number of local files (with their absolute paths)
// to a folder into the remote server. It also copies the local File Mode
// to the remote File Mode.
// The first argument is a Params struct, that is used for file content formatting.
// If it is not nil, the system will try to format the contents of the files with the params.
func (s *Server) SendFiles(params utils.Params, dest string, paths ...string) ([]string, error) {
	if s.sshClient == nil {
		return nil, ServerNotInitialized
	}
	sftpClient, err := s.GetSFTPClient()
	if err != nil {
		return nil, err
	}
	var remotePaths []string
	for _, aPath := range paths {
		fLocal, err := os.Open(aPath)
		if err != nil {
			return nil, err
		}
		if err := MkDirParents(sftpClient, dest); err != nil {
			fLocal.Close()
			return nil, err
		}
		splitPath := strings.Split(aPath, "/")
		filename := splitPath[len(splitPath)-1]
		remotePath := sftp.Join(dest, filename)
		fRemote, err := sftpClient.Create(remotePath)
		if err != nil {
			fLocal.Close()
			return nil, err
		}
		var fReader io.Reader
		if params != nil {
			fReader = params.FormatReader(fLocal)
		} else {
			fReader = fLocal
		}
		byteNum, err := fRemote.ReadFrom(fReader)
		logs.Log.WithFields(logrus.Fields{
			"filename": filename,
			"from":     aPath,
			"to":       dest,
		}).Infof("%d bytes copied successfully", byteNum)
		stat, err := fLocal.Stat()
		if err != nil {
			fLocal.Close()
			fRemote.Close()
			return nil, err
		}
		// Use the same permissions of source
		err = s.sftpClient.Chmod(remotePath, stat.Mode())
		if err != nil {
			fLocal.Close()
			fRemote.Close()
			return nil, err
		}
		remotePaths = append(remotePaths, remotePath)
		fLocal.Close()
		fRemote.Close()
	}
	return remotePaths, nil
}

// Closes the sftp connection and all the other connections which colud have been kept open.
// The object is reusable! (you can connect to it again)
func (s *Server) Close() (err error) {
	if s.sshClient == nil {
		return ServerNotInitialized
	}
	if err = s.sshClient.Close(); err != nil {
		return
	}
	s.sshClient = nil
	if s.sftpClient != nil {
		if err = s.sftpClient.Close(); err != nil {
			return
		}
	}
	s.sftpClient = nil
	for _, session := range s.Sessions {
		session.Close() // The error could be EOF (already closed), but we ignore that because we are closing the channels just in case.
	}
	s.Sessions = make([]*ssh.Session, 0)
	return
}
