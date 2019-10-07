package protocol_parser

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)

var sshRegexes = map[string]string {
	"software": "[a-z\\.\\-\\_ ]*", // Added space
	"version": "([0-9]+)((\\.|p)?[0-9]+)*",
}

// SSHParser parses version and software from SSH banner headers.
type SSHParser struct {
	Regexes map[string]*regexp.Regexp
	inited        bool
}

// init prepares the parser for its use.
// It should be automatically executed as first method of getVersion and getSoftware.
func (p *SSHParser) init() {
	if !p.inited {
		for regexType, regex := range sshRegexes {
			regex, err := regexp.Compile(regex)
			if err != nil {
				panic(panics.Info{
					Text:        fmt.Sprintf("%s Regex is not well defined: %s", regexType, regex),
					Err:         err,
				})
			}
			p.Regexes[regexType] = regex
		}
	}
	p.inited = true
}


// It returns the rest of the banner if the string starts with SSH, a hypen, something (ssh version grly) and another hypen
func (p *SSHParser) CheckSSH (banner string) string {
	firstLine := strings.Split(banner,"\n")[0]
	bannerParts := strings.SplitN(firstLine, "-", 3)
	if bannerParts[0] == "SSH" && len(bannerParts) > 2 {
		return bannerParts[2]
	}
	return ""
}

func (p *SSHParser) GetVersion(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckSSH(banner)))
}

func (p *SSHParser) GetSoftware(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckSSH(banner)))
}
