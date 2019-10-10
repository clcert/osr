package censys

// Represents a dumped Censys SMTP Entry
type IMAPEntry struct {
	*BaseEntry
	Banner string `json:"banner"`
}

func (e *IMAPEntry) GetBanner() string {
	return e.Banner
}