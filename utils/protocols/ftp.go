package protocols

import (
	"strings"
)

var ftpRegexes = map[string]string{
	"brackets":   " ?\\[.*\\]",
	"welcome":    " ?welcome to",
	"none":       " ?\\(none\\)", // Weird thing from GNU FTP server
	"ftpService": " ?ftp serv(er)|(ice).*",
	"dow":        " ?(mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun),? .*",
}

// It returns the rest of the banner if the string starts with SSH, a hypen, something (ssh version grly) and another hypen
func prepareFTP(p *BannerParser, banner string) string {
	banner = strings.Split(banner, "\n")[0]
	// trimming IPs and things in brackets
	banner = p.ExtraRegexes["brackets"].ReplaceAllString(banner, "")
	banner = p.ExtraRegexes["welcome"].ReplaceAllString(banner, "")
	banner = p.ExtraRegexes["none"].ReplaceAllString(banner, "")
	banner = p.ExtraRegexes["ftpService"].ReplaceAllString(banner, "")
	banner = p.ExtraRegexes["dow"].ReplaceAllString(banner, "")
	if len(banner) > len(p.StartString) && banner[:len(p.StartString)] == p.StartString {
		banner = strings.TrimSpace(banner[len(p.StartString)+1:])
	} else {
		return ""
	}
	return banner
}
