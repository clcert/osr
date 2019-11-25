package simple

import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

func process(entryChan chan sources.Entry, saver savers.Saver, wg *sync.WaitGroup, args *tasks.Context) {
	defer wg.Done()
	for file := range entryChan {
		logs.Log.WithFields(logrus.Fields{
			"file": file.Path(),
		}).Info("opening file")
		reader, err := file.Open()
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"file": file.Path(),
			}).Error("couldn't open file: %s", err)
			continue
		}
		csv, err := utils.NewHeadedCSV(reader, nil)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"file": file.Path(),
			}).Error("couldn't use file as csv: %s", err)
			continue
		}

		for {
			line, err := csv.NextRow()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					logs.Log.WithFields(logrus.Fields{
						"file": file.Path(),
						"line": line,
					}).Error("cannot read next line: %s. Skipping file...", err)
					break
				}
			}
			if err := saver.Save(savers.Savable{
				Object: line,
				Meta: map[string]string{
					"outID": file.Name(),
				},
			}); err != nil {
				logs.Log.WithFields(logrus.Fields{
					"file": file.Path(),
					"line": line,
				}).Error("cannot save line: %s", err)
			}
		}
	}
	logs.Log.WithFields(logrus.Fields{
	}).Info("thread done!")
}
