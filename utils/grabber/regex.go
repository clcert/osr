package grabber

import (
	"regexp"
)

var portStr = "[0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]" // 1-65535
var protocolStr = "(tcp)|(udp)"

type Regexes struct {
	Port       *regexp.Regexp
	Protocol   *regexp.Regexp
}

var Regex Regexes

func init() {
	portRegexp, err := regexp.Compile(portStr)
	if err != nil {
		return
	}
	Regex.Port = portRegexp
	protocolRegexp, err := regexp.Compile(protocolStr)
	if err != nil {
		return
	}
	Regex.Protocol = protocolRegexp
	return
}