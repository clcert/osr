package protocols

import (
	"strings"
)

var imapRegexes = map[string]string{
	"ready":    "( ?is )?ready.*", // Some IMAP Servers write "ready at DATE"
	"server":   " ?serv(er|ice)",
	"brackets": " ?\\[.*\\]",
	"lessmore": " ?<.*>",
	"dow":      " ?(mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun),? .*",
}

// It returns the rest of the banner if the string starts with * OK
func prepareIMAP(p *BannerParser, banner string) string {
	banner = strings.Split(banner, "\n")[0]
	banner = p.ExtraRegexes["ready"].ReplaceAllString(banner, "")
	banner = p.ExtraRegexes["server"].ReplaceAllString(banner, "")
	banner = p.ExtraRegexes["brackets"].ReplaceAllString(banner, "")
	banner = p.ExtraRegexes["lessmore"].ReplaceAllString(banner, "")
	banner = p.ExtraRegexes["dow"].ReplaceAllString(banner, "")
	if len(banner) > len(p.StartString) && banner[:len(p.StartString)] == p.StartString {
		banner = strings.TrimSpace(banner[len(p.StartString)+1:])
	} else {
		return ""
	}
	return banner
}
