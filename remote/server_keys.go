package remote

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/clcert/osr/logs"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

// Server, with password field. Used only in key creation routines
type ServerWithPassword struct {
	Server
	Password string
}

// Returns the path of the public key installed on the remote server.
// This key could not exist. So it's important to check for its existence before.
func (s *Server) GetPubKeyPath() (string, error) {
	keysPath, err := getKeysPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(keysPath, s.Name+PublicKeyExtension), nil
}

// Returns the path of the host key of the remote server.
// This key could not exist, so it's important to check for its existence before.
func (s *Server) GetHostKeyPath() (string, error) {
	keysPath, err := getKeysPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(keysPath, s.Name+HostPublicKeyExtension), nil
}

// Returns the path of the public key used for remote server authentication.
// This key could not exist, so it's important to check for its existence before.
func (s *Server) GetPrivKeyPath() (string, error) {
	keysPath, err := getKeysPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(keysPath, s.Name+PrivateKeyExtension), nil
}

// This function checks if the server created its keys before.
// It checks the private key and the host public key.
// If any of them doesn't exist, it returns false.
func (s *Server) HasKeys() (bool, error) {
	privKeyName, err := s.GetPrivKeyPath()
	if err != nil {
		return false, err
	}
	hostKeyName, err := s.GetHostKeyPath()
	if _, err := os.Stat(privKeyName); os.IsNotExist(err) {
		return false, nil
	}
	if _, err := os.Stat(hostKeyName); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

// Returns a callback used in the first SSH connection.
// This callback saves the host key as "trusted", allowing to use it
// in future connections with public key authentication.
func (s *Server) SaveHostKey() func(hostname string, remote net.Addr, key ssh.PublicKey) error {
	hostPath, err := s.GetHostKeyPath()
	if err != nil {
		return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return err
		}
	}
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return ioutil.WriteFile(hostPath, ssh.MarshalAuthorizedKey(key), 0655)
	}
}

// Returns the Private Key in signer format, allowing to use it directly
// in an SSH connection.
func (s *Server) getPrivKeySigner() (ssh.Signer, error) {
	privKeyPath, err := s.GetPrivKeyPath()
	if err != nil {
		return nil, err
	}
	privKeyFile, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}
	privKeySigner, err := ssh.ParsePrivateKey(privKeyFile)
	if err != nil {
		return nil, err
	}
	return privKeySigner, nil
}

// Returns the host key in a public key format, allowing to use it
// directly in a SSH connection.
func (s *Server) getHostKey() (ssh.PublicKey, error) {
	hostKeyPath, err := s.GetHostKeyPath()
	if err != nil {
		return nil, err
	}
	hostKeyFile, err := ioutil.ReadFile(hostKeyPath)
	if err != nil {
		return nil, err
	}
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey(hostKeyFile)
	if err != nil {
		return nil, err
	}
	return hostKey, nil
}

// Based on https://stackoverflow.com/questions/21151714/go-generate-an-ssh-public-key
// CreateServerKey creates a pair of public and private keys for SSH access to a server.
// Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// Private Key generated is PEM encoded and both files are saved on keys folder.
func (s *Server) CreateServerKey() error {
	privateKeyPath, err := s.GetPrivKeyPath()
	if err != nil {
		return err
	}
	pubKeyPath, err := s.GetPubKeyPath()
	if err != nil {
		return err
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	privateKeyFile, err := os.OpenFile(privateKeyPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
}

// Tries to delete the public, private and host keys of a server.
func (s *Server) DeleteKeys() error {
	pubKeyPath, err := s.GetPubKeyPath()
	if err != nil {
		return err
	}
	if err := os.Remove(pubKeyPath); err != nil {
		return err
	}
	privKeyPath, err := s.GetPrivKeyPath()
	if err != nil {
		return err
	}
	if err := os.Remove(privKeyPath); err != nil {
		return err
	}
	hostKeyPath, err := s.GetHostKeyPath()
	if err != nil {
		return err
	}
	if err := os.Remove(hostKeyPath); err != nil {
		return err
	}
	logs.Log.WithFields(logrus.Fields{
		"name":     s.Name,
		"address":  s.Address,
		"username": s.Username,
	}).Error("Public, private and host keys from server deleted!")

	return nil
}
