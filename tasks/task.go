package tasks

import (
	"fmt"
	"path"
	"strings"

	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/utils"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Defines a complete file with one or more process, savers and exports.
type TaskConfig struct {
	Name         string              // Name of the task
	Description  string              // Description of the task
	AbortOnError bool                // If true, task aborts if a process throws an error
	Incognito    bool                // If true, task is not registered and taskID is assigned to 0
	Parallel     bool                // Tasks are executed in parallel (?)
	Params       utils.Params        // A list of global parameters
	Processes    []*ProcessConfig    // A list of config for processes.
}

// A task defines the state of execution of a TaskConfig. It contains the stats of the execution.
type Task struct {
	*TaskConfig
	TaskSession *models.Task     // The task model used on this session
	Succeeded   []string         // A list with succeeded processes
	Failed      map[string]error // A list with failed processes
	Attachments []string         // A list with attachments created by processes
	DB          *pg.DB           // A pointer to a DB writer.
	CmdParams   utils.Params     // Params received by command line. They have the highest preference.
}

// GetSucceeded formats the names of the succeeded process related to the tasks.
func (task *Task) GetSucceeded() string {
	return strings.Join(task.Succeeded, ", ")
}

// GetFailed formats the names of the failed process related to the tasks.
func (task *Task) GetFailed() string {
	failed := make([]string, len(task.Failed))
	i := 0
	for process, _ := range task.Failed {
		failed[i] = process
		i++
	}
	return strings.Join(failed, ", ")
}

// Returns a new Task based on a task config.
func New(config *TaskConfig, params []string) (newTask *Task, err error) {
	logs.Log.Info("Initializing Database Connection...")
	dbHandler, err := databases.GetPostgresWriter()
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Couldn't connect to database")
		return
	}
	currentTask, err := models.NewTaskSession(dbHandler, !config.Incognito)
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Couldn't initialize imports: cannot write to database")
		return
	}
	cmdParams := utils.ListToParams(params)
	newTask = &Task{
		TaskConfig:  config,
		TaskSession: currentTask,
		DB:          dbHandler,
		Succeeded:   make([]string, 0),
		Failed:      make(map[string]error, 0),
		Attachments: make([]string, 0),
		CmdParams:   cmdParams,
	}
	logs.Log.WithFields(logrus.Fields{
		"importer":  newTask.TaskSession.ID,
		"incognito": config.Incognito,
	}).Info("Created new task session!")
	return
}

// ExecuteAll executes a entire task.
func (task *Task) Execute() {
	defer notify(task)
	processes := task.GetCommands()
	for index, process := range processes {
		err := task.execute(process, index)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"task": task.Name,
				"index": index,
				"process": process,
			}).Errorf("Task failed: %s", err)
			task.AddFailed(process, err)
			if task.AbortOnError {
				task.TaskSession.Failed()
				return
			}
		} else {
			task.AddSucceeded(process)
		}
	}
	logs.Log.WithFields(logrus.Fields{
		"task_id":   task.TaskSession.ID,
		"succeeded": task.GetSucceeded(),
		"failed":    task.GetFailed(),
	}).Info("All importer functions executed")
	task.TaskSession.Succeeded()
	if !task.Incognito {
		err := task.TaskSession.Save(task.DB)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Couldn't save import as succeeded in database")
		}
	} else {
		logs.Log.WithFields(logrus.Fields{
			"task_id":   task.TaskSession.ID,
			"succeeded": task.GetSucceeded(),
			"failed":    task.GetFailed(),
		}).Info("Task was not saved because incognito mode is on")
	}
	return
}

