// Mailer defines the routines to send notification mails related to OSR operation.
package mailer

import (
	"fmt"
	"github.com/spf13/viper"
)

// NotifyLevel defines the different notification levels.
type NotifyLevel int

const (
	DEBUG NotifyLevel = iota
	INFO
	WARN
	ERROR
	PANIC
	NONE
)

// Config defines a configuration related to mail sending.
type Config struct {
	Name     string // Name used as from
	Email    string // Email used as from
	Username string // Username used in the mail sender
	Password string // Password used in the mail sender
	Server   string // Mail server name, as defined in osr
	Port     uint16 // Mail server port.
}

// Returns the connection string to the mail server.
func (c Config) GetConnString() string {
	return fmt.Sprintf("%s:%d", c.Server, c.Port)
}

// Returns the address well formatted to be used in mail sending.
func (c Config) GetFormattedAddress() string {
	return fmt.Sprintf("%s <%s>", c.Name, c.Email)
}

// Returns a list with the servers defined in config file.
func GetNotifyCredentials() (*Config, error) {
	var credentials Config
	err := viper.UnmarshalKey("mailer.credentials", &credentials)
	return &credentials, err
}

// Returns a list with all the mails registered as notificable.
func GetNotifyEmails() ([]string, error) {
	emails := viper.GetStringSlice("mailer.emails")
	return emails, nil
}

// Returns the current notify level for the system.
func GetNotifyLevel() (level NotifyLevel, err error) {
	level = NotifyLevel(viper.GetInt("mailer.notifylevel"))
	return
}
