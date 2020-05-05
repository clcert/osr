package sources

// TODO: move this somewhere else
import (
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/query"
	"github.com/clcert/osr/remote"
	"github.com/clcert/osr/utils"
	"path"
)

func executeQueries(server *remote.Server, remotePath string, queryFiles query.Configs, params utils.Params) error {
	db, err := databases.GetPostgresReader()
	if err != nil {
		return err
	}
	defer db.Close()
	sftpClient, err := server.GetSFTPClient()
	if err != nil {
		return err
	}
	for _, queryConfig := range queryFiles {
		queries, err := queryConfig.Open()
		if err != nil {
			// TODO: Log this
			continue
		}
		formatted := make(map[string]*query.Query, 0)
		for _, aQuery := range queries {
			outPath := path.Join(remotePath, aQuery.Name+".csv")
			outFile, err := sftpClient.Create(outPath)
			if err != nil {
				if outFile != nil {
					outFile.Close()
				}
				return err
			}
			// Formatting query
			aQuery, err = aQuery.Format(params, formatted)
			if err != nil {
				return err
			}
			formatted[aQuery.Name] = aQuery
			chErr := aQuery.Export(db, outFile, false);
			if <-chErr != nil {
				if outFile != nil {
					outFile.Close()
				}
				return err
			}
			outFile.Close()
		}
	}
	return nil
}
