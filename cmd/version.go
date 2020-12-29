package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "4.0.3"

var Build string

// CreateDB creates all the models databases.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Checks OSR version",
	Long:  "Checks OSR version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("OSR version %s build %s\n", VERSION, Build)
	},
}
