package scans

import (
	"net"
	"time"
)

type ScanConfig struct {
	Since time.Time
	Until time.Time
	Blacklist map[uint16]struct{}
	SourceIP net.IP
}

func (conf *ScanConfig) IsPortAllowed(port uint16) bool {
	_, ok := conf.Blacklist[port]
	return !ok
}

func (conf *ScanConfig) IsDateInRange(date time.Time) bool {
	if (!conf.Since.IsZero() && date.Before(conf.Since)) ||
		(!conf.Until.IsZero() && date.After(conf.Until)) {
		return false
	}
	return true
}
