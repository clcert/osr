package censys

// Represents a dumped Censys SMTP Entry
type SMTPEntry struct {
	*BaseEntry
	Banner string `json:"banner"`
	Ehlo   string `json:"ehlo"`
}

func (e *SMTPEntry) GetBanner() string {
	return e.Banner
}
