package grabber

import (
	"crypto/x509"
	"github.com/clcert/osr/utils/protocols"
	"strconv"
	"strings"
	"time"
)

// CertProtocols shows the supported HTTPS protocols of a scan
type CertProtocols struct {
	SSL30 bool `json:"SSL_30"`
	TLS10 bool `json:"TLS_10"`
	TLS11 bool `json:"TLS_11"`
	TLS12 bool `json:"TLS_12"`
	TLS13 bool `json:"TLS_13"`
}

// Certificate includes all the fields scanned for a certificate.
type Certificate struct {
	CertificateAuthority string `json:"certificateAuthority"` // Certificate Emmisor
	SignatureAlgorithm   string `json:"signatureAlgorithm"`   // Signature Algorithm
	ExpiredTime          string `json:"expiredTime"`          // Expired Certificate Date
	OrganizationName     string `json:"organizationName"`     // Name related to the cert
	OrganizationURL      string `json:"organizationURL"`      // URL related to the cert
	KeyBits              string `json:"keyBits"`              // Number of key bits
	PemCert              string // The certificate itself in PEM format.
}

type HeartbleedData struct {
	Heartbeat  bool // True if heartbeat detected
	Heartbleed bool // True if Heartbleed vuln detected
}

const certTimeFormat = "Feb 01, 2006 3:04:05 PM"

// CertMeta represents a scan of a Certificate inside other scan.
type CertMeta struct {
	HeartbleedData `json:"heartbleedData"`                  // Status on Heartbleed
	Chain          []Certificate     `json:"chain"`         // Array of certificates representing the Certificate Chain of Trust
	Protocols      CertProtocols     `json:"protocols"`     // Certificate protocols available
	CiphersSuites  map[string]string `json:"ciphersSuites"` // Cipher suites available
}

func (m *CertMeta) IsAutosigned() bool {
	return len(m.Chain) < 2
}

func (m *CertMeta) GetKeySize() int {
	if len(m.Chain) == 0 {
		return 0
	}
	keySizeInt, err := strconv.ParseInt(m.Chain[0].KeyBits, 10, 16)
	if err != nil {
		return 0
	}
	return int(keySizeInt)
}

func (m *CertMeta) GetExpirationDate() time.Time {
	if len(m.Chain) == 0 {
		return time.Time{}
	}
	t, err := time.Parse(certTimeFormat, m.Chain[0].ExpiredTime)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (m *CertMeta) GetOrganizationName() string {
	if len(m.Chain) == 0 {
		return ""
	}
	return m.Chain[0].OrganizationName
}

func (m *CertMeta) GetOrganizationURL() string {
	if len(m.Chain) == 0 {
		return ""
	}
	return m.Chain[0].OrganizationURL
}

func (m *CertMeta) GetAuthority() string {
	if len(m.Chain) == 0 {
		return ""
	}
	return m.Chain[0].CertificateAuthority
}

func (m *CertMeta) GetSigAlgorithm() x509.SignatureAlgorithm {
	if len(m.Chain) == 0 {
		return x509.UnknownSignatureAlgorithm
	}
	sigName := strings.ReplaceAll(m.Chain[0].SignatureAlgorithm, "with", "-")
	sigAlgorithm := x509.UnknownSignatureAlgorithm + 1
	for ; sigAlgorithm <= x509.SHA512WithRSAPSS; sigAlgorithm++ {
		if sigAlgorithm.String() == sigName {
			return sigAlgorithm
		}
	}
	return x509.UnknownSignatureAlgorithm
}

func (m *CertMeta) GetTLSProtocol() protocols.TLSProto {
	if len(m.Chain) == 0 {
		return protocols.UnknownTLSPRoto
	}
	switch {
	case m.Protocols.TLS13:
		return protocols.TLS13
	case m.Protocols.TLS12:
		return protocols.TLS12
	case m.Protocols.TLS11:
		return protocols.TLS11
	case m.Protocols.TLS10:
		return protocols.TLS10
	case m.Protocols.SSL30:
		return protocols.SSL30
	}
	return protocols.UnknownTLSPRoto
}
