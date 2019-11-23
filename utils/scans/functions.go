package scans

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func ParseConf(params map[string]string, ip net.IP) (conf *ScanConfig, errors []error) {
	errors = make([]error, 0)
	blacklist, err := parseBlacklist(params)
	if err != nil {
		errors = append(errors, err)
	}
	since, err := parseDate(params, "since")
	if err != nil {
		errors = append(errors, err)
	}
	until, err := parseDate(params, "until")
	if err != nil {
		errors = append(errors, err)
	}

	conf = &ScanConfig{
		SourceIP:  ip, // Censys IP
		Blacklist: blacklist,
		Since:     since,
		Until:     until,
	}
	return
}

func parseBlacklist(params map[string]string) (blacklist map[uint16]struct{}, err error) {
	blacklist = make(map[uint16]struct{})
	unparsedPorts := make([]uint16, 0)
	if bl, ok := params["blacklist"]; ok {
		blSplit := strings.Split(bl, ",")
		for _, portStr := range blSplit {
			port, err := strconv.ParseInt(portStr, 10, 16)
			if err != nil {
				unparsedPorts = append(unparsedPorts, uint16(port))
				continue
			}
			blacklist[uint16(port)] = struct{}{}
		}
	}
	if len(unparsedPorts) > 0 {
		err = fmt.Errorf("couldn't parse the following ports: %v", unparsedPorts)
	}
	return
}

func parseDate(params map[string]string, key string) (date time.Time, err error) {
	sinceParam, ok := params[key]
	if ok {
		date, err = time.Parse("20060102", sinceParam)
	}
	return
}
