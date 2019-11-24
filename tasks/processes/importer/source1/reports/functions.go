package reports

import (
	"encoding/csv"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/filters"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

var scanToReport = map[string]func(string) (models.ReportTypeID, map[string]string, error){
	"bot":           newBot,
	"bots":          newBots,
	"bruteforce":    newBruteForce,
	"controller":    newC2,
	"darknet":       newDarknet,
	"honeypot":      newHoneypot,
	"openresolvers": newDNSResolver,
	"phishing":      newPhishing,
	"proxy":         newProxy,
	"spam":          newSpam,
}

const timeLayout = "2006-01-02 15:04:05"

func parseLine(line []string, args *tasks.Args, conf *filters.DateConfig) (*models.ReportEntry, error) {
	if len(line) < 6 {
		return nil, fmt.Errorf("line badly formatted")
	}
	scanFunction, ok := scanToReport[line[0]]
	if !ok {
		return nil, fmt.Errorf("cannot parse line")
	}
	ip := net.ParseIP(line[1])
	if ip == nil {
		return nil, fmt.Errorf("cannot parse IP")
	}
	date, err := time.Parse(timeLayout, line[3])
	if err != nil {
		return nil, err
	}
	if !conf.IsDateInRange(date) {
		return nil, fmt.Errorf("line is not in date range")
	}
	typeID, props, err := scanFunction(line[4])
	if err != nil {
		return nil, err
	}
	if len(props) == 0 {
		// we insert a dummy property
		props[""] = ""
	}
	entry := &models.ReportEntry{
		SourceID:     args.Process.Source,
		TaskID:       args.Task.ID,
		ReportTypeID: typeID,
		IP:           ip,
		Date:         date,
		Properties:   props,
	}
	return entry, nil
}

func saveReport(entry sources.Entry, saver savers.Saver, args *tasks.Args, conf *filters.DateConfig) error {
	reader, err := entry.Open()
	defer entry.Close()
	if err != nil {
		return err
	}
	csvReader := csv.NewReader(reader)
	csvReader.ReuseRecord = true
	csvReader.Comment = '#'
	csvReader.Comma = '|'
	args.Log.
		WithFields(logrus.Fields{
			"file": entry.Path(),
		}).
		Info("Reading file...")
	for {
		line, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			// Problems with this line, we skip it.
			return err
		}
		report, err := parseLine(line, args, conf)
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
