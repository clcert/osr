package geolite2

import (
	"encoding/csv"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strconv"
	"strings"
)

func parseSource(source sources.Source, saver savers.Saver, args *tasks.Args) {
	id, err := source.GetID()
	if err != nil {
		id = "unknown"
	}
	// prepareFTP the first one, if it's CSV, we scan all as CSV.
	// If it's ZIP, we assume that it is only one.
	entry := source.Next()
	if entry == nil {
		args.Log.WithFields(logrus.Fields{
			"source": id,
		}).Info("no files on this source. Skipping...")
		return
	}
	reader, err := entry.Open()
	if err != nil {
		args.Log.WithFields(logrus.Fields{
			"path":   entry.Path(),
			"source": id,
		}).Error("cannot open first file, skipping source")
		return
	}
	if strings.HasSuffix(entry.Name(), ".zip") {
		args.Log.WithFields(logrus.Fields{
			"path":   entry.Path(),
			"source": id,
		}).Info("entry is a zip. Assuming that there is only one entry on this source and continuing with zip...")
		zipConfig := &sources.ZipConfig{
			Path: entry.Path(),
			Filter: &sources.FilterConfig{
				Recursive: true,
			},
		}
		zipSource, err := zipConfig.NewFrom(source.GetName()+"_zip", reader, args.Params)
		if err != nil {
			// Log it and continue
			args.Log.WithFields(logrus.Fields{
				"entry": entry.Name(),
			}).Error("Cannot open file as zip")
			err := entry.Close()
			if err != nil {
				args.Log.WithFields(logrus.Fields{
					"entry": entry.Name(),
				}).Errorf("error closing entry: %s", err)
			}
			return
		}
		err = zipSource.Init()
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"path":   entry.Path(),
				"source": id,
			}).Error("could not init zip source")
		}
		zipID, err := zipSource.GetID()
		if err == nil {
			id = zipID
		}
		defer zipSource.Close()
		err = entry.Close()
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"entry": entry.Name(),
			}).Errorf("error closing entry: %s", err)
		}
		entry = zipSource.Next()
		if entry == nil {
			args.Log.WithFields(logrus.Fields{
				"source": id,
			}).Error("Empty zip")
			return
		}
		source = zipSource
	}
	for {
		saveFunc, ok := nameToFunc[entry.Name()]
		if ok {
			args.Log.WithFields(logrus.Fields{
				"entry": entry.Path(),
			}).Info("Reading entry...")

			if err := saveFunc(entry, saver, args); err != nil {
				args.Log.WithFields(logrus.Fields{
					"entry": entry.Path(),
					"error": err,
				}).Info("problems parsing this entry")
			}
			args.Log.WithFields(logrus.Fields{
				"entry": entry.Path(),
			}).Info("Done parsing this entry")
		} else {
			args.Log.WithFields(logrus.Fields{
				"entry": entry.Path(),
			}).Info("skipping entry...")
		}
		entry = source.Next()
		if entry == nil {
			break
		}
	}
	args.Log.WithFields(logrus.Fields{
		"id": id,
	}).Info("Done parsing this source")
}

func saveCountries(entry sources.Entry, saver savers.Saver, args *tasks.Args) error {
	reader, err := entry.Open()
	if err != nil {
		return err
	}
	defer entry.Close()
	countriesCSV := csv.NewReader(reader)
	countriesCSV.ReuseRecord = true
	// Discard first line
	if _, err := countriesCSV.Read(); err != nil {
		args.Log.Error("Err reading CSV file")
		return err
	}
	for {
		rec, err := countriesCSV.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			args.Log.Error("Err reading CSV file")
			return err
		}
		alpha2, name := rec[4], rec[5]
		if len(alpha2) == 0 { // A continent without country code
			// I hate the continents which think they're a country.
			alpha2, name = rec[2], rec[3]
			if alpha2 == "AS" {
				alpha2 = "AA" // We use AA as country code of AS continent.
			}
		}
		geoid, err := strconv.Atoi(rec[0])
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"geoid": rec[1],
				"error": err,
			}).Error("Error saving country")
			return err
		}
		if err := saver.Save(&models.Country{
			Alpha2:    alpha2,
			Name:      name,
			GeonameId: geoid,
		}); err != nil {
			args.Log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error saving country")
		}
	}
	return nil
}

func saveCountrySubnets(entry sources.Entry, saver savers.Saver, args *tasks.Args) error {
	reader, err := entry.Open()
	if err != nil {
		return err
	}
	defer entry.Close()
	countriesCSV := csv.NewReader(reader)
	countriesCSV.ReuseRecord = true
	// Discard first line (csv header)
	if _, err := countriesCSV.Read(); err != nil {
		args.Log.Error("Err reading CSV file")
		return err
	}
	for {
		rec, err := countriesCSV.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			args.Log.Error("Err reading CSV file")
			return err
		}
		_, subnet, err := net.ParseCIDR(rec[0])
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"cidr": rec[0],
			}).Error("Err Parsing Network Address")
			return err
		}
		if len(rec[1]) > 0 {
			geoid, err := strconv.Atoi(rec[1])
			if err != nil {
				args.Log.WithFields(logrus.Fields{
					"geoid": rec[1],
				}).Error("Err converting geoname Number")
				return err
			}
			if err := saver.Save(&models.SubnetCountry{
				TaskID:           args.GetTaskID(),
				SourceID:         args.GetSourceID(),
				Subnet:           subnet,
				CountryGeonameId: geoid,
			}); err != nil {
				args.Log.WithFields(logrus.Fields{
					"error": err,
				}).Error("Error saving subnet")
			}
		}
	}
	return nil
}

func saveASNSubnets(entry sources.Entry, saver savers.Saver, args *tasks.Args) error {
	reader, err := entry.Open()
	if err != nil {
		return err
	}
	defer entry.Close()
	countriesCSV := csv.NewReader(reader)
	countriesCSV.ReuseRecord = true
	// Discard first line
	if _, err := countriesCSV.Read(); err != nil {
		args.Log.Error("Error reading CSV file")
		return err
	}
	for {
		rec, err := countriesCSV.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			args.Log.Error("Error reading CSV file")
			return err
		}
		_, subnet, err := net.ParseCIDR(rec[0])
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"cidr": rec[0],
			}).Error("Err Parsing Network Address")
			return err
		}
		asn, err := strconv.Atoi(rec[1])
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"asnID": rec[1],
			}).Error("Err converting AS Number")
			return err
		}
		if err := saver.Save(&models.SubnetASN{
			TaskID:   args.GetTaskID(),
			SourceID: args.GetSourceID(),
			Subnet:   subnet,
			AsnID:    asn,
		}); err != nil {
			args.Log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error saving subnet")
		}
	}
	return nil
}
