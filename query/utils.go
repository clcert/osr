package query

import (
	"github.com/spf13/viper"
	"path/filepath"
)

// Helper function which returns the queries locations using the config values.
func GetQueriesPath() (string, error) {
	home := viper.GetString("folders.home")
	logs := viper.GetString("folders.queries")
	return filepath.Join(home, logs), nil
}
