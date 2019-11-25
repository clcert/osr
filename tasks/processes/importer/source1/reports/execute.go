package reports

import (
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/filters"
	"github.com/sirupsen/logrus"
)

func Execute(args *tasks.Context) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	conf, errs := filters.NewDateConfig(args.Params)
	for err := range errs {
		args.Log.Errorf("Error parsing config: %s", err)
	}
	for {
		entry := source.Next()
		if entry == nil {
			break
		}
		if err := saveReport(entry, saver, args, conf); err != nil {
			args.Log.WithFields(logrus.Fields{
				"file": entry.Path(),
			}).Errorf("Error reading file: %s", err)
			continue
		}
	}
	return nil
}
