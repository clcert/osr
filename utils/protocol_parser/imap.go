package protocol_parser

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)


var imapRegexes = map[string]string {
	"software": "[a-z\\.\\-\\_]*",
	"version": "[0-9\\.\\_]*",
	"ready": "( ?is )?ready ?.*", // Some IMAP Servers write "ready at DATE"
	"server": " ?serv(er|ice) ?",
	"brackets": "\\[.*\\]",
}

// IMAPParser parses version and software from IMAP banners.
type IMAPParser struct {
	Regexes map[string]*regexp.Regexp
	inited bool
}

// init prepares the parser for its use.
// It should be automatically executed as first method of getVersion and getSoftware.
func (p *IMAPParser) init() {
	if !p.inited {
		for regexType, regex := range imapRegexes {
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

// It returns the rest of the banner if the string starts with * OK
func (p *IMAPParser) CheckIMAP(banner string) string {
	banner = strings.Split(banner, "\n")[0]
	banner = p.Regexes["ready"].ReplaceAllString(banner, "")
	banner = p.Regexes["server"].ReplaceAllString(banner, "")
	banner = p.Regexes["brackets"].ReplaceAllString(banner, "")
	if len(banner) > 4 && banner[:4] == "* OK" {
		banner = strings.TrimSpace(banner[5:])
	}else {
		return ""
	}
	domainSoftware := strings.SplitN(banner, " ", 2)
	banner = domainSoftware[len(domainSoftware) - 1]
	return banner
}


func (p *IMAPParser) GetVersion(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckIMAP(banner)))
}

func (p *IMAPParser) GetSoftware(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckIMAP(banner)))
}