package grabber

// SSHFile represents a scan of a SSH server
type SSHEntry struct {
	BaseEntry
	Banner string `json:"banner"` // Protocol banner
}


func (e *SSHEntry) GetBanner() string {
	return e.Banner
}
