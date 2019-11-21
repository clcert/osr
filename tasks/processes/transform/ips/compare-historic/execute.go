package compare_historic

import (
	"fmt"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/transform/ips"
	"github.com/clcert/osr/utils"
	"github.com/sirupsen/logrus"
	"strings"
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


func CompareNextPair(csv1, csv2 *utils.HeadedCSV, saver savers.Saver, args *tasks.Args) error {
	if !csv1.HasHeader("ip") || !csv2.HasHeader("date") || !csv2.HasHeader("port_number") {
		return fmt.Errorf("file must have ip, port number and date headers")
	}
	csv1Name := strings.TrimSuffix(csv1.Name, ".csv")
	csv2Name := strings.TrimSuffix(csv2.Name, ".csv")
	chan1 := utils.CSVToRowChan(csv1)
	chan2 := utils.CSVToRowChan(csv2)
	for chan1.IsOpen() || chan2.IsOpen() {
		port, date, err := GetMinPortAndDate(chan1, chan2)
		if err != nil {
			args.Log.Errorf("cannot get min date: %s", err)
			break
		}
		args.Log.Infof("Logging date %s and port %d", date, port)

		if date.Year() == 2019 && date.Month() == 9 && date.Day() == 23 {
			fmt.Printf("nooo");
		}

		joinChan := chan1.Join(chan2, IPCompareUntilPortAndDate(port, date, dateFormat))
		countIPs := make(map[string]int)
		for joinChan.IsOpen() {
			ip := joinChan.Get()
			if _, ok := countIPs[ip["tag"]]; !ok {
				countIPs[ip["tag"]] = 0
			}
			countIPs[ip["tag"]]++
		}
		line := map[string]string{
			"date":   date.String(),
			"port":   fmt.Sprintf("%d", port),
			csv1Name: fmt.Sprintf("%d", countIPs[csv1.Name]),
			csv2Name: fmt.Sprintf("%d", countIPs[csv2.Name]),
			"both":   fmt.Sprintf("%d", countIPs["both"]),
		}
		err = saver.Save(savers.Savable{
			Object: line,
			Meta: map[string]string{
				"outID": fmt.Sprintf(
					"compare_%s_%s",
					csv1Name,
					csv2Name),
			},
		})
		if err != nil {
			args.Log.Errorf("cannot save entry: %s", err)
		}
	}
	return nil
}
