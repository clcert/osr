package grabber

// FTPEntry represents a scan of a FTP server
type FTPEntry struct {
	BaseEntry
	Banner string `json:"banner"` // Protocol banner
}

func (e *FTPEntry) GetBanner() string {
	return e.Banner
}
