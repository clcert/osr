package sources

import (
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/query"
	"github.com/clcert/osr/remote"
	"github.com/clcert/osr/utils"
	"golang.org/x/crypto/ssh"
	"io"
	"path"
)

// TODO: More logging on this source
// ScriptConfig defines the configuration that a Script Source uses
type ScriptConfig struct {
	ServerName string        // Name of the server as defined in OSR config
	Name       string        // Name of the virtual file
	Script     string        // Local Path to script that is going to be executed on server
	Path       string        // Working Remote Directory used by the command and output of the queries
	Queries    query.Configs // List of queries to execute and upload to the Path folder
}

// ScriptSource represents a source that is a command in execution.
type ScriptSource struct {
	*ScriptConfig           // Script Configuration
	name     string         // Name of source
	consumed bool           // True if the command was retrieved with next
	server   *remote.Server // Remote server
	file     *ScriptFile    // The Script itself, as a file.
	log      *logs.OSRLog   // Source log file
	params   utils.Params   // Process Parameters
}

// ScriptFile represents the command being executed.
type ScriptFile struct {
	source  *ScriptSource
	script  string
	session *ssh.Session
}

// New creates a new ScriptSource from a ScriptConfig.
func (config *ScriptConfig) New(name string, params utils.Params) (source *ScriptSource, err error) {
	err = config.Format(params)
	if err != nil {
		return
	}
	log, err := logs.NewLog(name)
	if err != nil {
		return
	}
	source = &ScriptSource{
		name:         name,
		ScriptConfig: config,
		log:          log,
		params:       params,
	}
	return
}

// Format formats the configuration, using the params defined in the task
func (config *ScriptConfig) Format(params utils.Params) error {
	if params == nil {
		// do nothing
		return nil
	}
	config.ServerName = params.FormatString(config.ServerName)
	config.Name = params.FormatString(config.Name)
	config.Script = params.FormatString(config.Script)
	config.Path = params.FormatString(config.Path)
	config.Queries = config.Queries.Format(params)
	return nil
}

func (source *ScriptSource) Init() error {
	if source.ServerName == "" {
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
	scriptRoot, err := remote.GetScriptsPath()
	if err != nil {
		return err
	}
	paths, err := source.server.SendFiles(source.params, source.Path, path.Join(scriptRoot, source.Script))
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return err
	}
	source.file = &ScriptFile{
		source: source,
		script: paths[0],
	}

	if source.Queries != nil {
		if err := executeQueries(source.server, source.Path, source.Queries, source.params); err != nil {
			return err
		}
	}
	return nil
}

func (source *ScriptSource) Next() Entry {
	if !source.consumed {
		source.consumed = true
		return source.file
	} else {
		return nil
	}
}

func (source *ScriptSource) Reset() {
	source.consumed = false
}

func (source *ScriptSource) Close() error {
	err := source.file.Close()
	if err != nil {
		return err
	}
	return source.server.Close()
}

func (source *ScriptSource) GetID() (string, error) {
	if source.server == nil {
		return "", fmt.Errorf("source not initialized")
	}

	addr, err := source.server.GetAddr()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (source *ScriptSource) GetName() string {
	return source.name
}

func (source *ScriptSource) GetAttachments() []string {
	return []string{source.log.Path}
}

func (srcFile *ScriptFile) Open() (io.Reader, error) {
	sshSession, err := srcFile.source.server.NewSession()
	if err != nil {
		return nil, err
	}
	srcFile.session = sshSession
	out, err := sshSession.StdoutPipe()
	if err != nil {
		return nil, err
	}
	sshSession.Stderr = srcFile.source.log.Out
	err = sshSession.Start("cd " + path.Dir(srcFile.script) + "; ./" + path.Base(srcFile.script))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (srcFile *ScriptFile) Name() string {
	// custom name
	return srcFile.source.Name
}

func (srcFile *ScriptFile) Dir() string {
	// Same dir as script
	return path.Dir(srcFile.source.Path)
}

func (srcFile *ScriptFile) Path() string {
	return path.Join(srcFile.Dir(), srcFile.Name())
}

func (srcFile *ScriptFile) Close() error {
	return srcFile.session.Wait()
}
