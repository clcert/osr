package cmd

import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/clcert/osr/panics"
	"github.com/clcert/osr/query"
	"github.com/spf13/cobra"
)

var inFiles []string
var outFolder string
var queries []string
var headers bool

func init() {
	RawQueryCmd.Flags().StringSliceVarP(&inFiles, "input-files", "i", []string{}, "Input folder absolute path.")
	RawQueryCmd.Flags().StringVarP(&outFolder, "output-folder", "o", "./", "Output folder absolute path.")
	RawQueryCmd.Flags().StringSliceVarP(&queries, "queries", "q", []string{}, "whitelisted query names. Default is all.")
	
	RawQueryCmd.Flags().BoolVarP(&headers, "headers", "H", false, "Output CSV with column headers.")
	_ = RawQueryCmd.MarkFlagRequired("input-files")
	_ = RawQueryCmd.MarkFlagRequired("output-folder")
}

// Models command groups all importer related to models.
var RawQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "executes PGSql queries from YAML files",
	Long:  "executes PGSql queries from YAML files",
	Run: func(cmd *cobra.Command, args []string) {
		err := query.Execute(inFiles, outFolder, queries, headers)
		if err != nil {
			panic(&panics.Info{
				Text:        "fatal error executing the raw query",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log},
			})
		}
	},
}
