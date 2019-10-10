package grabber

import (
	"fmt"
	"github.com/clcert/osr/utils/protocols"
	"net"
	"time"
)

// Base Grabber Entry
// It includes a IP scanned and an Error string.
// The error string could be not empty but still have some data around there
type BaseEntry struct {
	IP    string `json:"ip"`    // Scanned IP
	Error string `json:"error"` // Reported Error
	*CertMeta
}

func (e *BaseEntry) GetIP() net.IP {
	return net.ParseIP(e.IP)
}

func (e *BaseEntry) GetTime(format string, defaultTime time.Time) time.Time {
	return defaultTime // Grabber results don't have date :(
}

func (e *BaseEntry) GetError() error {
	if e.Error == "" {
		return nil
	}
	return fmt.Errorf(e.Error)
}

func (e *BaseEntry) GetCertificate() (protocols.Certificate, error) {
	if e.CertMeta == nil {
		return nil, fmt.Errorf("cert not found")
	}
	return e.CertMeta, nil
}