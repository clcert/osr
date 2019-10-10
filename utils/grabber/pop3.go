package grabber

// IMAPEntry represents a scan of a POP3 server
type POP3Entry struct {
	BaseEntry
	Banner string `json:"banner"` // Protocol banner (port 110)
}


func (e *POP3Entry) GetBanner() string {
	return e.Banner
}
