package grabber

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
	SignatureAlgorithm    string `json:"signatureAlgorithm"`   // Signature Algorithm
	ExpiredTime          string `json:"expiredTime"`          // Expired Certificate Date
	OrganizationURL      string `json:"organizationURL"`      // URL related to the cert
	KeyBits              string `json:"keyBits"`              // Number of key bits
	PemCert              string // The certificate itself in PEM format.
}

type HeartbleedData struct {
	Heartbeat  bool // True if heartbeat detected
	Heartbleed bool // True if Heartbleed vuln detected
}

// CertMeta represents a scan of a Certificate inside other scan.
type CertMeta struct {
	HeartbleedData                   `json:"heartbleedData"`// Status on Heartbleed
	Chain          []Certificate     `json:"chain"`// Array of certificates representing the Certificate Chain of Trust
	Protocols      CertProtocols     `json:"protocols"`// Certificate protocols available
	CiphersSuites  map[string]string `json:"ciphersSuites"`// Cipher suites available
}
