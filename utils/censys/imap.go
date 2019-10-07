package censys

import (
	"github.com/clcert/osr/utils/protocol_parser"
)

// Represents a dumped Censys SMTP Entry
type IMAPEntry struct {
	Parser protocol_parser.IMAPParser
	*BasicEntry
	Banner string `json:"banner"`
}

func (e *IMAPEntry) GetBasicEntry() *BasicEntry {
	return e.BasicEntry
}

func (e *IMAPEntry) GetSoftware() string {
	return e.Parser.GetSoftware(e.Banner)
}

func (e *IMAPEntry) GetVersion() string {
	return e.Parser.GetVersion(e.Banner)
}

func (e *IMAPEntry) GetExtra() string {
	return e.Banner
}
