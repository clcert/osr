package censys

import (
	"github.com/clcert/osr/utils/protocols"
)

// Represents a dumped Censys SMTP Entry
type SMTPEntry struct {
	Parser protocols.SMTPParser
	*BasicEntry
	Banner string `json:"banner"`
	Ehlo   string `json:"ehlo"`
}

func (e *SMTPEntry) GetBanner() string {
	return e.Banner
}
