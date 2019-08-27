package ips

import (
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/clcert/osr/utils/ips"
	"io"
	"net"
)

func GetIPs(csv *utils.HeadedCSV, args *tasks.Args) ips.IPChan {
	ch := make(ips.IPChan)
	go func(){
		for {
			next, err := csv.NextRow()
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			ch <- net.ParseIP(next["ip"])
		}
		close(ch)
	}()
	return ch
}

func CompareIPs(ch1, ch2 ips.IPChan) (chCombined, chMore, chLess ips.IPChan) {
	chCombined = make(ips.IPChan)
	chMore = make(ips.IPChan)
	chLess = make(ips.IPChan)
	go func () {
		var ip1, ip2 net.IP
		var open bool
		for {
			if ip1 == nil {
				ip1, open = <-ch1
				if !open {
					break
				}
			}
			if ip2 == nil {
				ip2, open = <-ch2
				if !open {
					break
				}
			}
			switch ips.CompareBytes(ip1, ip2) {
			case -1:
				chLess <- ip1
				ip1 = nil
			case 0:
				chCombined <- ip1
				ip1 = nil
				ip2 = nil
			case 1:
				chMore <- ip2
				ip2 = nil
			}
		}
		close(chCombined)
		for val := range ch2 {
			chMore <- val
		}
		close(chMore)
		for val := range ch1 {
			chLess <- val
		}
		close(chLess)
	}()
	return
}
