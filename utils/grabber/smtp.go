package grabber

// SMTPEntry represents a scan of a SMTP server
type SMTPEntry struct {
	BaseEntry
	CertMeta // Certificate Metainformation (SMTP port 25, 465)
	Banner string `json:"banner"`// Banner of file (port 25)
	Help   string `json:"help"`// Message when HELP command is sent (port 25)
	Ehlo   string `json:"ehlo"`// Message when EHLO command is sent (port 25)
}


func (e *SMTPEntry) GetBanner() string {
	return e.Banner
}
