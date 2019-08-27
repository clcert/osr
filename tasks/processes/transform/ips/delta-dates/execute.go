package compare

import (
	"fmt"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/transform/ips"
	"github.com/clcert/osr/utils"
	ips2 "github.com/clcert/osr/utils/ips"
)

func Execute(args *tasks.Args) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	file1 := source.Next()
	if file1 == nil {
		return fmt.Errorf("no files in source. Cannot compare")
	}
	file2 := source.Next()
	if file2 == nil {
		return fmt.Errorf("only one file in source. Cannot compare")
	}
	reader1, err := file1.Open()
	if err != nil {
		return fmt.Errorf("couldn't open first file: %s", err)
	}
	reader2, err := file2.Open()
	if err != nil {
		return fmt.Errorf("couldn't open second file: %s", err)
	}
	csv1, err := utils.NewHeadedCSV(reader1, nil)
	if err != nil {
		return fmt.Errorf("couldn't use first file as CSV: %s", err)
	}
	csv2, err := utils.NewHeadedCSV(reader2, nil)
	if err != nil {
		return fmt.Errorf("couldn't use second file as CSV: %s", err)
	}
	var common, less, more ips2.IPChan
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
		common, less, more = commonNet.ToIPChan(), lessNet.ToIPChan(), moreNet.ToIPChan()
	} else if csv1.HasHeader("ip") && csv2.HasHeader("ip") {
		ips1 := ips.GetIPs(csv1, args)
		ips2 := ips.GetIPs(csv2, args)
		common, less, more = ips.CompareIPs(ips1, ips2)
	} else {
		return fmt.Errorf("both files must have the same header between this list of headers: [ip, subnet], but they have this headers: file1: %+v; file2: %+v", csv1.Headers, csv2.Headers)
	}
	ips.SaveIPs(saver, common, "CommonIPs", args)
	ips.SaveIPs(saver, less, "LessIPs", args)
	ips.SaveIPs(saver, more, "MoreIPs", args)
	// TODO: log done!
	return nil
}
