package sources

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/query"
	"github.com/clcert/osr/remote"
	"github.com/clcert/osr/utils"
	"github.com/pkg/sftp"
	"io"
	"path"
)

// TODO: More loging on this source
// SFTPConfig defines a configuration for a SFTP Source
type SFTPConfig struct {
	ServerName string        // Name of the server as declared on OSR config
	Path       string        // Path where the process will start to look for the entries
	Script     string        // Local script Path. This script will be uploaded to Path.
	Queries    query.Queries // Saved Queries on files on queries folder.
	Filter     *FilterConfig // Filter configuration
}

// SFTPSource defines a remote source of files, connected via SFTP.
type SFTPSource struct {
	*SFTPConfig           // Configuration
	name   string         // Source name
	server *remote.Server // Server name as defined on OSR Config
	filter *Filter        // Filter to use on document finding
	files  chan Entry     // Channel for files found.
	log    *logs.OSRLog   // Logs
	params utils.Params   // Process Parameters
}

// SFTPFile represents a SFTP entry source.
type SFTPFile struct {
	path   string      // Complete Path of the file
	source *SFTPSource // Related source
	file   *sftp.File  // File object
	buffer io.Reader   // Buffered extReader
}

// Creates a new SFTP source from a configuration
func (config *SFTPConfig) New(name string, params utils.Params) (source *SFTPSource, err error) {
	if config.Filter == nil {
		config.Filter = &FilterConfig{}
	}
	err = config.Format(params)
	if err != nil {
		return
	}
	var filter *Filter
	filter, err = config.Filter.New()
	if err != nil {
		return
	}
	log, err := logs.NewLog(name)
	if err != nil {
		return
	}
	source = &SFTPSource{
		name:       name,
		SFTPConfig: config,
		files:      make(chan Entry),
		filter:     filter,
		log:        log,
		params:     params,
	}
	return
}

// Format formats the configuration, using the params defined in the task
func (config *SFTPConfig) Format(params utils.Params) error {
	if params == nil {
		// do nothing
		return nil
	}
	config.Path = params.FormatString(config.Path)
	config.Queries = config.Queries.Format(params)
	config.Script = params.FormatString(config.Script)
	config.Filter = config.Filter.Format(params)
	return nil
}

func (source *SFTPSource) Init() error {
	if source.ServerName == "" || source.Path == "" {
		return fmt.Errorf("mandatory config fields not initialized")
	}

	// Connect to ServerName
	server, err := remote.GetServer(source.ServerName)
	if err != nil {
		return err
	}
	err = server.Connect()
	if err != nil {
		return err
	}
	source.server = server

	if source.Path == "" {
		// We use the default one
		client, err := server.GetSFTPClient()
		if err != nil {
			return err
		}
		folder, err := remote.CreateDefaultProcessFolder(client, source.name)
		if err != nil {
			return err
		}
		source.Path = folder
	}

	if source.Queries != nil {
		if err := executeQueries(source.server, source.Path, source.Queries, source.params); err != nil {
			return err
		}
	}

	if source.Script != "" {
		if err := source.server.ExecuteLocalScript(source.Path, source.Script, source.params); err != nil {
			return err
		}
	}
	go source.retrieveFiles()

	return nil
}

func (source *SFTPSource) Close() error {
	return source.server.Close()
}

func (source *SFTPSource) Next() Entry {
	entry, ok := <-source.files
	if !ok {
		return nil
	}
	return entry
}

func (source *SFTPSource) retrieveFiles() {
	defer close(source.files)
	sftpClient, err := source.server.GetSFTPClient()
	if err != nil {
		return
	}
	walker := sftpClient.Walk(source.Path)
	for walker.Step() {
		if walker.Stat() == nil {
			// TODO: log this
			break // folder doesn't exist
		}
		if walker.Stat().IsDir() {
			if walker.Path() != source.Path && !source.filter.Recursive {
				walker.SkipDir()
				continue
			}
		} else {
			aFile := &SFTPFile{
				source: source,
				path:   walker.Path(),
			}
			if len(source.filter.Patterns) == 0 {
				source.files <- aFile
			}
			for _, regex := range source.filter.Patterns {
				if regex.MatchString(walker.Path()) {
					source.files <- aFile
					break
				}
			}
		}
	}
	return
}

func (source *SFTPSource) GetID() (string, error) {
	if source.server == nil {
		return "", fmt.Errorf("source not initialized")
	}

	addr, err := source.server.GetAddr()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (source *SFTPSource) GetName() string {
	return source.name
}

func (source *SFTPSource) GetAttachments() []string {
	return []string{source.log.Path, source.server.Output.Path}
}

func (srcFile *SFTPFile) Open() (io.Reader, error) {
	sftpClient, err := srcFile.source.server.GetSFTPClient()
	if err != nil {
		return nil, err
	}
	if srcFile.file == nil {
		file, err := sftpClient.Open(srcFile.path)
		if err != nil {
			return nil, err
		}
		srcFile.file = file
	}
	if srcFile.buffer == nil {
		srcFile.buffer = bufio.NewReader(srcFile.file)
	}
	return srcFile.buffer, nil
}

func (srcFile *SFTPFile) Name() string {
	return path.Base(srcFile.path)
}

func (srcFile *SFTPFile) Dir() string {
	return path.Dir(srcFile.path)
}

func (srcFile *SFTPFile) Path() string {
	return srcFile.path
}

func (srcFile *SFTPFile) Close() error {
	srcFile.buffer = nil
	return srcFile.file.Close()
}
