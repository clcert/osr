package sources

import (
	"fmt"
	"io"
	"path"

	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/query"
	"github.com/clcert/osr/utils"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

// TODO: More logging on this source
// QueryListConfig defines the configuration that a Query Source uses
type QueryListConfig struct {
	Path    string        // Virtual Path
	Queries query.Configs // List of queries to execute and upload to the Path folder
}

// QuerySource represents a source that is a group of sql queries.
type QuerySource struct {
	*QueryListConfig
	conn   *pg.DB          // DB Connection
	name   string          // Source name
	log    *logs.OSRLog    // Source log file
	params utils.Params    // Process Parameters
	files  chan *QueryFile // Channel of files
}

// QueryFile represents the command being executed.
type QueryFile struct {
	source *QuerySource
	query  *query.Query
	reader *io.PipeReader
	writer *io.PipeWriter
}

// New creates a new QuerySource from a QueryListConfig.
func (config *QueryListConfig) New(name string, params utils.Params) (source *QuerySource, err error) {
	err = config.Format(params)
	if err != nil {
		return
	}
	log, err := logs.NewLog(name)
	if err != nil {
		return
	}
	source = &QuerySource{
		name:            name,
		QueryListConfig: config,
		log:             log,
		params:          params,
		files:           make(chan *QueryFile),
	}
	return
}

// Format formats the configuration, using the params defined in the task
func (config *QueryListConfig) Format(params utils.Params) error {
	if params == nil {
		return nil
	}
	config.Path = params.FormatString(config.Path)
	config.Queries = config.Queries.Format(params)
	return nil
}

// Init initializes the source.
func (source *QuerySource) Init() error {
	if source.Queries == nil {
		return fmt.Errorf("there are no queries to execute")
	}
	// Connect to ServerName
	db, err := databases.GetPostgresReader()
	if err != nil {
		return err
	}
	source.conn = db

	go func() {
		for _, queryConfig := range source.Queries {
			queries, err := queryConfig.Open()
			if err != nil {
				source.log.WithFields(logrus.Fields{
					"query": queryConfig.Path,
					"whitelist": queryConfig.Whitelist,
					"params": queryConfig.Params,
				}).Infof("error opening query: %s", err)
				continue
			}
			for _, aQuery := range queries {
				reader, writer := io.Pipe()
				source.files <- &QueryFile{
					source: source,
					query:  aQuery,
					reader: reader,
					writer: writer,
				}
			}
		}
		close(source.files)
	}()
	return nil
}

func (source *QuerySource) Next() Entry {
	entry, ok := <-source.files
	if !ok {
		return nil
	}
	return entry
}

func (source *QuerySource) Close() error {
	err := source.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (source *QuerySource) GetID() (string, error) {
	if source.conn == nil {
		return "", fmt.Errorf("source not initialized")
	}

	opt := source.conn.Options()
	if opt == nil {
		return "", fmt.Errorf("cannot find IP")
	}
	return opt.Addr, nil
}

func (source *QuerySource) GetName() string {
	return source.name
}

// TODO: add queries as log
func (source *QuerySource) GetAttachments() []string {
	return []string{source.log.Path}
}

func (srcFile *QueryFile) Open() (io.Reader, error) {
	go func() {
		chErr := srcFile.query.Export(srcFile.source.conn, srcFile.writer, true)
		_ = srcFile.writer.CloseWithError(<-chErr)
	}()
	return srcFile.reader, nil
}

func (srcFile *QueryFile) Name() string {
	return srcFile.query.Name + ".csv" // we are csv aren't we?
}

func (srcFile *QueryFile) Dir() string {
	return srcFile.source.Path
}

func (srcFile *QueryFile) Path() string {
	return path.Join(srcFile.Dir(), srcFile.Name())
}

func (srcFile *QueryFile) Close() error {
	return srcFile.reader.Close()
}
