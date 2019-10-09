package protocols

import "github.com/clcert/osr/models"

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

}

var ProtocolToPorts = map[string][]uint16{
	"ftp": {21},
	"ssh": {22},
	"http": {80, 8000, 8080},
	"https": {443, 8443},
	"smtp": {25, 465},
	"imap": {143, 993},
	"pop3": {110, 995},
}

// returns UDP if the port scanned is related to an UDP protocol.
func GetTransport(port uint16) models.PortProtocol {
	switch port {
	case 53, 623, 123, 520, 1900, 20000, 47808:
		return models.UDP
	default:
		return models.TCP
	}
}