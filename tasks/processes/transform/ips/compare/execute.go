package compare

import (
	"fmt"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/transform/ips"
	"github.com/clcert/osr/utils"
	"github.com/sirupsen/logrus"
)

func Execute(args *tasks.Args) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	pairsCompared := 0
	pairChan := ips.ToHeadedCSVs(source)
	for pair := range pairChan {
		var common, less, more *utils.RowChan
		if err := pair.HasError(); err != nil {
			args.Log.WithFields(logrus.Fields{
				"pairs_compared": pairsCompared,
			}).Infof("Error getting a pair: %s", err)
			continue
		}
		csv1, csv2 := pair.CSV1, pair.CSV2
		if csv1.HasHeader("subnet") && csv2.HasHeader("subnet") {
			subnet1, err := ips.GetSubnets(csv1, args)
			if err != nil {
				return fmt.Errorf("cannot parse first file: %s", err)
			}
			subnet2, err := ips.GetSubnets(csv2, args)
			if err != nil {
				return fmt.Errorf("cannot parse second file: %s", err)
			}
			// With subnet intersection
			commonNet, lessNet, moreNet := ips.CompareSubnets(subnet1, subnet2)
			common, less, more = commonNet.ToRowChan(), lessNet.ToRowChan(), moreNet.ToRowChan()
		} else if csv1.HasHeader("ip") && csv2.HasHeader("ip") {
			chan1 := utils.CSVToRowChan(csv1)
			chan2 := utils.CSVToRowChan(csv2)
			common, less, more = chan1.Compare(chan2, ips.IPCompare)
		} else {
			return fmt.Errorf("both files must have the same header between this list of headers: [ip, subnet], but they have this headers: file1: %+v; file2: %+v", csv1.Headers, csv2.Headers)
		}
		ips.SaveIPs(saver, common, "CommonIPs", args)
		ips.SaveIPs(saver, less, "LessIPs", args)
		ips.SaveIPs(saver, more, "MoreIPs", args)
	}
	return nil
}
