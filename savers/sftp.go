package savers

import (
	"encoding/csv"
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/remote"
	"github.com/clcert/osr/utils"
	"github.com/pkg/sftp"
	"path"
	"strings"
	"sync"
)

// TODO: More logging on this source
// SFTPConfig defines a configuration for a SFTP saver.
type SFTPConfig struct {
	ServerName string                     // Name of the server as defined in the server configuration of OSR
	Path       string                     // Init path. If the folder doesn't exist, it's created
	ForceFlush bool                       // If it's true, the file lines are commited immediately. Its slower but safer (if the routine crashes, it saves some data).
	FileConfig map[string]*SFTPFileConfig // A map with the outID of the savables as the key, and a configuration as a value.
}

// SFTPFileConfig defines a configuration for a specific file or outID.
type SFTPFileConfig struct {
	FileName string   // Name of the file. It should be unique because all the files are saved on the same folder
	Fields   []string // Fields to save. TODO: make this case insensitive.
}

// SFTPSaver defines a saver which saves the objects in a CSV file, on a specific folder in a remote server.
type SFTPSaver struct {
	*SFTPConfig                     // SFTPConfig related to the saver
	sync.Mutex                      // A mutex for map operations
	name       string               // Name for the saver instance
	params     utils.Params         // Parameters configured in file
	server     *remote.Server       //the server object
	sftpClient *sftp.Client         // The sftp Connection
	outFiles   map[string]*SFTPFile // A map with file writers, where the key is the outID.
	finished   chan bool            // True if channel finished parsing files
	objects    chan Savable         // List of objects to save
	inserted   int                  // number of inserted data rows
	errors     []error              // List of errors
	log        *logs.OSRLog         // Saver log
}

// SFTPFile defines a specific file where to save the objects.
type SFTPFile struct {
	file   *sftp.File
	writer *csv.Writer
}

// New uses the configuration defined in a SFTPConfig to create a new instance.
func (config *SFTPConfig) New(name string, params utils.Params) (*SFTPSaver, error) {
	err := config.Format(params)
	if err != nil {
		return nil, err
	}
	log, err := logs.NewLog(name)
	if err != nil {
		return nil, err
	}
	return &SFTPSaver{
		SFTPConfig: config,
		name:       name,
		finished:   make(chan bool, 1),
		objects:    make(chan Savable),
		errors:     make([]error, 0),
		outFiles:   make(map[string]*SFTPFile, 0),
		inserted:   0,
		log:        log,
	}, nil
}

func (config *SFTPConfig) Format(params utils.Params) error {
	if params == nil {
		return nil
	}
	config.ServerName = params.FormatString(config.ServerName)
	config.Path = params.FormatString(config.Path)
	newFileConfig := make(map[string]*SFTPFileConfig)
	for k, v := range config.FileConfig {
		newFileConfig[params.FormatString(k)] = v.Format(params)
	}
	config.FileConfig = newFileConfig
	return nil
}

func (saver *SFTPSaver) Start() error {

	// Connect to server
	server, err := remote.GetServer(saver.ServerName)
	if err != nil {
		return err
	}
	saver.server = server
	err = saver.server.Connect()
	if err != nil {
		return err
	}
	sftpClient, err := saver.server.GetSFTPClient()
	if err != nil {
		saver.server.Close()
		return err
	}
	saver.sftpClient = sftpClient

	// Create outFiles and prepare serializers

	for outName, outConfig := range saver.FileConfig {
		if err := saver.createOutFile(outName, outConfig); err != nil {
			return err
		}
	}
	go func() {
		for newObject := range saver.objects {
			saver.writeToFile(newObject)
		}
		saver.finished <- true
	}()
	return nil
}

// There are no messages (yet) for this saver
func (saver *SFTPSaver) SendMessage(msg interface{}) error {
	return nil
}

func (saver *SFTPSaver) Save(objs ...interface{}) error {
	for _, obj := range objs {
		switch obj.(type) {
		case Savable:
			saver.objects <- obj.(Savable)
		default:
			saver.objects <- Savable{Object: obj}
		}
	}
	return nil
}

func (saver *SFTPSaver) Finish() error {
	close(saver.objects)
	<-saver.finished
	for _, outFile := range saver.outFiles {
		outFile.writer.Flush()
		outFile.file.Close()
	}
	err := saver.server.Close()
	return err
}

func (saver *SFTPSaver) GetErrors() []error {
	return saver.errors
}

func (saver *SFTPSaver) GetAttachments() []string {
	return []string{saver.log.Path}
}

func (saver *SFTPSaver) GetName() string {
	return saver.name
}

func (saver *SFTPSaver) writeToFile(savable Savable) {
	outID := savable.GetOutID()
	file, ok := saver.outFiles[outID]
	if !ok {
		structName := savable.StructName()
		file, ok = saver.outFiles[structName]
		if !ok {

			saver.FileConfig[outID] = &SFTPFileConfig{
				FileName: strings.Replace(outID, "/", "-", -1),
				Fields:   savable.FieldNames(), // all the fields
			}

			if err := saver.createOutFile(outID, saver.FileConfig[outID]); err != nil {
				fmt.Errorf("object could not be saved. The file didn't exist and it was impossible to create it")
				return
			}

			// It should exist now
			file = saver.outFiles[outID]
		}
	}
	allFields := savable.Fields()
	outConfig, ok := saver.FileConfig[outID]
	if !ok {
		saver.errors = append(saver.errors, fmt.Errorf("there is no config for this file type"))
	}
	values := make([]string, len(outConfig.Fields))
	for i, field := range outConfig.Fields {
		fieldValue, ok := allFields[field]
		if ok {
		values[i] = fmt.Sprintf("%v", fieldValue)
		}
	}
	err := file.writer.Write(values)
	if saver.ForceFlush {
		file.writer.Flush()
	}
	if err != nil {
		saver.errors = append(saver.errors, err)
	}
}

func (saver *SFTPSaver) createOutFile(name string, config *SFTPFileConfig) error {
	if saver.sftpClient == nil {
		return fmt.Errorf("sftp connection not initialized")
	}
	if len(config.Fields) == 0 {
		return fmt.Errorf("must declare fields to write on file")
	}
	// create folder in remote
	err := saver.sftpClient.MkdirAll(saver.Path)
	if err != nil {
		return err
	}
	// Create file in remote
	aFile, err := saver.sftpClient.Create(path.Join(saver.Path, config.FileName))
	if err != nil {
		return err
	}
	saver.outFiles[name] = &SFTPFile{
		file:   aFile,
		writer: csv.NewWriter(aFile),
	}
	err = saver.outFiles[name].writer.Write(config.Fields)
	if saver.ForceFlush {
		saver.outFiles[name].writer.Flush()
	}
	if err != nil {
		saver.outFiles[name].file.Close()
		return err
	}
	return nil
}

func (fileConfig *SFTPFileConfig) Format(params utils.Params) *SFTPFileConfig {
	newFileConfig := &SFTPFileConfig{
		FileName: params.FormatString(fileConfig.FileName),
		Fields:   make([]string, len(fileConfig.Fields)),
	}
	for i, v := range fileConfig.Fields {
		newFileConfig.Fields[i] = params.FormatString(v)
	}
	return newFileConfig
}
