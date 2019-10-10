package grabber_protocol_scan

import (
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

// This function reads a remote server and savers all the
// results contained in folders with a signal empty file.
// []struct -> error
func Execute(args *tasks.Args) (err error) {
	source := args.Sources[0]
	srcAddr, err := source.GetID()
	if err != nil {
		return err
	}
	srcIPStr := strings.Split(srcAddr, ":")[0]
	srcIP := net.ParseIP(srcIPStr)
	filesRead := 0
	for {
		file := source.Next()
		if file == nil {
			break
		}
		args.Log.WithFields(logrus.Fields{
			"files_read": filesRead,
		}).Info("Reading files...")
		err := parseFile(file, args, srcIP)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
			}).Errorf("Error reading file: %s", err)
		}
		filesRead++
	}
	return
}
