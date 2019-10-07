package protocol_parser

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)


var smtpRegexes = map[string]string {
	"software": "[a-z\\.\\-\\_]*",
	"version": "[0-9\\.\\_]*",
	"ready": "ready.*", // Some SMTP Servers write "ready at DATE"
}

// SMTPParser parses version and software from SMTP banner.
type SMTPParser struct {
	Regexes map[string]*regexp.Regexp
	inited bool
}

// init prepares the parser for its use.
// It should be automatically executed as first method of getVersion and getSoftware.
func (p *SMTPParser) init() {
	if !p.inited {
		for regexType, regex := range smtpRegexes {
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

// It returns the rest of the banner if the string starts with 220
func (p *SMTPParser) CheckSMTP(banner string) string {
	banner = strings.Split(banner, "\n")[0]
	banner = p.Regexes["ready"].ReplaceAllString(banner, "")
	if len(banner) > 3 && banner[:3] == "220" {
		banner = strings.TrimSpace(banner[4:])
	} else {
		return ""
	}
	domainSoftware := strings.SplitN(banner, " ", 2)
	banner = domainSoftware[len(domainSoftware) - 1]
	return banner
}


func (p *SMTPParser) GetVersion(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckSMTP(banner)))
}

func (p *SMTPParser) GetSoftware(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckSMTP(banner)))
}