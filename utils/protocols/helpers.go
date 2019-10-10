package protocols

// PortProtocol represents the transport protocol checked in a port scan.
type PortProtocol int

const (
	UnknownProtocol PortProtocol = iota
	TCP
	UDP
)

var PortToProtocol = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	80:   "http",
	8080: "http",
	8000: "http",
	443:  "https",
	8443: "https",
	25:   "smtp",
	465:  "smtp",
	110:  "pop3",
	995:  "pop3",
	143:  "imap",
	993:  "imap",
}

var ProtocolToPorts = map[string][]uint16{
	"ftp":   {21},
	"ssh":   {22},
	"http":  {80, 8000, 8080},
	"https": {443, 8443},
	"smtp":  {25, 465},
	"imap":  {143, 993},
	"pop3":  {110, 995},
}

// returns UDP if the port scanned is related to an UDP protocol.
func GetTransport(port uint16) PortProtocol {
	switch port {
	case 53, 623, 123, 520, 1900, 20000, 47808:
		return UDP
	default:
		return TCP
	}
}
