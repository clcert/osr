
// logs groups the methods used for logging in the application.
package logs

import (
	"github.com/spf13/viper"
	"path/filepath"
)

// This is the default logger of the app, used by the logs modules.
var Log *OSRLog

// InitDefaultLog initializes the default logger for the application.
func InitDefaultLog() (err error) {
	Log, err = NewLog("osr")
	return
}

// getLogsPath is a helper function which returns the log locations using the config values.
func getLogsPath() (string, error) {
	home := viper.GetString("folders.home")
	logs := viper.GetString("folders.logs")
	return filepath.Join(home, logs), nil
}
