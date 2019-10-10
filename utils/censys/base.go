package censys

import (
	"net"
	"time"
)

type BaseEntry struct {
	Timestamp string `json:"timestamp"`
	IP        string `json:"ip"`
}

func (e *BaseEntry) GetIP() net.IP {
	return net.ParseIP(e.IP)
}

func (e *BaseEntry) GetTime(formatter string, defaultDate time.Time) time.Time {
	t, err := time.Parse(e.Timestamp, formatter)
	if err != nil {
		t = defaultDate
	}
	return t
}

func (e *BaseEntry) HasError() bool {
	return false // Censys doesnt report errors
}
