package sources

// TODO: move this somewhere else
import (
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/query"
	"github.com/clcert/osr/remote"
	"github.com/clcert/osr/utils"
	"path"
)

func executeQueries(server *remote.Server, remotePath string, queryFiles query.FileMap, params utils.Params) error {
	db, err := databases.GetPostgresReader()
	if err != nil {
		return err
	}
	defer db.Close()
	sftpClient, err := server.GetSFTPClient()
	if err != nil {
		return err
	}
	for queryFile, whitelist := range queryFiles {
		queries, err := query.OpenFile(queryFile, whitelist...)
		if err != nil {
			// TODO: Log this
			continue
		}
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
			aQuery = aQuery.Format(params)
			if err := aQuery.Export(db, outFile, false); err != nil {
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
