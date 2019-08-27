package ips

import (
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/ips"
	"net"
)

func SaveIPs(saver savers.Saver, ips ips.IPChan, id string, args *tasks.Args) {
	for ip := range ips {
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
