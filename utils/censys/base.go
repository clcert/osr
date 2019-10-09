package censys

import (
	"net"
	"time"
)

type BasicEntry struct {
	Timestamp string `json:"timestamp"`
	IP        string `json:"ip"`
}

func (e *BasicEntry) GetIP() net.IP {
	return net.ParseIP(e.IP)
}

func (e *BasicEntry) GetTime(formatter string, defaultDate time.Time) time.Time {
	t, err := time.Parse(e.Timestamp, formatter)
	if err != nil {
		t = defaultDate
	}
	return t
}
