package grabber

import "strings"

// HTTPEntry represents an HTTP Header Scan
type HTTPEntry struct {
	BaseEntry
	Header      map[string][]string `json:"header"`      // HTTP Headers
	Index       string              `json:"index"`       // Text representation of index webpage
	TlsProtocol string              `json:"tlsProtocol"` // TLS protocol used on connection
	CipherSuite string              `json:"cipherSuite"` // Cipher Suite used
}

func (e *HTTPEntry) GetBanner() string {
	// Need to do this: server header is not case sensitive
	for k, v := range e.Header {
		if strings.ToLower(k) == "server" {
			return v[0]
		}
	}

	return ""
}