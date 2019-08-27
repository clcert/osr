package cmd

import (
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/clcert/osr/panics"
	"github.com/clcert/osr/remote"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	RemoteCmd.AddCommand(RemoteInitCmd)
	RemoteCmd.AddCommand(ListCmd)
	RemoteCmd.AddCommand(SendCmd)
	RemoteCmd.AddCommand(LocalExecCmd)
	RemoteCmd.AddCommand(RemoteExecCmd)
}

// Remote manages the scanner servers and checks their health.
var RemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Remote Servers Operations",
	Long:  "Remote Servers Operations",
}

// Init command connects to scanner servers vía password auth, and creates
// SSH keys to connect to them without the passwords.
var RemoteInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates the keys for the SSH connections",
	Long:  "Creates the keys for the SSH connections",
	Run: func(cmd *cobra.Command, args []string) {
		err := remote.CreateKeysForRemote()
		if err != nil {
			panic(&panics.Info{
				Text:        "couldn't create keys for remote",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log},
			})
		}
	},
}

// List command shows the health of all (or some) scanner servers.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List status of servers",
	Long:  "List status of servers",
	Run: func(cmd *cobra.Command, args []string) {
		servers, err := remote.GetServers(args...)
		if err != nil {
			panic(&panics.Info{
				Text:        "error getting server list",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log},
			})
		}
		if err := remote.ListServersInfo(servers); err != nil {
			attachments := make([]mailer.Attachable, len(servers)+1)
			for i, server := range servers {
				attachments[i] = server
			}
			attachments[len(servers)] = logs.Log
			panic(&panics.Info{
				Text:        "error listing servers info",
				Err:         err,
				Attachments: attachments,
			})
		}
	},
}

// send command sends files to a remote server.
var SendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends files to a remote server.",
	Long:  "Sends files to a remote server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := remote.GetServer(args[0])
		if err != nil {
			panic(&panics.Info{
				Text:        "cannot get server list",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log, server},
			})
		}
		// TODO: add params?
		_, err = server.SendFiles(nil, remote.ResourcesPath, args[1:]...)
		if err != nil {
			panic(&panics.Info{
				Text:        fmt.Sprintf("cannot send files to server: %s", strings.Join(args[1:], ", ")),
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log, server},
			})
		}
	},
}

// LocalExec command sends a file to a remote server and executes it.
var LocalExecCmd = &cobra.Command{
	Use:   "local-exec",
	Short: "Sends a local script in 'scripts' folder to a remote server and then executes it.",
	Long:  "Sends a local script in 'scripts' folder to a remote server and then executes it",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := remote.GetServer(args[0])
		if err != nil {
			panic(&panics.Info{
				Text:        "cannot get server list",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log, server},
			})
		}
		err = server.Connect()
		if err != nil {
			panic(&panics.Info{
				Text:        "cannot connect to server",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log, server},
			})
		}
		defer server.Close()
		// Remove command name
		args = args[1:]
		for _, script := range args {
			logs.Log.Info("executing script")
			// TODO: use params someday?
			err = server.ExecuteLocalScript(remote.TempPath, script, nil)
			if err != nil {
				panic(&panics.Info{
					Text:        "local execution failed",
					Err:         err,
					Attachments: []mailer.Attachable{logs.Log, server},
				})
			}
		}
	},
}

// exec command executes a command in a remote server.
var RemoteExecCmd = &cobra.Command{
	Use:   "remote-exec",
	Short: "Sends local files to a remote server and then executes them.",
	Long:  "Sends local files to a remote server and then executes them",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := remote.GetServer(args[0])
		if err != nil {
			panic(&panics.Info{
				Text:        "cannot get server list",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log, server},
			})
		}
		err = server.Connect()
		if err != nil {
			panic(&panics.Info{
				Text:        "cannot connect to server",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log, server},
			})
		}
		defer server.Close()
		err = server.RemoteExecute(strings.Join(args[1:], " "))
		if err != nil {
			panic(&panics.Info{
				Text:        "remote execution failed",
				Err:         err,
				Attachments: []mailer.Attachable{logs.Log, server},
			})
		}
	},
}
