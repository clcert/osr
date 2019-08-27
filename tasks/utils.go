package tasks

import (
	"github.com/spf13/viper"
	"path/filepath"
)

// Helper function which returns the Task locations using the config values.
func GetTasksPath() (string, error) {
	home := viper.GetString("folders.home")
	tasks := viper.GetString("folders.tasks")
	return filepath.Join(home, tasks), nil
}
