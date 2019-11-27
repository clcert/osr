package grabber

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/clcert/osr/models"
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

const certTimeFormat = "Jan 2, 2006 3:04:05 PM"

// CertMeta represents a scan of a Certificate inside other scan.
type CertMeta struct {
	HeartbleedData `json:"heartbleedData"`                  // Status on Heartbleed
	Chain          []Certificate     `json:"chain"`         // Array of certificates representing the Certificate Chain of Trust
	Protocols      CertProtocols     `json:"protocols"`     // Certificate protocols available
	CiphersSuites  map[string]string `json:"ciphersSuites"` // Cipher suites available
}

func (m *CertMeta) CheckValid(now time.Time) (models.CertStatus, error) {
	if len(m.Chain) > 0 {
		// get first of chain
		block, _ := pem.Decode([]byte(m.Chain[0].PemCert))
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return models.CertUnparseable, err
		}
		if len(m.Chain) == 1 && cert.Issuer.String() == cert.Subject.String() {
			return models.CertSelfSigned, nil
		}
		intermediatePool := x509.NewCertPool()
		for i := 1; i < len(m.Chain); i++ {
			block, _ := pem.Decode([]byte(m.Chain[i].PemCert))
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return models.CertUnparseable, err
			}
			intermediatePool.AddCert(cert)
		}
		_, err = cert.Verify(x509.VerifyOptions{
			Intermediates: intermediatePool,
			CurrentTime:   now,
		})
		if err == nil {
			return models.CertValid, nil
		} else {
			switch e := err.(type) {
			case x509.CertificateInvalidError:
				switch e.Reason {
				case x509.Expired:
					return models.CertExpired, err
				case x509.NotAuthorizedToSign:
					return models.CertNotAuthorizedToSign, err
				default:
					return models.CertUnknownError, err
				}
			case x509.UnknownAuthorityError:
				return models.CertUnknownAuthority, err
			default:
				return models.CertUnknownError, err
			}
		}
	} else {
		return models.CertEmptyChain, nil
	}
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

func (m *CertMeta) GetTLSProtocol() models.TLSProto {
	if len(m.Chain) == 0 {
		return models.UnknownTLSPRoto
	}
	switch {
	case m.Protocols.SSL30:
		return models.SSL30
	case m.Protocols.TLS10:
		return models.TLS10
	case m.Protocols.TLS11:
		return models.TLS11
	case m.Protocols.TLS12:
		return models.TLS12
	case m.Protocols.TLS13:
		return models.TLS13
	}
	return models.UnknownTLSPRoto
}
