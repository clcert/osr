package databases

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Defines the structure of a clickhouse connection string.
const clickhouseConnStr = "tcp://%s:%d?database=%s&username=%s&password=%s"

// Returns a sqlx struct with a "connection" to a clickhouse database
// with read permissions.
func GetClickhouseReader() (*sqlx.DB, error) {
	conf, err := GetDBConfig("clickhouse")
	if err != nil {
		return nil, err
	}
	fConnStr := fmt.Sprintf(clickhouseConnStr,
		conf.Server,
		conf.Port,
		conf.DBName,
		conf.Reader.Username,
		conf.Reader.Password)
	return sqlx.Open("clickhouse", fConnStr)
}

// Returns a sqlx struct with a "connection" to a clickhouse database
// with read and write permissions
func GetClickhouseWriter() (*sqlx.DB, error) {
	conf, err := GetDBConfig("clickhouse")
	if err != nil {
		return nil, err
	}
	fConnStr := fmt.Sprintf(clickhouseConnStr,
		conf.Server,
		conf.Port,
		conf.DBName,
		conf.Reader.Username,
		conf.Reader.Password)
	return sqlx.Open("clickhouse", fConnStr)
}
