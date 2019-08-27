package remote

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

// Returns a list with the servers defined in config file.
func GetRemoteConfig() ([]*Server, error) {
	var serverList []*Server
	err := viper.UnmarshalKey("remote", &serverList)
	return serverList, err
}

// Returns the combined path for the remote server keys.
func getKeysPath() (string, error) {
	home := viper.GetString("folders.home")
	keys := viper.GetString("folders.keys")
	return filepath.Join(home, keys), nil
}

// Returns an specific server configuration, based on the name provided.
func GetServer(serverName string) (*Server, error) {
	servers, err := GetServers(serverName)
	if err != nil {
		return nil, err
	}
	return servers[0], nil
}

// Returns an specific servers configuration, based on the names provided.
// If serverNames is empty, returns the information of all the servers.
func GetServers(serverNames ...string) ([]*Server, error) {
	var servers []*Server
	config, err := GetRemoteConfig()
	if err != nil {
		return servers, err
	}
	if len(serverNames) == 0 {
		return config, nil
	}
	for _, name := range serverNames {
		for _, server := range config {
			if server.Name == name {
				servers = append(servers, server)
				break
			}
		}
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("server not found in config file")
	}
	return servers, nil
}
