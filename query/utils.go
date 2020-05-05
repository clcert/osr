package query

import (
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)


// Export executes queries defined in inFolder. If whitelist is not empty, it uses only the queries of it.
func Execute(queryFiles []string, outFolder string, whitelist []string, headers bool, params []string) error {
	db, err := databases.GetPostgresReader()
	if err != nil {
		return err
	}
	cmdParams := utils.ListToParams(params)
	for _, queryFile := range queryFiles {

		logs.Log.WithFields(logrus.Fields{
			"queryFile": queryFile,
		}).Info("Opening query file...")
		outFolderPath := filepath.Join(outFolder, strings.TrimSuffix(queryFile, ".yaml"))
		logs.Log.WithFields(logrus.Fields{
			"output_folder": outFolderPath,
		}).Info("Creating output folder")
		if err := os.MkdirAll(outFolderPath, 0755); err != nil {
			return err
		}
		queries, err := OpenFile(queryFile, cmdParams)
		if err != nil {
			return err
		}
		whitelisted := make(map[string]*Query)
		if len(whitelist) > 0 {
			for _, queryName := range whitelist {
				if query, ok := queries[queryName]; ok {
					whitelisted[queryName] = query
				}
			}
		} else {
			whitelisted = queries
		}
		for _, query := range whitelisted {
			newFilePath := filepath.Join(outFolderPath, query.Name+".csv")
			logs.Log.WithFields(logrus.Fields{
				"new_file_path": newFilePath,
			}).Info("Opening destination file...")
			newFile, err := os.Create(newFilePath)
			if err != nil {
				return err
			}
			chErr := query.Export(db, newFile, headers);
			if <-chErr != nil {
				return err
			}
		}
	}
	return nil
}



// Helper function which returns the queries locations using the config values.
func GetQueriesPath() (string, error) {
	home := viper.GetString("folders.home")
	logs := viper.GetString("folders.queries")
	return filepath.Join(home, logs), nil
}

// Opens a query file and returns the queries specified in whitelist
func OpenFile(queryFilename string, params utils.Params) (map[string]*Query, error) {
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
	queries := make(map[string]*Query, 0)
	for _, query := range queryFile.Queries {
		newQuery, err := query.Format(params, queries)
		if err != nil {
			return nil, err
		}
		queries[newQuery.Name] = newQuery
	}
	return queries, nil
}

// formatQuery formats a query using the queries formatted until now
func formatQuery(params utils.Params, query *Query, queries map[string]*Query) (string, error) {
	tmpl, err := template.New(query.Name).Parse(query.SQL)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, &FormatArgs{
		FormatArgs: utils.NewFormatArgs(params),
		Queries:    queries,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
