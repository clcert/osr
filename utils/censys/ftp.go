package censys

import (
	"github.com/clcert/osr/utils/protocol_parser"
)

// Represents a Dumped Censys FTP Entry
type FTPEntry struct {
	Parser protocol_parser.FTPParser
	*BasicEntry
	Banner string `json:"banner"`
}

func (e *FTPEntry) GetBasicEntry() *BasicEntry {
	return e.BasicEntry
}

func (e *FTPEntry) GetSoftware() string {
	return e.Parser.GetSoftware(e.Banner)
}

func (e *FTPEntry) GetVersion() string {
	return e.Parser.GetVersion(e.Banner)
}

func (e *FTPEntry) GetExtra() string {
	return e.Banner
}
