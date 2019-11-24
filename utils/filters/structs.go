package filters

import (
	"net"
	"time"
)

type DateConfig struct {
	Since time.Time
	Until time.Time
}

type ScanConfig struct {
	*DateConfig
	Blacklist map[uint16]struct{}
	SourceIP  net.IP
}

func (conf *ScanConfig) IsNotInBlacklist(port uint16) bool {
	_, ok := conf.Blacklist[port]
	return !ok
}

func (conf *DateConfig) IsDateInRange(date time.Time) bool {
	if (!conf.Since.IsZero() && date.Before(conf.Since)) ||
		(!conf.Until.IsZero() && date.After(conf.Until)) {
		return false
	}
	return true
}
