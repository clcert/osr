package rankings

import (
	"encoding/csv"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

func saveCSV(entry sources.Entry, saver savers.Saver, args *tasks.Args) error {
	var tlds []string
	if _, ok := args.Params["tlds"]; ok {
		tlds = strings.Split(args.Params["tlds"], ",")
	} else {
		tlds = []string{"cl"}
	}
	args.Log.WithFields(logrus.Fields{
		"fields": tlds,
	}).Info("TLDs to consider")
	reader, err := entry.Open()
	if err != nil {
		return err
	}
	defer entry.Close()
	args.Log.Info("Reading file...")
	countriesCSV := csv.NewReader(reader)
	countriesCSV.ReuseRecord = true
	rank := 1
	for {
		rec, err := countriesCSV.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		subdomain, name, tld, err := utils.SplitDomain(rec[1])
		if err != nil {
			args.Log.Error("Error splitting domain %s: %v", rec[1], err)
			continue
		}
		for _, allowedTLD := range tlds {
			if strings.TrimSpace(tld) == allowedTLD {
				err := saver.Save(&models.DomainRanking{
					TaskID:          args.GetTaskID(),
					SourceID:        args.GetSourceID(),
					Ranking:         int64(rank),
					DomainSubdomain: subdomain,
					DomainName:      name,
					DomainTLD:       tld,
				})
				if err != nil {
					args.Log.Error("Error saving domain %s: %v", rec[1], err)
				}
				rank++
				break
			}
		}
	}
	return nil
}
