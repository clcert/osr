package grabber

import (
	"net"
	"time"
)

// Base Grabber Entry
// It includes a IP scanned and an Error string.
// The error string could be not empty but still have some data around there
type BaseEntry struct {
	IP    string `json:"ip"`    // Scanned IP
	Error string `json:"error"` // Reported Error
}

func (e *BaseEntry) GetIP() net.IP {
	return net.ParseIP(e.IP)
}

func (e *BaseEntry) GetTime(format string, defaultTime time.Time) time.Time {
	return defaultTime // Grabber results don't have date :(
}

func (e *BaseEntry) HasError() bool {
	return e.Error != ""
}

