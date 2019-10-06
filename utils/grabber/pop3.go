package grabber

// IMAPFile represents a scan of a POP3 server
type POP3File struct {
	BaseFile
	CertMeta // port 110, 995
	Banner string `json:"banner"` // Protocol banner (port 110)
}
