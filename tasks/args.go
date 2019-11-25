package tasks

import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/utils"
)

// Context defines a set of args
type Context struct {
	Process *Process // A reference to the Process struct which called the Task Method.
	Index   int
	Task    *Task            // A reference to the task context.
	Sources []sources.Source // A list of sources with the data to import.
	Savers  []savers.Saver   // A list of savers used to store the processed data
	Params  utils.Params     // A list of args
	Log     *logs.OSRLog     // Log used in the process(es) execution
}

// AddSources adds more sources to an args list
func (context *Context) AddSources(sourcesList []sources.Source) {
	for _, v := range sourcesList {
		context.Sources = append(context.Sources, v)
	}
}

// AddSavers add new savers to the savers list
func (context *Context) AddSavers(saversList []savers.Saver) {
	for _, v := range saversList {
		context.Savers = append(context.Savers, v)
	}
}

func (context *Context) GetProcessConfig() *ProcessConfig {
	return context.Task.GetConfig(context.Index)
}

func (context *Context) GetSourceID() models.DataSourceID {
	conf := context.GetProcessConfig()
	if conf != nil && conf.SourceID > 0 {
		return conf.SourceID
	}
	return context.Process.DefaultSourceID
}

func (context *Context) GetTaskID() int {
	return context.Task.TaskSession.ID
}
