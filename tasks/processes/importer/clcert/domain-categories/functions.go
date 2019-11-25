package domain_categories

import (
	"encoding/csv"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"io"
)

func ImportCategories(source sources.Source, saver savers.Saver, args *tasks.Context) error {
	file := source.Next()
	if file == nil {
		return fmt.Errorf("source empty")
	}
	reader, err := file.Open()
	if err != nil {
		return err
	}
	csvReader := csv.NewReader(reader)
	csvReader.ReuseRecord = true
	args.Log.
		WithFields(logrus.Fields{
			"file": file.Path(),
		}).
		Info("Reading file...")
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if pErr, ok := err.(*csv.ParseError); ok {
			args.Log.WithFields(logrus.Fields{
				"start_line": pErr.StartLine,
				"line":       pErr.Line,
				"column":     pErr.Column,
			}).Error("Couldn't read from body of file.")
			return err
		}

		if len(line) < 4 {
			args.Log.
				WithFields(logrus.Fields{
					"num_fields": len(line),
					"line":       line,
				}).Error("Missing fields in row.")
			continue
		}
		subdomain, domain, tld, category := line[0], line[1], line[2], line[3]
		err = saver.Save(&models.DomainToCategory{
			TaskID:             args.GetTaskID(),
			SourceID:           args.GetSourceID(),
			DomainSubdomain:    subdomain,
			DomainName:         domain,
			DomainTLD:          tld,
			DomainCategorySlug: category,
		})
		if err != nil {
			//TODO Log this
		}
	}
	return nil
}
