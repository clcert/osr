package protocols

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)

var ftpRegexes = map[string]string{
	"software": "[a-z\\.\\-\\_ ]*", // Added space
	"version":  "([0-9]+)((\\.|p)?[0-9]+)*",
	"brackets": "\\[.*\\]",
	"welcome":  "welcome to",
	"none": "\\(none\\)", // Weird thing from GNU FTP server
	"ftpService": "ftp serv(er|ice)",
}

// FTPParser parses version and software from FTP banner headers.
type FTPParser struct {
	Regexes map[string]*regexp.Regexp
	inited  bool
}

// If message starts with 220, almost certainly is a SSH server
func (p *FTPParser) IsValid(banner string) bool {
	p.init()
	return len(banner) >= 3 && banner[:3] == "220"
}


// init prepares the parser for its use.
// It should be automatically executed as first method of getVersion and getSoftware.
func (p *FTPParser) init() {
	if !p.inited {
		for regexType, regex := range ftpRegexes {
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

// It returns the rest of the banner if the string starts with SSH, a hypen, something (ssh version grly) and another hypen
func (p *FTPParser) CheckFTP(banner string) string {
	banner = strings.Split(banner, "\n")[0]
	// trimming IPs and things in brackets
	banner = p.Regexes["brackets"].ReplaceAllString(banner, "")
	banner = p.Regexes["welcome"].ReplaceAllString(banner, "")
	banner = p.Regexes["none"].ReplaceAllString(banner, "")
	if len(banner) > 3 && banner[:3] == "220" {
		banner = strings.TrimSpace(banner[4:])
	} else {
		return ""
	}
	return banner
}

func (p *FTPParser) GetVersion(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["software"].FindString(p.CheckFTP(banner)))
}

func (p *FTPParser) GetSoftware(banner string) string {
	p.init()
	banner = strings.ToLower(banner)
	return strings.TrimSpace(p.Regexes["version"].FindString(p.CheckFTP(banner)))
}
