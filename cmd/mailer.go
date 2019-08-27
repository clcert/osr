package cmd

import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/clcert/osr/panics"
	"github.com/spf13/cobra"
)

func init() {
	MailerCmd.AddCommand(MailerTestCmd)
}

// Mailer command groups mailer-related subcommands
var MailerCmd = &cobra.Command{
	Use:   "mailer",
	Short: "executes raw queries from YAML files",
	Long:  "executes raw queries from YAML files",
}

// Test command allows to test email notification system.
var MailerTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Sends a test mail",
	Long:  "Sends a test mail",
	Run: func(cmd *cobra.Command, args []string) {
		err := mailer.SendTestMail()
		if err != nil {
			panic(&panics.Info{
				Text:        "couldn't send a test mail",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log},
			})
		}
		logs.Log.Info("mail sent!")
	},
}
