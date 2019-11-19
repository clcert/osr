package ips

import (
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/clcert/osr/utils/ips"
	"net"
)

func SaveIPs(saver savers.Saver, rows *utils.RowChan, id string, args *tasks.Args) {
	for rows.IsOpen() {
		row := rows.Get()
		ipStr, ok := row["ip"];
		if !ok {
			args.Log.Errorf("cannot save ip: ip field not found")
			continue
		}
		ip := net.ParseIP(ipStr)
		err := saver.Save(savers.Savable{
			Object: struct{IP net.IP `structs:",string"`}{ip},
			Meta: map[string]string{
				"outID": id,
			},
		})
		if err != nil {
			args.Log.Errorf("cannot save ip: %s", err)
		}
	}
}

func SaveSubnets(saver savers.Saver, subnets ips.SubnetList, id string, args *tasks.Args) {
	for _, subnet := range subnets {
		err := saver.Save(savers.Savable{
			Object: struct{Subnet *net.IPNet `structs:",string"`}{subnet},
			Meta: map[string]string{
				"outID": id,
			},
		})
		if err != nil {
			args.Log.Errorf("cannot save subnet: %s", err)
		}
	}
}
