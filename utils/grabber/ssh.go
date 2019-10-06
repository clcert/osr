package grabber

// SSHFile represents a scan of a SSH server
type SSHFile struct {
	BaseFile
	Banner string `json:"banner"` // Protocol banner
}
