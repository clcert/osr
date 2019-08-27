package tasks

import (
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/utils"
	"strings"
)

type Processes map[string]*Process

// The default list of registered process
var Registered = make(Processes)

const InfiniteSources = -1

// ProcessConfig defines the configuration specific for a process
type ProcessConfig struct {
	Command string           // Name of the command
	Sources []sources.Config // List of sources related to the command
	Savers  []savers.Config  // List of savers related to the command
	Params  utils.Params     // List of specific params. They override the params of the global file.
}

// Process defines completely a Task.
// Task importer have now a [provider]->[data] structure, but this
// implementation doesn't demand it.
type Process struct {
	Name        string              // Readable name for the importer command
	Command     string              // The name you write when you want to execute the command
	Description string              // A description for the importer routine
	URL         string              // If exists, a URL related with the source of the data
	Source      models.DataSourceID // Provider ID. Check Providers model to get it or register it.
	Execute     func(*Args) error   // An action to be executed when this command is called.
	NumSources  int                 // Number of allowed sources on this task. If negative, it's unlimited.
	NumSavers   int                 // Number of allowed savers on this task. If negative, it's unlimited.
}

// Registers a process to the global dictionary.
func (processes Processes) Register(processList ...*Process) {
	for _, process := range processList {
		processes[process.Command] = process
	}
}

// Creates a new object of type Args, used in process execution.
// It also creates a logger for the process and adds the task params to its
// params list.
func (process *Process) NewArgs(task *Task) (*Args, error) {
	log, err := logs.NewLog(fmt.Sprintf("%s_%s", task.GetSafeName(), process.GetSafeName()))
	if err != nil {
		return nil, err
	}
	// Adding log to task attachments for future notifications
	task.AddAttachments(log)
	args := &Args{
		Process: process,
		Sources: make([]sources.Source, 0),
		Savers:  make([]savers.Saver, 0),
		Params:  make(map[string]string, 0),
		Task:    task.TaskSession,
		Log:     log,
	}
	// Adding params of task
	args.AddParams(task.Params)
	return args, nil
}

func (process *Process) GetSafeName() string {
	return strings.Replace(process.Command, "/", "-", -1)
}