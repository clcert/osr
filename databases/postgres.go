package databases

import (
	"context"
	"fmt"
	"strings"

	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	fmt.Println(q.FormattedQuery())
	return c, nil
}

func enablePartman(db *pg.DB, tablePermissions []string, creds *Credentials) error {
	partmanStmts := []string{
		"CREATE SCHEMA IF NOT EXISTS partman",
		"CREATE EXTENSION IF NOT EXISTS pg_partman SCHEMA partman",
		fmt.Sprintf("GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA partman TO %s", creds.Username),
		fmt.Sprintf("GRANT EXECUTE ON ALL PROCEDURES IN SCHEMA partman TO %s", creds.Username),
		fmt.Sprintf("GRANT ALL ON ALL TABLES IN SCHEMA partman TO %s", creds.Username),
		fmt.Sprintf("GRANT ALL ON SCHEMA partman TO %s", creds.Username),
	}
	for _, stmt := range partmanStmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("Error on statement %s: %s", stmt, err)
		}
	}
	return nil
}

func fixUserPermissions(db *pg.DB, schema string, tablePermissions, seqPermissions []string, creds *Credentials) error {
	userStmts := []string{
		fmt.Sprintf("GRANT %s ON ALL TABLES IN SCHEMA %s TO %s",
			strings.Join(tablePermissions, ", "), schema, creds.Username),
		fmt.Sprintf("GRANT %s ON ALL SEQUENCES IN SCHEMA %s TO %s",
			strings.Join(seqPermissions, ", "), schema, creds.Username),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT %s ON TABLES TO %s",
			schema, strings.Join(tablePermissions, ", "), creds.Username),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT %s ON SEQUENCES TO %s",
			schema, strings.Join(seqPermissions, ", "), creds.Username),
	}
	for _, stmt := range userStmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("Error on statement %s: %s", stmt, err)
		}
	}
	return nil
}

// NewPostgresUser creates a new postgres user using the given DB connection, with a given username and a list of default table and sequence permissions.
// If the user already exists, it finishes with an error
func NewPostgresUser(db *pg.DB, username string, tablePermissions, seqPermissions []string) (*Credentials, error) {
	creds := &Credentials{
		Username: username,
		Password: utils.GenerateRandomString(24),
	}

	logs.Log.WithFields(logrus.Fields{
		"username": username,
	}).Info("Checking if user exists...")

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

	if _, err := db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", creds.Username, creds.Password)); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": creds.Username,
		}).Errorf("User creation failed: %s", err)
		return nil, err
	}

	// Fixing permissions
	logs.Log.WithFields(logrus.Fields{
		"username": username,
	}).Info("Fixing permissions...")
	if err := fixUserPermissions(db, "public", tablePermissions, seqPermissions, creds); err != nil {
		return nil, err
	}

	// Optional: Partman extension (?)
	logs.Log.Info("Trying to enable partman extension (for automatic partitioning...)")
	if err := enablePartman(db, tablePermissions, creds); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"username": creds.Username,
		}).Errorf("%s", err)
	} else {
		logs.Log.WithFields(logrus.Fields{
			"username": username,
		}).Info("Fixing partman permissions...")
		if err := fixUserPermissions(db, "partman", tablePermissions, seqPermissions, creds); err != nil {
			logs.Log.WithFields(logrus.Fields{
				"username": creds.Username,
			}).Errorf("%s", err)
		}
	}

	return creds, nil
}

// GetPostgresReader returns a pg.DB struct with a "connection" to a Postgres database with read permissions.
func GetPostgresReader() (*pg.DB, error) {
	conf, err := GetDBConfig("postgres")
	if err != nil {
		return nil, err
	}
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Server, conf.Port),
		User:     conf.Reader.Username,
		Password: conf.Reader.Password,
		Database: conf.DBName,
		OnConnect: func(ctx context.Context, cn *pg.Conn) error {
			_, err := cn.Exec("set TIMEZONE='America/Santiago'")
			return err
		},
	})
	// db.AddQueryHook(dbLogger{})
	return db, err
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
	// 	db.AddQueryHook(dbLogger{})
	return db, err
}
