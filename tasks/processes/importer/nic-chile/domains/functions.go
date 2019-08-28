package domains

import (
	"bufio"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func processSource(source sources.Source, saver savers.Saver, args *tasks.Args) error {
	for {
		var err error
		entry := source.Next()
		if entry == nil {
			return nil
		}
		lowerPath := strings.ToLower(entry.Path())
		if strings.Contains(lowerPath, "ultimos") || strings.Contains(lowerPath, "nuevos") {
			err = saveNewDomains(entry, saver, args)
			saver.SendMessage(map[string]string{"close": entry.Path()})
		} else if strings.Contains(strings.ToLower(entry.Path()), "eliminados") {
			err = saveDeletedDomains(entry, saver, args)
			saver.SendMessage(map[string]string{"close": entry.Path()})
		} else {
			args.Log.
				WithFields(logrus.Fields{
					"file": entry.Path(),
				}).
				Info("skipping file...")
		}
		if err != nil {
			args.Log.
				WithFields(logrus.Fields{
					"file": entry.Path(),
					"err":  err,
				}).
				Info("couldn't save domains")
		}
	}
}

func saveNewDomains(entry sources.Entry, saver savers.Saver, args *tasks.Args) error {
	reader, err := entry.Open()
	if err != nil {
		return err
	}
	defer entry.Close()
	scanner := bufio.NewScanner(reader)
	args.Log.
		WithFields(logrus.Fields{
			"file": entry.Path(),
		}).
		Info("Scanning file...")
	var firstRowRead bool
	for scanner.Scan() {
		if !firstRowRead {
			firstRowRead = true
			continue
		}
		line := scanner.Text()
		splitLine := strings.Split(line, ",")
		if len(splitLine) == 1 {
			args.Log.WithFields(logrus.Fields{
				"line": line,
				"file": entry.Path(),
			}).Error("bad formatted line")
			continue
		}
		_, domain, tld, err := utils.SplitDomain(splitLine[0])
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"line": splitLine[0],
				"file": entry.Path(),
			}).Error("error splitting domain")
			continue
		}
		splitDate := strings.Split(splitLine[1], ".")[0]
		date, err := time.Parse(NicTimeLayout, splitDate)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"date": splitDate,
				"line": line,
				"file": entry.Path(),
			}).Error("Couldn't parse date from file")
			continue
		}
		aDomain := &models.Domain{
			TaskID:           args.Task.ID,
			SourceID:         args.Process.Source,
			Subdomain:        "",
			Name:             domain,
			TLD:              tld,
			RegistrationDate: date,
		}
		if err := saver.Save(savers.Savable{
			Object: aDomain,
			Meta:   map[string]string{"outID": entry.Path()},
		}); err != nil {
			args.Log.WithFields(logrus.Fields{
				"file": entry.Path(),
			}).Error("Couldn't save object")
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		args.Log.Error("Couldn't read from body of file.")
		return err
	}
	return nil
}

func saveDeletedDomains(entry sources.Entry, saver savers.Saver, args *tasks.Args) error {
	reader, err := entry.Open()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(reader)
	var date time.Time
	args.Log.
		WithFields(logrus.Fields{
			"file": entry.Path(),
		}).
		Info("Scanning file...")
	numEntries := 0
	for scanner.Scan() {
		if date.IsZero() {
			dateString := strings.Split(scanner.Text(), ": ")[1]
			date, err = time.Parse(NicTimeLayout, dateString)
			if err != nil {
				args.Log.WithFields(logrus.Fields{
					"date": dateString,
				}).Error("Couldn't parse date from header of file.")
				return err
			}
			continue
		}
		_, domain, tld, err := utils.SplitDomain(scanner.Text())
		if err != nil {
			return err
		}
		aDomain := &models.Domain{
			TaskID:       args.Task.ID,
			SourceID:     args.Process.Source,
			Subdomain:    "",
			Name:         domain,
			TLD:          tld,
			DeletionDate: date,
		}
		if err := saver.Save(savers.Savable{
			Object: aDomain,
			Meta:   map[string]string{"outID": entry.Path()},
		}); err != nil {
			return err
		}
		numEntries++
	}
	if err := scanner.Err(); err != nil {
		args.Log.Error("Couldn't read from body of file.")
		return err

	}
	return nil
}
