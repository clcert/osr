package grabber

// IMAPFile represents a scan of a IMAP server
type IMAPFile struct {
	BaseFile
	CertMeta // Port 143, 993
	Banner string `json:"banner"` // Protocol banner
}
