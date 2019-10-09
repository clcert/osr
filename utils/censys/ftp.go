package censys

// Represents a Dumped Censys FTP Entry
type FTPEntry struct {
	*BasicEntry
	Banner string `json:"banner"`
}

func (e *FTPEntry) GetBanner() string {
	return e.Banner
}
