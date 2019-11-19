package compare_historic

import (
	"fmt"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/clcert/osr/utils/ips"
	"net"
	"time"
)

func IPCompareUntilDate(date time.Time, dateFormat string) utils.RowChanCompareFunc {
	return func(map1, map2 map[string]string) (cmp int8, err error) {
		date1str, ok := map1["date"]
		if !ok {
			err = fmt.Errorf("date key not found on row")
			return
		}
		date2str, ok := map2["date"]
		if !ok {
			err = fmt.Errorf("date key not found on row")
			return
		}
		date1, err := time.Parse(dateFormat, date1str)
		if err != nil {
			return
		}
		date2, err := time.Parse(dateFormat, date2str)
		if err != nil {
			return
		}
		if date1.After(date) && date2.After(date) {
			err = fmt.Errorf("remaining values are from a future date")
			return
		} else if date1.After(date) {
			cmp = 1
		} else if date2.After(date) {
			cmp = -1
		} else {
			ip1str, ok := map1["ip"]
			if !ok {
				err = fmt.Errorf("ip key not found on row")
				return
			}
			ip2str, ok := map2["ip"]
			if !ok {
				err = fmt.Errorf("ip key not found on row")
				return
			}
			ip1, ip2 := net.ParseIP(ip1str), net.ParseIP(ip2str)
			cmp = ips.CompareBytes(ip1, ip2)
		}
		return
	}
}

func getMinDate(chan1, chan2 *utils.RowChan) (date time.Time, err error) {
	var date1, date2 time.Time
	if !chan1.IsOpen() && !chan2.IsOpen() {
		err = fmt.Errorf("both channels closed")
		return
	}
	if chan1.IsOpen() {
		dateStr1 := chan1.Get()
		date1, err = time.Parse(dateFormat, dateStr1["date"])
		if err != nil {
			return
		}
	}
	if chan2.IsOpen() {
		dateStr2 := chan2.Get()
		date2, err = time.Parse(dateFormat, dateStr2["date"])
		if err != nil {
			return
		}
	}
	if !date1.IsZero() && date1.Before(date2) {
		date = date1
		return
	} else if !date2.IsZero() && date2.Before(date1) {
		date = date2
		return
	} else {
		err = fmt.Errorf("unknown error: this should not happen")
		return
	}
}

func CompareNextPair(csv1, csv2 *utils.HeadedCSV, saver savers.Saver, args *tasks.Args) error {
	if !csv1.HasHeader("ip") || !csv2.HasHeader("date") {
		return fmt.Errorf("file must have ip and date headers")
	}
	chan1 := utils.CSVToRowChan(csv1)
	chan2 := utils.CSVToRowChan(csv2)
	for chan1.IsOpen() || chan2.IsOpen() {
		date, err := getMinDate(chan1, chan2)
		if err != nil {
			args.Log.Infof("Logging date %s...", date)
			break
		}
		both, chan1Uniq, chan2Uniq := chan1.Compare(chan2, IPCompareUntilDate(date, dateFormat))
		line := map[string]interface{}{
			"date":    date,
			"both":    both.Count(),
			csv1.Name: chan1Uniq.Count(),
			csv2.Name: chan2Uniq.Count(),
		}
		err = saver.Save(line)
		if err != nil {
			args.Log.Errorf("cannot save entry: %s", err)
		}
	}
	return nil
}
