package grabber

// IMAPEntry represents a scan of a IMAP server
type IMAPEntry struct {
	BaseEntry
	CertMeta // Port 143, 993
	Banner string `json:"banner"` // Protocol banner
}

func (e *IMAPEntry) GetBanner() string {
	return e.Banner
}
