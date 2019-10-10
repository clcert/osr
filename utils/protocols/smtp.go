package protocols

import (
	"strings"
)

var smtpRegexes = map[string]string{
	"ready":     " ?ready.*", // Some SMTP Servers write "ready at DATE"
	"randomhex": " ?\\([0-9a-f]*\\)",
	"dow":       " ?(mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun),? .*",
}

// It returns the rest of the banner if the string starts with 220
func prepareSMTP(p *BannerParser, banner string) string {
	banner = strings.Split(banner, "\n")[0]
	if len(banner) > len(p.StartString) && banner[:len(p.StartString)] == p.StartString {
		banner = strings.TrimSpace(banner[len(p.StartString)+1:])
		banner = p.ExtraRegexes["ready"].ReplaceAllString(banner, "")
		banner = p.ExtraRegexes["randomhex"].ReplaceAllString(banner, "")
		banner = p.ExtraRegexes["dow"].ReplaceAllString(banner, "")
	} else {
		return ""
	}
	return banner
}
