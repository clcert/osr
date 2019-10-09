package grabber

// HTTPEntry represents an HTTP Header Scan
type HTTPEntry struct {
	BaseEntry
	Header      map[string][]string `json:"header"`      // HTTP Headers
	Index       string              `json:"index"`       // Text representation of index webpage
	TlsProtocol string              `json:"tlsProtocol"` // TLS protocol used on connection
	CipherSuite string              `json:"cipherSuite"` // Cipher Suite used
	CertMeta
}

func (e *HTTPEntry) GetBanner() string {
	if s, ok := e.Header["Server"]; ok  && len(s) >=1 {
		return s[0]
	}
	return ""
}