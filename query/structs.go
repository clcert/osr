package query

// TODO comment all this file
import (
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"io"
)

type File struct {
	Queries []*Query `yaml:"queries"`
}

type Query struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	SQL         string `yaml:"query"`
}

type FormatArgs struct {
	*utils.FormatArgs
	Queries map[string]*Query
}

type Config struct {
	Path      string
	Whitelist []string
	Params    utils.Params
}

// Defines a query file configuration.
// It is composed of a filename (with its path from the config folder) and a list of queries to execute.
type Configs []*Config

func (config *Config) Open() (map[string]*Query, error) {
	queries, err := OpenFile(config.Path, config.Params)
	if err != nil {
		return nil, err
	}
	whitelisted := make(map[string]*Query)
	if len(config.Whitelist) > 0 {
		for _, queryName := range config.Whitelist {
			if query, ok := queries[queryName]; ok {
				whitelisted[queryName] = query
			}
		}
	} else {
		whitelisted = queries
	}
	return whitelisted, nil
}

func (entry *Query) Execute(db *pg.DB, params ...interface{}) (pg.Result, error) {
	return db.Exec(entry.SQL, params...)
}

func (entry *Query) Export(db *pg.DB, file io.Writer, headers bool) chan error {
	chErr := make(chan error)
	logs.Log.WithFields(logrus.Fields{
		"query":       entry.Name,
		"description": entry.Description,
		"SQL":         entry.SQL,
	}).Info("Executing query...")

	stmt := "COPY (" + entry.SQL + ") TO STDOUT WITH CSV"
	if headers {
		stmt += " HEADER"
	}
	go func() {
		result, err := db.CopyTo(file, stmt)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"query": entry.Name,
			}).Errorf("Error executing query: %s", err)
		} else {
			logs.Log.WithFields(logrus.Fields{
				"query":         entry.Name,
				"rows_affected": result.RowsAffected(),
			}).Info("Query done!")
		}
		chErr <- err
		close(chErr)
	}()
	return chErr
}

// Creates a new Configs instance with the edited fields.
func (queries Configs) Format(params utils.Params) Configs {
	newList := make(Configs, len(queries))
	for i, queryConfig := range queries {
		newParams := params.Join(queryConfig.Params)
		newList[i] = &Config{
			Path:      params.FormatString(queryConfig.Path),
			Params:    newParams,
			Whitelist: params.FormatStringArray(queryConfig.Whitelist),
		}
	}
	return newList
}

// Creates a new Query instance with the edited fields.
func (entry *Query) Format(params utils.Params, queries map[string]*Query) (*Query, error) {
	q := &Query{
		Name:        params.FormatString(entry.Name),
		Description: params.FormatString(entry.Description),
	}
	sql, err := formatQuery(params, entry, queries)
	if err != nil {
		return nil, err
	}
	q.SQL = sql
	return q, nil
}

func (file *File) Format(params utils.Params) (*File, error) {
	newFile := &File{
		Queries: make([]*Query, 0),
	}
	queries := make(map[string]*Query)
	if file.Queries != nil {
		for _, query := range file.Queries {
			newQuery, err := query.Format(params, queries)
			if err != nil {
				return nil, err
			}
			queries[newQuery.Name] = newQuery
			newFile.Queries = append(newFile.Queries, newQuery)
		}
	}
	return newFile, nil
}

// Query returns a query already parsed in the file (defined before)
func (args *FormatArgs) Query(name string) string {
	if _, ok := args.Queries[name]; !ok {
		return ""
	}
	return fmt.Sprintf("(%s)", args.Queries[name].SQL)
}
