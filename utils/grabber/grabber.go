package grabber

// Base Grabber File
// It includes a IP scanned and an Error string.
// The error string could be not empty but still have some data around there
type BaseFile struct {
	IP    string `json:"ip"`    // Scanned IP
	Error string `json:"error"` // Reported Error
}


