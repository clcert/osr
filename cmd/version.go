package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const VERSION = "3.7.0"

// CreateDB creates all the models databases.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Checks OSR version",
	Long:  "Checks OSR version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("OSR version %s\n", VERSION)
	},
}
