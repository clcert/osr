package protocol_parser

// BannerParser defines a struct capable to parse version and software from a Banner.
type BannerParser interface {
	// Returns the version of the service obtained from the banner
	GetVersion(string) string
	// Returns the software name of the service obtained from the banner
	GetSoftware(string) string
}

// BannerParsers represent the map with all registered parsers
var BannerParsers = make(map[string]BannerParser)

// init registers the bannerParsers.
func init() {
	BannerParsers["http"] = &HTTPParser{}
	BannerParsers["https"] = BannerParsers["http"]
	BannerParsers["ssh"] = &SSHParser{}
	BannerParsers["ftp"] = &FTPParser{}
	BannerParsers["pop3"] = &POP3Parser{}
	BannerParsers["imap"] = &IMAPParser{}
	BannerParsers["smtp"] = &SMTPParser{}

}

