package grabber

// IMAPEntry represents a scan of a IMAP server
type IMAPEntry struct {
	BaseEntry
	Banner string `json:"banner"` // Protocol banner
}

func (e *IMAPEntry) GetBanner() string {
	return e.Banner
}
