package protocols

import (
	"strings"
)

// It returns the rest of the banner if the string starts with SSH, a hypen, something (ssh version grly) and another hypen
func prepareSSH(p *BannerParser, banner string) string {
	firstLine := strings.Split(banner,"\n")[0]
	banner = strings.ReplaceAll(firstLine, "-", " ")
	bannerParts := strings.SplitN(banner, " ", 3)
	if bannerParts[0] == "ssh" && len(bannerParts) > 2 {
		return strings.ReplaceAll(bannerParts[2], "_", " ")
	}
	return ""
}
