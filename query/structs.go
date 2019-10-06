package query

// TODO comment all this file
import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"path"
)

type File struct {
	Queries []*Query `yaml:"queries"`
}

type Query struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	SQL         string `yaml:"query"`
}

type QueryConfig struct {
	Path      string
	Whitelist []string
	Params    utils.Params
}

// Defines a query file configuration.
// It is composed of a filename (with its path from the config folder) and a list of queries to execute.
type Queries []*QueryConfig

func (config *QueryConfig) Open() ([]*Query, error) {
	return OpenFile(config.Path, config.Whitelist, config.Params)
}

func (entry *Query) Execute(db *pg.DB, params ...interface{}) (pg.Result, error) {
	return db.Exec(entry.SQL, params...)
}

func (entry *Query) Export(db *pg.DB, file io.Writer, headers bool) error {
	logs.Log.WithFields(logrus.Fields{
		"query":       entry.Name,
		"description": entry.Description,
	}).Info("Executing query...")

	stmt := "COPY (" + entry.SQL + ") TO STDOUT WITH CSV"
	if headers {
		stmt += " HEADER"
	}
	result, err := db.CopyTo(file, stmt)
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"query": entry.Name,
		}).Errorf("Error executing query: %s", err)
		return err
	}
	logs.Log.WithFields(logrus.Fields{
		"query":         entry.Name,
		"rows_affected": result.RowsAffected(),
	}).Info("Query done!")
	return nil
}

// Opens a query file and returns the queries specified in whitelist
func OpenFile(queryFilename string, whitelist []string, params utils.Params) ([]*Query, error) {
	if params == nil {
		params = make(utils.Params)
	}
	queryFile := &File{}
	filepath, err := GetQueriesPath()
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(path.Join(filepath, queryFilename))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, queryFile)
	if err != nil {
		return nil, err
	}
	queries := make([]*Query, 0)
	whitelistMap := make(map[string]struct{})
	for _, whitelisted := range whitelist {
		whitelistMap[whitelisted] = struct{}{}
	}
	for _, query := range queryFile.Queries {
		if _, ok := whitelistMap[query.Name]; ok || len(whitelistMap) == 0 {
			queries = append(queries, query.Format(params))
		}
	}
	return queries, nil
}

// Creates a new Queries instance with the edited fields.
func (queries Queries) Format(params utils.Params) Queries {
	newList := make(Queries, len(queries))
	for i, queryConfig := range queries {
		newParams := params.Join(queryConfig.Params)
		newList[i] = &QueryConfig{
			Path:      params.FormatString(queryConfig.Path),
			Params:    newParams,
			Whitelist: params.FormatStringArray(queryConfig.Whitelist),
		}
	}
	return newList
}

// Creates a new Query instance with the edited fields.
func (entry *Query) Format(params utils.Params) *Query {
	return &Query{
		Name:        params.FormatString(entry.Name),
		Description: params.FormatString(entry.Description),
		SQL:         params.FormatString(entry.SQL),
	}
}

func (file *File) Format(params utils.Params) *File {
	newFile := &File{
		Queries: make([]*Query, len(file.Queries)),
	}
	if file.Queries != nil {
		for i, query := range file.Queries {
			newFile.Queries[i] = query.Format(params)
		}
	}
	return newFile
}
