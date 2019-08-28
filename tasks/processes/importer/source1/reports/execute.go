package reports

import (
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
)

func Execute(args *tasks.Args) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	for {
		entry := source.Next()
		if entry == nil {
			break
		}
		if err := saveReport(entry, saver, args); err != nil {
			args.Log.WithFields(logrus.Fields{
				"file": entry.Path(),
			}).Errorf("Error reading file: %s", err)
			continue
		}
	}
	return nil
}
