package censys

import (
	"fmt"
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

func (e *BaseEntry) GetCertificate() (protocols.Certificate, error) {
	return nil, nil // We don't parse certs from censys
}

func (e *BaseEntry) GetError() error {
	if e.Timestamp == "" {
		return fmt.Errorf("empty date")
	} else if e.GetIP() == nil {
		return fmt.Errorf("cannot parse IP")
	}
	return nil
}
