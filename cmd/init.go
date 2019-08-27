package cmd

import (
	"fmt"
	"github.com/clcert/osr/databases"
	"github.com/spf13/cobra"
)

// Init creates the initial users for the database.
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init configuration and users of the DB OSR uses",
	Long:  "Init configuration and users of the DB OSR uses",
	Run: func(cmd *cobra.Command, args []string) {
		if err := databases.GetDBConfigData(); err != nil {
			fmt.Print(err)
		}
	},
}
