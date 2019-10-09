package databases

import (
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"strings"
)

// NewPostgresUser creates a new postgres user using the given DB connection, with a given username and a list of default table and sequence permissions.
// If the user already exists, it finishes with an error
func NewPostgresUser(db *pg.DB, username string, tablePermissions []string, seqPermissions []string) (*Credentials, error) {
	creds := Credentials{
		Username: username,
		Password: utils.GenerateRandomString(24),
	}
	prepStmt, err := db.Prepare("SELECT count(*) FROM pg_roles where rolname = $1")
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Error("Err trying to check if user exists (preparing statement).")
		return nil, err
	}

	var userExists int
	if _, err := prepStmt.QueryOne(pg.Scan(&userExists), creds.Username); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Error("Err trying to check if user exists (executing query).")
		return nil, err
	}
	if userExists > 0 {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Error("Cannot continue: User already exists!")
		return nil, fmt.Errorf("user %s already exists", username)
	}
	logs.Log.WithFields(logrus.Fields{
		"username": username,
	}).Info("Creating user...")
	if _, err := db.Exec(
		fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", creds.Username, creds.Password)); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Errorf("User creation failed: %s", err)
		return nil, err
	}
	if _, err := db.Exec(
		fmt.Sprintf("GRANT %s ON ALL TABLES IN SCHEMA public TO %s",
			strings.Join(tablePermissions, ", "), creds.Username)); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Errorf("privileges config (tables) failed: %s", err)
		return nil, err
	}
	if _, err := db.Exec(
		fmt.Sprintf("GRANT %s ON ALL SEQUENCES IN SCHEMA public TO %s",
			strings.Join(seqPermissions, ", "), creds.Username)); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Errorf("privileges config (sequences) failed: %s", err)
		return nil, err
	}
	if _, err := db.Exec(
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT %s ON TABLES TO %s",
			strings.Join(tablePermissions, ", "), creds.Username)); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Errorf("default privileges config (sequences) failed: %s", err)
		return nil, err
	}
	if _, err := db.Exec(
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT %s ON SEQUENCES TO %s",
			strings.Join(seqPermissions, ", "), creds.Username)); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Errorf("default privileges config (sequences) failed: %s", err)
		return nil, err
	}
	return &creds, nil
}

// GetPostgresReader returns a pg.DB struct with a "connection" to a Postgres database with read permissions.
func GetPostgresReader() (*pg.DB, error) {
	conf, err := GetDBConfig("postgres")
	return pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Server, conf.Port),
		User:     conf.Reader.Username,
		Password: conf.Reader.Password,
		Database: conf.DBName,
	}), err
}

// GetPostgresWriter returns a pg.DB struct with a "connection" to a Postgres database
// with read and write permissions.
func GetPostgresWriter() (*pg.DB, error) {
	conf, err := GetDBConfig("postgres")
	if err != nil {
		return nil, err
	}
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Server, conf.Port),
		User:     conf.Writer.Username,
		Password: conf.Writer.Password,
		Database: conf.DBName,
	})
	return db, err
}
