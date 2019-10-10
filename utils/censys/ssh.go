package censys

// Represents a Dumped Censys SSH Entry
type SSHEntry struct {
	*BaseEntry
	Raw string `json:"raw"`
}

func (e *SSHEntry) GetBanner() string {
	return e.Raw
}