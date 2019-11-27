package protocols

import (
	"crypto/x509"
	"github.com/clcert/osr/models"
	"net"
	"time"
)

type Entry interface {
	GetIP() net.IP
	GetTime(string, time.Time) time.Time
	GetError() error
	GetBanner() string
	GetCertificate() (Certificate, error)
}

type Certificate interface {
	CheckValid(time.Time) (models.CertStatus, error)
	GetKeySize() int
	GetExpirationDate() time.Time
	GetOrganizationName() string
	GetOrganizationURL() string
	GetAuthority() string
	GetSigAlgorithm() x509.SignatureAlgorithm
	GetTLSProtocol() models.TLSProto
}
