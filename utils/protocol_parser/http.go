package protocol_parser

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)


var httpRegexes = map[string]string {
	"software": "[a-z\\.\\-\\_]*",
	"version": "[0-9\\.]*",
}

// HTTPParser parses version and software from HTTP "Server" headers.
type HTTPParser struct {
	Regexes map[string]*regexp.Regexp
	inited bool
}

// HTTP is always valid, after all we got to connect correctly
func (p *HTTPParser) IsValid(string) interface{} {
	return true
}

// init prepares the parser for its use.
// It should be automatically executed as first method of getVersion and getSoftware.
func (p *HTTPParser) init() {
	if !p.inited {
		for regexType, regex := range httpRegexes {
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

func (p *HTTPParser) GetVersion(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(banner))
}

func (p *HTTPParser) GetSoftware(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(banner))
}