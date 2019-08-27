package cmd

import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/panics"
	"github.com/spf13/cobra"
)

func init() {
	ModelsCmd.AddCommand(CreateDBCommand)
}

// Models command groups all importer related to models.
var ModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Shows models registered",
	Long:  "Shows models registered",
}

// CreateDB creates all the models databases.
var CreateDBCommand = &cobra.Command{
	Use:   "createdb",
	Short: "Creates the databases",
	Long:  "Creates the databases",
	Run: func(cmd *cobra.Command, args []string) {
		err := models.DefaultModels.CreateTables()
		if err != nil {
			panic(&panics.Info{
				Text:        "couldn't create the tables of the databases",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log},
			})
		}

	},
}
