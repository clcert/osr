package censys

import "github.com/clcert/osr/utils/protocol_parser"

// Represents a dumped Censys HTTP Entry
type HTTPEntry struct {
	Parser protocol_parser.HTTPParser
	*BasicEntry
	ContentLength string `json:"content_length"`
	ContentType   string `json:"content_Type"`
	Server        string `json:"server"`
	StatusLine    string `json:"status_line"`
}

func (e *HTTPEntry) GetBasicEntry() *BasicEntry {
	return e.BasicEntry
}

func (e *HTTPEntry) GetSoftware() string {
	return e.Parser.GetSoftware(e.Server)
}

func (e *HTTPEntry) GetVersion() string {
	return e.Parser.GetVersion(e.Server)
}

func (e *HTTPEntry) GetExtra() string {
	// TODO: insert more data as extra?
	return e.Server
}
