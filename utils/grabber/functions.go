package grabber

import (
	"fmt"
	"github.com/clcert/osr/utils/protocols"
	"strconv"
	"strings"
	"time"
)

func ParseDate(format, path string) (date time.Time, err error) {
	dirSlice := strings.Split(path, "/")
	for i := len(dirSlice) - 1; i >= 0; i-- {
		date, err = time.Parse(format, dirSlice[i])
		if err == nil && !date.IsZero() {
			break
		}
	}
	return
}

func ParsePort(name string) (uint16, error) {
	nameArr := strings.Split(name, "port")
	if len(nameArr) > 1 {
		name = nameArr[len(nameArr) - 1] // Protocol is into the last part of the name, after the "port" string
	}
	portNumberStr := Regex.Port.FindString(name)
	if len(portNumberStr) == 0 {
		return 0, fmt.Errorf("port not found in name")
	}
	port, err := strconv.Atoi(portNumberStr)
	if err != nil {
		return 0, err
	}
	return uint16(port), nil
}

func ParseProtocol(name string) (protocol protocols.PortProtocol) {
	protocol = protocols.TCP
	protocolStr := Regex.Protocol.FindString(name)
	if len(protocolStr) != 0 {
		if protocolStr == "tcp" {
			protocol = protocols.TCP
		} else if protocolStr == "udp" {
			protocol = protocols.UDP
		}
	}
	return
}
