package compare_historic

import (
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/transform/ips"
	"github.com/sirupsen/logrus"
)

const dateFormat = "2006-01-02 15:04:05-07"

func Execute(args *tasks.Args) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	pairsCompared := 0
	pairChan := ips.ToHeadedCSVs(source)
	for pair := range pairChan {
		if err := pair.HasError(); err != nil {
			args.Log.WithFields(logrus.Fields{
				"pairs_compared": pairsCompared,
			}).Infof("Error getting a pair: %s", err)
			continue
		}
		csv1, csv2 := pair.CSV1, pair.CSV2
		args.Log.WithFields(logrus.Fields{
			"pairs_compared": pairsCompared,
			"file1":          csv1.Name,
			"file2":          csv2.Name,
		}).Infof("Comparing files...")
		if err := CompareNextPair(csv1, csv2, saver, args); err != nil {
			return err
		}
		pairsCompared++
	}
	args.Log.WithFields(logrus.Fields{
		"pairs_compared": pairsCompared,
	}).Infof("Done!")
	return nil
}
