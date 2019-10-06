package grabber

// SMTPFile represents a scan of a SMTP server
type SMTPFile struct {
	BaseFile
	CertMeta // Certificate Metainformation (SMTP port 25, 465)
	Banner string `json:"banner"`// Banner of file (port 25)
	Help   string `json:"help"`// Message when HELP command is sent (port 25)
	Ehlo   string `json:"ehlo"`// Message when EHLO command is sent (port 25)
}
