package protocols

import "strings"

func prepareHTTP(p *BannerParser, banner string) string {
	return strings.ReplaceAll(banner, "/", " ")
}
