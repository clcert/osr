package protocol_scan

import (
	"github.com/clcert/osr/tasks"
)

// This function reads a remote server and savers all the
// results contained in folders with a signal empty file.
// []struct -> error
func Execute(args *tasks.Args) (err error) {
	source := args.Sources[0]
	saver := args.Savers[0]
	err = parseFiles(source, saver, args)

	return
}
