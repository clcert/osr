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

// Defines a query file configuration.
// It is composed of a filename (with its path from the config folder) and a list of queries to execute.
type FileMap map[string][]string

func (entry *Query) Execute(db *pg.DB, params... interface{}) (pg.Result, error) {
	return db.Exec(entry.SQL, params...)
}

func (entry *Query) Export(db *pg.DB, file io.Writer, headers bool) error {
	logs.Log.WithFields(logrus.Fields{
		"query":       entry.Name,
		"description": entry.Description,
	}).Info("Executing query...")

	stmt := "COPY (" + entry.SQL + ") TO STDOUT"
	if headers {
		stmt += " WITH CSV HEADER"
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
func OpenFile(queryFilename string, whitelist ...string) ([]*Query, error) {
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
	if len(whitelist) == 0 {
		return queryFile.Queries, nil
	} else {
		queries := make([]*Query, 0)
		for _, query := range queryFile.Queries {
			for _, whitelisted := range whitelist {
				if query.Name == whitelisted {
					queries = append(queries, query)
					break
				}
			}
		}
		return queries, nil
	}
}

// Creates a new FileMap instance with the edited fields.
func (fileMap FileMap) Format(params utils.Params) FileMap {
	newMap := make(FileMap)
	for k, vList := range fileMap {
		newMap[params.FormatString(k)] = params.FormatStringArray(vList)
	}
	return newMap
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
