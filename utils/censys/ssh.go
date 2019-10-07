package censys

import (
	"github.com/clcert/osr/utils/protocol_parser"
)

// Represents a Dumped Censys SSH Entry
type SSHEntry struct {
	Parser protocol_parser.SSHParser
	*BasicEntry
	Banner string `json:"banner"`
}

func (e *SSHEntry) GetBasicEntry() *BasicEntry {
	return e.BasicEntry
}

func (e *SSHEntry) GetSoftware() string {
	return e.Parser.GetSoftware(e.Banner)
}

func (e *SSHEntry) GetVersion() string {
	return e.Parser.GetVersion(e.Banner)
}

func (e *SSHEntry) GetExtra() string {
	return e.Banner
}