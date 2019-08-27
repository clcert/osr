package remote

import (
	"bytes"
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Creates the keys used in PublicKey authentication for all the servers
// which doesn't have them.
func CreateKeysForRemote() error {
	config, err := GetRemoteConfig()
	if err != nil {
		return err
	}
	for _, server := range config {
		keyExists, err := server.HasKeys()
		if err != nil {
			server.DeleteKeys()
			return err
		}
		if !keyExists {
			logs.Log.WithFields(logrus.Fields{
				"server":   server.Name,
				"username": server.Username,
				"address":  server.Address,
			}).Info("Key for server doesn't exist. Creating pair...")
			err := server.CreateServerKey()
			if err != nil {
				server.DeleteKeys()
				return err
			}
			tries := 1
			for {
				err = server.ConnectWithPassword()
				if err != nil {
					if err, ok := err.(ConnectionError); ok && tries < 3 {
						tries++
						logs.Log.WithFields(logrus.Fields{
							"tries":    tries,
							"name":     err.Name,
							"address":  err.Address,
							"username": err.Username,
						}).Error("cannot log in to server: invalid credentials. Trying again...")
						continue
					}
					_ = server.DeleteKeys()
					return err
				} else {
					break
				}
			}

			session, err := server.NewSession()
			if err != nil {
				_ = server.DeleteKeys()
				return err
			}
			pubKeyPath, err := server.GetPubKeyPath()
			if err != nil {
				session.Close()
				_ = server.DeleteKeys()
				return err
			}
			pubKeyReader, err := os.Open(pubKeyPath)
			if err != nil {
				session.Close()
				_ = server.DeleteKeys()
				return err
			}
			session.Stdin = pubKeyReader
			// one liner for installing a new authorized key
			if err := session.Run("umask 077 && mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"); err != nil {
				session.Close()
				_ = server.DeleteKeys()
				return err
			}
			logs.Log.WithFields(logrus.Fields{
				"server_name": server.Name,
				"address":     server.Address,
				"username":    server.Username,
			}).Info("Public Key Added!")

		} else {
			logs.Log.WithFields(logrus.Fields{
				"server_name": server.Name,
				"address":     server.Address,
				"username":    server.Username,
			}).Info("Key for server already exists!")
		}

	}
	return nil
}

// Executes getInfo in all the servers given as argument.
func ListServersInfo(servers []*Server) error {
	for _, server := range servers {
		err := server.Connect()
		if err != nil {
			return err
		}
		info, err := server.GetInfo()
		if err != nil {
			_ = server.Close()
			return err
		}
		notify(server, info)
		server.Output.Println(info.String())
		server.Close()
	}
	return nil
}

// Creates all the non existent folders.
// Based on MkDirParents implementation in Golang.
func MkDirParents(client *sftp.Client, dir string) (err error) {
	var parents string

	if path.IsAbs(dir) {
		// Otherwise, an absolute path given below would be turned in to a relative one
		// by splitting on "/"
		parents = "/"
	}

	for _, name := range strings.Split(dir, "/") {
		if name == "" {
			// Paths with double-/ in them should just move along
			// this will also catch the case of the first character being a "/", i.e. an absolute path
			continue
		}
		parents = path.Join(parents, name)
		err = client.Mkdir(parents)
		if status, ok := err.(*sftp.StatusError); ok {
			if status.Code == uint32(sftp.ErrSshFxFailure) {
				var fi os.FileInfo
				fi, err = client.Stat(parents)
				if err == nil {
					if !fi.IsDir() {
						return fmt.Errorf("file exists: %s", parents)
					}
				}
			}
		}
		if err != nil {
			break
		}
	}
	return err
}

// Parses the output of ``df'' commmand and stores it in a array of Disk structs.
func ParseDf(b bytes.Buffer) (disks []*Disk, err error) {
	lines := 0
	var line string
	disks = make([]*Disk, 0)
	for {
		lines++
		line, err = b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		if lines == 1 {
			continue
		}
		values := strings.Fields(line)
		var blocks, used, available, capacity int
		blocks, err = strconv.Atoi(values[1])
		if err != nil {
			return
		}
		used, err = strconv.Atoi(values[2])
		if err != nil {
			return
		}
		available, err = strconv.Atoi(values[3])
		if err != nil {
			return
		}
		capacity, err = strconv.Atoi(values[4][:len(values[4])-1]) // Removing the % symbol
		if err != nil {
			return
		}
		disk := &Disk{
			Name:       values[0],
			Blocks:     blocks,
			Used:       used,
			Available:  available,
			Capacity:   capacity,
			Mountpoint: values[5],
		}
		disks = append(disks, disk)
	}
	return
}

// Creates a default Process Folder, using the task number.
func CreateDefaultProcessFolder(client *sftp.Client, name string) (folder string, err error) {
	// Default folder is ~/osr/YYYY-MM-DD/<name>/ (assuming ~ as the SFTP initial directory)
	formattedDate := time.Now().Format("2006-01-02")
	folder = fmt.Sprintf("osr/%s/%s/", formattedDate, name)
	err = MkDirParents(client, folder)
	return
}

// Returns the path of the public key installed on the remote server.
// This key could not exist. So it's important to check for its existence before.
func GetScriptsPath() (string, error) {
	home := viper.GetString("folders.home")
	scripts := viper.GetString("folders.scripts")
	return filepath.Join(home, scripts), nil
}
