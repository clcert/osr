package reports

import (
	"fmt"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/clcert/osr/utils/filters"
	"github.com/sirupsen/logrus"
	"io"
)

func saveReport(entry sources.Entry, saver savers.Saver, args *tasks.Args, conf *filters.DateConfig) error {
	reader, err := entry.Open()
	defer entry.Close()
	if err != nil {
		return err
	}
	if !fileToReport.Exists(entry.Name()) {
		return fmt.Errorf("file not supported for importing")
	}
	headedCSV, err := utils.NewHeadedCSV(reader, nil)
	if err != nil {
		return err
	}
	args.Log.
		WithFields(logrus.Fields{
			"file": entry.Path(),
		}).
		Info("Reading file...")
	for {
		line, err := headedCSV.NextRow()
		if err != nil {
			if err == io.EOF {
				break
			}
			// Problems with this line, we skip it.
			continue
		}
		report, err := fileToReport.Parse(entry.Name(), line, args, conf)
		if err != nil {
			args.Log.
				WithFields(logrus.Fields{
					"file": entry.Path(),
					"line": line,
				}).
				Errorf("cannot parse line: %s", err)
			continue
		}
		if err := saver.Save(report); err != nil {
			args.Log.
				WithFields(logrus.Fields{
					"file":   entry.Path(),
					"report": report,
				}).
				Errorf("cannot save report: %s", err)
			continue
		}
	}
	return nil
}
