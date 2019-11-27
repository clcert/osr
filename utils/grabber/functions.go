package grabber

import (
	"fmt"
	"github.com/clcert/osr/models"
	"strconv"
	"strings"
	"time"
)

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("America/Santiago")
}

func ParseDate(format, path string) (date time.Time, err error) {
	dirSlice := strings.Split(path, "/")
	for i := len(dirSlice) - 1; i >= 0; i-- {
		date, err = time.ParseInLocation(format, dirSlice[i], loc)
		if err == nil && !date.IsZero() {
			break
		}
	}
	return
}

func ParsePort(name string) (uint16, error) {
	nameArr := strings.Split(name, "port")
	if len(nameArr) > 1 {
		name = nameArr[len(nameArr)-1] // Protocol is into the last part of the name, after the "port" string
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

func ParseProtocol(name string) (protocol models.PortProtocol) {
	protocol = models.TCP
	protocolStr := Regex.Protocol.FindString(name)
	if len(protocolStr) != 0 {
		if protocolStr == "tcp" {
			protocol = models.TCP
		} else if protocolStr == "udp" {
			protocol = models.UDP
		}
	}
	return
}
