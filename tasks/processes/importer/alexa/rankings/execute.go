package rankings

import (
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"strings"
)

// name of the file used in the data source.
const rankingFilename = "top-1m.csv"

// Imports chilean-domains information from Geolite2 Source
func Execute(args *tasks.Args) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	// I expect only one file: the csv
	entry := source.Next()
	if entry == nil {
		args.Log.Info("No more rankings left")
		return nil
	}
	rankReader, err := entry.Open()
	defer entry.Close()
	if err != nil {
		return err
	}
	var thisEntry sources.Entry
	if strings.HasSuffix(entry.Name(), ".zip") {
		zipConfig := &sources.ZipConfig{
			Path: entry.Path(),
			Filter: &sources.FilterConfig{
				Recursive: false,
				Patterns: []string{
					rankingFilename,
				},
			},
		}
		zipSource, err := zipConfig.NewFrom(source.GetName()+"_zip", rankReader, args.Params)
		if err != nil {
			// Log it and continue
			args.Log.WithFields(logrus.Fields{
				"entry": entry.Name(),
			}).Error("Cannot open file as zip")
			return err
		}
		err = zipSource.Init()
		if err != nil {
			return err
		}
		defer zipSource.Close()
		// I expect only one file: the csv
		zipEntry := zipSource.Next()
		if zipEntry == nil {
			args.Log.Info("No more rankings left")
			return nil
		}
		thisEntry = zipEntry
	} else {
		thisEntry = entry
	}
	err = saveCSV(thisEntry, saver, args)
	if err != nil {
		return err
	}
	args.Log.Info("Done parsing file, exiting...")
	return nil
}
