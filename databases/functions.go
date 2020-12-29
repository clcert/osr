package databases

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

// GetDBConfigData asks interactively for Postgres config parameters.
// It asks first for DB connection data and the name for the new DB to create.
// It also asks for admin user credenmtials, used only once.
func GetDBConfigData() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the address and port of the PostgresDB:")
	fmt.Print("Address: [localhost] ")
	scanner.Scan()
	address := scanner.Text()
	if len(address) == 0 {
		address = "localhost"
	}
	fmt.Print("Port: [5432] ")
	scanner.Scan()
	port := scanner.Text()
	if len(port) == 0 {
		port = "5432"
	}
	fmt.Print("Database Name: [osr] ")
	scanner.Scan()
	dbname := scanner.Text()
	if len(dbname) == 0 {
		dbname = "osr"
	}
	fmt.Println("Please enter the username and password of the admin user of the Postgres DB:")
	fmt.Print("Username: [postgres] ")
	scanner.Scan()
	username := scanner.Text()
	if len(username) == 0 {
		username = "postgres"
	}
	fmt.Print("Password (hidden): ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	password := string(bytePassword)
	fmt.Println("")
	return InitDatabase(address, port, dbname, username, password)
}

// InitDatabase initializes the database using the admin credentials, creating reader and writer users and saving their credentials into the config file.
func InitDatabase(address, port, dbname, username, password string) error {
	intPort, err := strconv.ParseInt(port, 10, 16)
	if err != nil {
		return err
	}
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", address, intPort),
		User:     username,
		Password: password,
	})
	defer db.Close()
	prepStmt, err := db.Prepare("SELECT count(*) FROM pg_database WHERE datname= $1")
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"database": dbname,
		}).Error("Err trying to check if database exists (preparing statement).")
		return err
	}

	var dbExists int
	if _, err := prepStmt.QueryOne(pg.Scan(&dbExists), dbname); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"database": dbname,
		}).Error("Err trying to check if database exists (executing query).")
		return err
	}
	if dbExists != 0 {
		logs.Log.WithFields(logrus.Fields{
			"database": dbname,
		}).Error("Database already exists. Try with another name or delete it before retrying.")
		return fmt.Errorf("database %s already exists", dbname)
	}
	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)); err != nil {
		logs.Log.WithFields(logrus.Fields{
			"database": dbname,
		}).Error("Couldn't create database. has the user provided enough permissions to create a database?")
		return err
	}
	db = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", address, port),
		User:     username,
		Password: password,
		Database: dbname,
	})
	defer db.Close()
	readerUser := fmt.Sprintf("%s_reader_%s", dbname, utils.GenerateRandomHex(6))
	readerCreds, err := NewPostgresUser(db, readerUser,
		[]string{"SELECT"},
		[]string{"USAGE", "SELECT"},
	)
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"database":      dbname,
			"readuser_name": readerUser,
		}).Errorf("Couldn't create read user: %s (has the user provided enough permissions to create an user?)", err)
		// Remove created database.
		if _, err := db.Exec(fmt.Sprintf("DROP DATABASE %s", dbname)); err != nil {
			logs.Log.WithFields(logrus.Fields{
				"database": dbname,
			}).Error("Couldn't remove database. has the user provided enough permissions to remove a database?")
			return err
		}

		return err
	}
	writerUser := fmt.Sprintf("%s_writer_%s", dbname, utils.GenerateRandomHex(6))
	writerCreds, err := NewPostgresUser(db, writerUser,
		[]string{"SELECT", "INSERT", "UPDATE", "DELETE", "TRUNCATE", "REFERENCES", "TRIGGER"},
		[]string{"USAGE", "SELECT", "UPDATE"},
	)
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"database":       dbname,
			"writeuser_name": writerUser,
		}).Error("Couldn't write user. has the user provided enough permissions to create an user?")

		// Remove created database.
		if _, err := db.Exec(fmt.Sprintf("DROP DATABASE %s", dbname)); err != nil {
			logs.Log.WithFields(logrus.Fields{
				"database": dbname,
			}).Error("Couldn't remove database. has the user provided enough permissions to remove a database?")
			return err
		}

		return err
	}
	conf := Config{
		Server: address,
		Port:   int(intPort),
		Reader: *readerCreds,
		Writer: *writerCreds,
		DBName: dbname,
	}
	return WriteDBConf("postgres", conf)
}
