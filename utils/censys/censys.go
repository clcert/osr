package censys

type BasicEntry struct {
	Timestamp string `json:"timestamp"`
	IP string `json:"ip"`
}

// Represents a Dumped Censys FTP Entry
type FTPEntry struct {
	BasicEntry
	Banner string `json:"banner"`
}

// Represents a Dumped Censys FTP Entry
type SSHEntry struct {
	BasicEntry
	Banner string `json:"banner"`
}

// Represents a dumped Censys SMTP Entry
type SMTPEntry struct {
	BasicEntry
	Banner string `json:"banner"`
	Ehlo string `json:"ehlo"`
}

// Represents a dumped Censys POP3 Entry
type POP3Entry struct {
	BasicEntry
	Banner string `json:"banner"`
}

// Represents a dumped Censys IMAP Entry
type IMAPEntry struct {
	BasicEntry
	Banner string `json:"banner"`
}

// Represents a dumped Censys HTTP Entry
type HTTPEntry struct {
	BasicEntry
	ContentLength string `json:"content_length"`
	ContentType string `json:"content_Type"`
	Server string `json:"server"`
	StatusLine string `json:"status_line"`
}

