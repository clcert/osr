package grabber

// HTTPFile represents an HTTP Header Scan
type HTTPFile struct {
	BaseFile
	Header      map[string][]string `json:"header"`      // HTTP Headers
	Index       string              `json:"index"`       // Text representation of index webpage
	TlsProtocol string              `json:"tlsProtocol"` // TLS protocol used on connection
	CipherSuite string              `json:"cipherSuite"` // Cipher Suite used
	CertMeta
}