// Export executes a specific process name in a task
func (task *Task) execute(processName string, processIndex int) error {
	config := task.GetConfig(processIndex)
	if config == nil {
		logs.Log.WithFields(logrus.Fields{
			"command": processName,
			"index": processIndex,
		}).Error("Process not found on task")
		return fmt.Errorf("process config not found on task: %s", processName)
	}

	process, ok := Registered[processName]
	if !ok {
		logs.Log.WithFields(logrus.Fields{
			"command": processName,
			"index": processIndex,
		}).Error("Process not defined in system")
		return fmt.Errorf("process not defined in system: %s", processName)
	}
	if process.Execute == nil {
		return fmt.Errorf("proccess command is not defined")
	}

	// Add args
	args, err := process.NewArgs(task, processIndex)
	if err != nil {
		return err
	}

	// Add process specific params
	args.Params = args.Params.Join(config.Params)
	// Add command line specific params
	args.Params = args.Params.Join(task.CmdParams)

	// parse and initialize sourcesList
	if process.NumSources >= 0 && len(config.Sources) != process.NumSources {
		return fmt.Errorf("there should be only %d sources(s) with this process and there are %d", process.NumSources, len(config.Sources))
	}
	sourcesList := make([]sources.Source, 0)
	for i, sourceConf := range config.Sources {
		sourceName := fmt.Sprintf("%s_%s_source-%d", task.GetSafeName(), process.GetSafeName(), i)
		source, err := sourceConf.New(sourceName, args.Params)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"command": processName,
				"index": processIndex,
			}).Error("source list error on index %d: %s", i, err)
			return fmt.Errorf("source list error on index %d: %s", i, err)
		}
		sourcesList = append(sourcesList, source)
		if err := source.Init(); err != nil {
			return fmt.Errorf("error initializing source: %s", err)
		}
		task.AddAttachments(source)
	}
	args.AddSources(sourcesList)

	defer func() {
		for _, source := range args.Sources {
			_ = source.Close()
		}
	}()

	// parse and initialize saversList
	if process.NumSavers >= 0 && len(config.Savers) != process.NumSavers {
		return fmt.Errorf("there should be only %d saver(s) with this process and there are %d", process.NumSavers, len(config.Savers))
	}
	saversList := make([]savers.Saver, 0)
	for i, saverConf := range config.Savers {
		saverName := fmt.Sprintf("%s_%s_saver_%d", task.GetSafeName(), process.GetSafeName(), i)
		saver, err := saverConf.New(saverName, args.Params)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"command": processName,
				"index": processIndex,
			}).Errorf("saver list error on index %d: %s", i, err)
			return fmt.Errorf("saver list error on index %d: %s", i, err)
		}
		saversList = append(saversList, saver)
		if err := saver.Start(); err != nil {
			logs.Log.WithFields(logrus.Fields{
				"command": processName,
				"index": processIndex,
			}).Errorf("Cannot start saver with index %d: %s", i, err)
			return err
		}
		task.AddAttachments(saver)
	}
	args.AddSavers(saversList)

	defer func() {
		for _, saver := range args.Savers {
			_ = saver.Finish()
		}
	}()

	logs.Log.WithFields(logrus.Fields{
		"command": processName,
		"index": processIndex,
	}).Info("executing process")
	err = process.Execute(args)
	if err == nil {
		logs.Log.WithFields(logrus.Fields{
			"command": processName,
			"index": processIndex,
		}).Info("process succeeded")
	} else {
		logs.Log.WithFields(logrus.Fields{
			"command": processName,
			"index": processIndex,
		}).Error("process failed")
	}
	return err
}

// Close closes the task database.
func (task *Task) Close() error {
	return task.DB.Close()
}

// HasErrors returns true if the task had any errors on its execution.
func (task *Task) HasErrors() bool {
	return len(task.Failed) > 0
}

// AddSucceeded adds a name to succeeded list
func (task *Task) AddSucceeded(name string) {
	task.Succeeded = append(task.Succeeded, name)
}

// AddFailed adds an error to failed list
func (task *Task) AddFailed(name string, err error) {
	task.Failed[name] = err
}

// GetConfig returns the configuration with an specific name in the task.
func (task *Task) GetConfig(index int) *ProcessConfig {
	if index < len(task.TaskConfig.Processes) {
		return task.TaskConfig.Processes[index]
	}
	return nil
}

// GetCommands returns a list with the names of the commands defined in the task.
func (task *Task) GetCommands() []string {
	commands := make([]string, len(task.Processes))
	for i, process := range task.Processes {
		commands[i] = process.Command
	}
	return commands
}

// AddAttachments adds an attachment to the task state.
func (task *Task) AddAttachments(attachable mailer.Attachable) {
	for _, attachment := range attachable.GetAttachments() {
		task.Attachments = append(task.Attachments, attachment)
	}
}

func (task *Task) GetAttachments() []string {
	return task.Attachments
}

// ParseConfig parses the config related to the task from a specific file.
func ParseConfig(name string) (*TaskConfig, error) {
	viperConfig := viper.New()
	taskRoot, err := GetTasksPath()
	if err != nil {
		return nil, err
	}
	viperConfig.SetConfigFile(path.Join(taskRoot, name))
	if err := viperConfig.ReadInConfig(); err != nil {
		return nil, err
	}
	var config TaskConfig

	if err := viperConfig.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// GetSafeName returns the safe name for the task.
func (task *Task) GetSafeName() string {
	return fmt.Sprintf("task-%d", task.TaskSession.ID)
}
