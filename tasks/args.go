package tasks

import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/utils"
)

// Args defines a set of args
type Args struct {
	Process *Process         // A reference to the Process struct which called the Task Method.
	Task    *Task            // A reference to the task context.
	Sources []sources.Source // A list of sources with the data to import.
	Savers  []savers.Saver   // A list of savers used to store the processed data
	Params  utils.Params     // A list of args
	Log     *logs.OSRLog     // Log used in the process(es) execution
}

// AddSources adds more sources to an args list
func (args *Args) AddSources(sourcesList []sources.Source) {
	for _, v := range sourcesList {
		args.Sources = append(args.Sources, v)
	}
}

// AddSavers add new savers to the savers list
func (args *Args) AddSavers(saversList []savers.Saver) {
	for _, v := range saversList {
		args.Savers = append(args.Savers, v)
	}
}

func (args *Args) GetSourceID() models.DataSourceID {
	setSourceID := args.Task.GetConfig(args.Process.Command).SourceID
	if  setSourceID > 0 {
		return setSourceID
	}
	return args.Process.DefaultSourceID
}

func (args *Args) GetTaskID() int {
	return args.Task.TaskSession.ID
}
