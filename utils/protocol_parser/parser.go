package protocol_parser

// BannerParser defines a struct capable to parse version and software from a Banner.
type BannerParser interface {
	// Returns true if the banner is valid
	IsValid(string) bool
	// Returns the version of the service obtained from the banner
	GetVersion(string) string
	// Returns the software name of the service obtained from the banner
	GetSoftware(string) string
}

