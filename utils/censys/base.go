package censys

import (
	"github.com/clcert/osr/utils/protocols"
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

func (e *BaseEntry) GetCertificate() protocols.Certificate {
	return nil // We don't parse certs from censys
}

func (e *BaseEntry) GetError() error {
	return nil // Censys doesnt report errors
}
