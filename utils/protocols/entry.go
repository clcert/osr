package protocols

import (
	"net"
	"time"
)

type Entry interface {
	GetIP() net.IP
	GetTime(string, time.Time) time.Time
	GetBanner() string
}
