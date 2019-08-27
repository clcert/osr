package port_scan

import (
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"regexp"
	"strconv"
	"strings"
)

var portStr = "[0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]" // 1-65535
var protocolStr = "(tcp)|(udp)"

type Regexes struct {
	Port       *regexp.Regexp
	Protocol   *regexp.Regexp
	UDPSpecial []*regexp.Regexp
}

func NewRegex(args *tasks.Args) (regexes *Regexes, err error) {
	regexes = &Regexes{
		UDPSpecial: make([]*regexp.Regexp, 0),
	}
	portRegexp, err := regexp.Compile(portStr)
	if err != nil {
		return
	}
	regexes.Port = portRegexp
	protocolRegexp, err := regexp.Compile(protocolStr)
	if err != nil {
		return
	}
	regexes.Protocol = protocolRegexp
	udpParams, ok := args.Params["udp"]
	if ok {
		regexpStrings := strings.Split(udpParams, ",")
		for _, regexpStr := range regexpStrings {
			regex, err := regexp.Compile(regexpStr)
			if err != nil {
				continue
			}
			regexes.UDPSpecial = append(regexes.UDPSpecial, regex)
		}
	}
	return
}

func (regex *Regexes) GetPort(name string) (uint16, error) {
	portNumberStr := regex.Port.FindString(name)
	if len(portNumberStr) == 0 {
		return 0, fmt.Errorf("port not found in name")
	}
	port, err := strconv.Atoi(portNumberStr)
	if err != nil {
		return 0, err
	}
	return uint16(port), nil
}

func (regex *Regexes) GetProtocol(name string) (protocol models.PortProtocol) {
	protocol = models.TCP
	protocolStr := regex.Protocol.FindString(name)
	if len(protocolStr) == 0 {
		for _, regex := range regex.UDPSpecial {
			if regex.MatchString(name) {
				protocol = models.UDP
				break
			}
		}
	} else {
		if protocolStr == "tcp" {
			protocol = models.TCP
		} else if protocolStr == "udp" {
			protocol = models.UDP
		}
	}
	return
}
