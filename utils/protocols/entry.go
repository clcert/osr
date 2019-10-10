package protocols

import (
	"crypto/x509"
	"net"
	"time"
)

type Entry interface {
	GetIP() net.IP
	GetTime(string, time.Time) time.Time
	GetError() error
	GetBanner() string
	GetCertificate() Certificate
}

type Certificate interface {
	IsAutosigned() bool
	GetKeySize() int
	GetExpirationDate() time.Time
	GetOrganizationName() string
	GetOrganizationURL() string
	GetAuthority() string
	GetSigAlgorithm() x509.SignatureAlgorithm
	GetTLSProtocol() TLSProto
}

type TLSProto int

const (
	UnknownTLSPRoto TLSProto = iota
	SSL30
	TLS10
	TLS11
	TLS12
	TLS13
)

