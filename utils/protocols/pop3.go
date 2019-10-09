package protocols

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)

var pop3Regexes = map[string]string{
	"software": "[a-z\\.\\-\\_]*",
	"version":  "[0-9\\.\\_]*",
	"ready":    "( ?is )?ready ?.*", // Some POP3 Servers write "ready at DATE"
	"server":   " ?serv(er|ice) ?",
}

// POP3Parser parses version and software from POP3 banners.
type POP3Parser struct {
	Regexes map[string]*regexp.Regexp
	inited  bool
}

// If message starts with +OK, almost certainly is a SSH server
func (p *POP3Parser) IsValid(banner string) bool {
	return len(banner) >= 3 && banner[:3] == "+OK"
}


// init prepares the parser for its use.
// It should be automatically executed as first method of getVersion and getSoftware.
func (p *POP3Parser) init() {
	if !p.inited {
		for regexType, regex := range pop3Regexes {
			regex, err := regexp.Compile(regex)
			if err != nil {
				panic(panics.Info{
					Text: fmt.Sprintf("%s Regex is not well defined: %s", regexType, regex),
					Err:  err,
				})
			}
			p.Regexes[regexType] = regex
		}
	}
	p.inited = true
}

// It returns the rest of the banner if the string starts with +OK
func (p *POP3Parser) CheckPOP3(banner string) string {
	banner = strings.Split(banner, "\n")[0]
	banner = p.Regexes["ready"].ReplaceAllString(banner, "")
	banner = p.Regexes["server"].ReplaceAllString(banner, "")
	if len(banner) > 3 && banner[:3] == "+OK" {
		banner = strings.TrimSpace(banner[4:])
	} else {
		return ""
	}
	domainSoftware := strings.SplitN(banner, " ", 2)
	banner = domainSoftware[len(domainSoftware)-1]
	return banner
}

func (p *POP3Parser) GetVersion(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckPOP3(banner)))
}

func (p *POP3Parser) GetSoftware(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckPOP3(banner)))
}
