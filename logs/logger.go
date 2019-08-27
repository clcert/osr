package logs

import (
	"github.com/spf13/viper"
	"path/filepath"
)

// This is the default logger of the app, used by the logs modules.
// In a near future, there should be a way to create new loggers in some situations,
// (like a logger special for a scanner device).
var Log *OSRLog

// Inits the default logger for the application.
func InitDefaultLog() (err error) {
	Log, err = NewLog("osr")
	return
}

// Helper function which returns the log locations using the config values.
func getLogsPath() (string, error) {
	home := viper.GetString("folders.home")
	logs := viper.GetString("folders.logs")
	return filepath.Join(home, logs), nil
}
