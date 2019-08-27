// OSR is a tool for scanning and processing data from internet scans.
package main

import (
	"fmt"
	"github.com/clcert/osr/cmd"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/panics"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/osr/")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".osr"))
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(&panics.Info{
			Text: "fatal error reading config file",
			Err:  err,
		})
	}
	if err := CreateOSRDirs(); err != nil {
		panic(&panics.Info{
			Text: "fatal error reading folders.home",
			Err:  err,
		})
	}
	if err := logs.InitDefaultLog(); err != nil {
		panic(&panics.Info{
			Text: "fatal error creating logs",
			Err:  err,
		})
	}
}

func CreateOSRDirs() error {
	viper.SetDefault("folders.home", filepath.Join(os.Getenv("HOME"), ".osr"))
	home := viper.GetString("folders.home")
	folders := []string{"scripts", "queries", "keys", "logs", "tasks"}
	for _, folder := range folders {
		viper.SetDefault(fmt.Sprintf("folders.%s", folder), folder)
		key := fmt.Sprintf("folders.%s", folder)
		val := viper.GetString(key)
		if err := os.MkdirAll(path.Join(home, val), 0744); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	cmd.Execute()
}
