package cmd

import (
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/clcert/osr/panics"
	"github.com/clcert/osr/tasks"
	_ "github.com/clcert/osr/tasks/registered"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var params []string

func init() {
	TaskCmd.Flags().StringSliceVarP(&params, "params", "p", []string{}, "Parameters")

}

// Process executes a batch of process defined in a conf file.
var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Executes a process task",
	Long:  "Init configuration and users of the DB OSR uses",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no task file in args")
		}
		for _, configName := range args {
			config, err := tasks.ParseConfig(configName)
			if err != nil {
				logs.Log.WithFields(logrus.Fields{
					"name": configName,
				}).Infof("file parsing error: %s", err)
				panic(panics.Info{
					Text:        fmt.Sprintf("File parsing error: %s", configName),
					Err:         err,
					Attachments: []mailer.Attachable{logs.Log},
				})
			}
			logs.Log.WithFields(logrus.Fields{
				"name":          config.Name,
				"description":   config.Description,
				"global_params": config.Params,
				"file":          configName,
			}).Info("executing task file")

			task, err := tasks.New(config, params)
			if err != nil {
				panic(&panics.Info{
					Text:        fmt.Sprintf("error executing task with name %s", configName),
					Err:         err,
					Attachments: []mailer.Attachable{logs.Log},
				})
			}
			task.Execute()
		}
		return nil
	},
}
