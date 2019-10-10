package censys

// Represents a dumped Censys HTTP Entry
type HTTPEntry struct {
	*BaseEntry
	ContentLength string `json:"content_length"`
	ContentType   string `json:"content_Type"`
	Server        string `json:"server"`
	StatusLine    string `json:"status_line"`
}

func (e *HTTPEntry) GetBanner() string {
	return e.Server
}