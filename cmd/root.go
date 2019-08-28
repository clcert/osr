// The package cmd defines the importer for OSR, using Cobra Commander library.
package cmd

import (
	"fmt"
	"github.com/clcert/osr/cmd/plot"
	"github.com/clcert/osr/panics"
	"github.com/spf13/cobra"
	"os"
)

const VERSION = "2.1.0"

func init() {
	RootCmd.AddCommand(ModelsCmd)
	RootCmd.AddCommand(RawQueryCmd)
	RootCmd.AddCommand(RemoteCmd)
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(MailerCmd)
	RootCmd.AddCommand(TaskCmd)
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(plot.PlotCmd)
}

// The root of osr.
var RootCmd = &cobra.Command{
	Use:          "osr",
	Short:        "OSR manages data and statistics related to chilean network scans",
	Long:         "OSR manages data and statistics related to chilean network scans",
	SilenceUsage: true,
}

func Execute() {
	defer panics.NotifyPanic()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
