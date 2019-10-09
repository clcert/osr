package censys

import (
	"github.com/clcert/osr/utils/protocols"
)

// Represents a Dumped Censys SSH Entry
type SSHEntry struct {
	Parser protocols.SSHParser
	*BasicEntry
	Banner string `json:"banner"`
}

func (e *SSHEntry) GetBanner() string {
	return e.Banner
}