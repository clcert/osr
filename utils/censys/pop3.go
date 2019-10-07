package censys

import (
	"github.com/clcert/osr/utils/protocol_parser"
)

// Represents a dumped Censys POP3 Entry
type POP3Entry struct {
	Parser protocol_parser.POP3Parser
	*BasicEntry
	Banner string `json:"banner"`
}

func (e *POP3Entry) GetBasicEntry() *BasicEntry {
	return e.BasicEntry
}

func (e *POP3Entry) GetSoftware() string {
	return e.Parser.GetSoftware(e.Banner)
}

func (e *POP3Entry) GetVersion() string {
	return e.Parser.GetVersion(e.Banner)
}

func (e *POP3Entry) GetExtra() string {
	return e.Banner
}
