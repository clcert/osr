package ips

import (
	"fmt"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/clcert/osr/utils/ips"
	"io"
	"net"
	"sync"
)

// HeadedCSVPair represents two csv-like files already been opened.
type HeadedCSVPair struct {
	CSV1, CSV2     *utils.HeadedCSV
	Error1, Error2 error
	sync.WaitGroup
}

// Returns an error in any of both files if exists, nil otherwise.
func (pair *HeadedCSVPair) HasError() error {
	if pair.Error1 != nil {
		return pair.Error1
	}
	return pair.Error2
}

// Transforms a source into a csv channel. It processes a new pair immediately when a pair is requested, maintaining
// the order of the entries.
func ToHeadedCSVs(source sources.Source) (csvChan chan *HeadedCSVPair) {
	csvChan = make(chan *HeadedCSVPair)
	go func() {
		defer close(csvChan)
		for {
			file1 := source.Next()
			if file1 == nil {
				return
			}
			file2 := source.Next()
			if file2 == nil {
				return
			}
			pair := &HeadedCSVPair{}
			pair.Add(2)
			go func() {
				defer pair.Done()
				csv, err := EntryToHeadedCSV(file1)
				if err != nil {
					pair.Error1 = err
					return
				}
				pair.CSV1 = csv
			}()
			go func() {
				defer pair.Done()
				csv, err := EntryToHeadedCSV(file2)
				if err != nil {
					pair.Error2 = err
					return
				}
				pair.CSV2 = csv
			}()
			pair.Wait()
			csvChan <- pair
		}
	}()
	return
}

// Transforms an individual entry to a headed CSV
func EntryToHeadedCSV(file sources.Entry) (*utils.HeadedCSV, error) {
	reader, err := file.Open()
	if err != nil {
		err = fmt.Errorf("couldn't open file: %s", err)
		return nil, err
	}
	csv, err := utils.NewHeadedCSV(reader, &utils.HeadedCSVOptions{
		Name: file.Name(),
	})
	if err != nil {
		err = fmt.Errorf("couldn't open file as CSV: %s", err)
		return nil, err
	}
	return csv, nil
}

func Compare(map1, map2 map[string]string) (cmp int8, err error) {
	if map1 == nil {
		cmp = 1
		return
	} else if map2 == nil {
		cmp = -1
		return
	}
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
	return
}

func GetSubnets(csv *utils.HeadedCSV, args *tasks.Args) (ips.SubnetList, error) {
	list := make(ips.SubnetList, 0)
	for {
		line, err := csv.NextRow()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		_, subnet, err := net.ParseCIDR(line["subnet"]) // header should be subnet
		if err != nil {
			args.Log.Error("cannot save ip: %s", err)
			continue
		}
		list = append(list, subnet)
	}
	return list, nil
}
