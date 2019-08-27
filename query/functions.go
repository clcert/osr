package query

import (
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/logs"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// Export executes queries defiend in inFolder. If whitelist is not empty, it uses only the queries of it.
func Execute(queryFiles []string, outFolder string, whitelist []string, headers bool) error {
	db, err := databases.GetPostgresReader()
	if err != nil {
		return err
	}

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
		queries, err := OpenFile(queryFile, whitelist...)
		if err != nil {
			return err
		}
		for _, query := range queries {
			newFilePath := filepath.Join(outFolderPath, query.Name+".csv")
			logs.Log.WithFields(logrus.Fields{
				"new_file_path": newFilePath,
			}).Info("Opening destination file...")
			newFile, err := os.Create(newFilePath)
			if err != nil {
				return err
			}
			if err := query.Export(db, newFile, headers); err != nil {
				return err
			}
		}
	}
	return nil
}
