package censys

import (
	"github.com/clcert/osr/utils/protocol_parser"
)

// Represents a dumped Censys SMTP Entry
type SMTPEntry struct {
	Parser protocol_parser.SMTPParser
	*BasicEntry
	Banner string `json:"banner"`
	Ehlo   string `json:"ehlo"`
}

func (e *SMTPEntry) GetBasicEntry() *BasicEntry {
	return e.BasicEntry
}

func (e *SMTPEntry) GetSoftware() string {
	return e.Parser.GetSoftware(e.Banner)
}

func (e *SMTPEntry) GetVersion() string {
	return e.Parser.GetVersion(e.Banner)
}

func (e *SMTPEntry) GetExtra() string {
	return e.Banner
}

