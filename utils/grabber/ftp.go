package grabber

// FTPFile represents a scan of a FTP server
type FTPFile struct {
	BaseFile
	Banner string `json:"banner"` // Protocol banner
}
