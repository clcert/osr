// This package manages the databases creation and usage.
package databases

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// Config struct represents each item in the "databases" configuration in
// the config file.
type Config struct {
	Server string      // Host of the database server
	Port   int         // Port of the database server
	DBName string      // name of the database
	Writer Credentials // Writer credentials for the database
	Reader Credentials // Reader credentials for the database
}

// Credentials represents an username and a password for a database
type Credentials struct {
	Username string // Username of the database user
	Password string // Password of the database user
}

// Unmarshals the config defined in viper as a Config struct.
func GetDBConfig(name string) (Config, error) {
	var dbConf map[string]Config
	err := viper.UnmarshalKey("databases", &dbConf)
	return dbConf[name], err
}

// Rewrites the credentials in "databases" config property, in the given
// subproperty.
func WriteDBConf(db string, conf Config) error {
	viper.Set(fmt.Sprintf("databases.%s", db), conf)
	return viper.WriteConfig()
}
