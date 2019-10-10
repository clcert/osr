package protocols

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)

const trimmable = " .,;\n\t\r-_="

// BannerParser defines a struct capable to parse version and software from a Banner.
type BannerParser struct {
	regexes       map[string]string
	SoftwareRegex *regexp.Regexp
	VersionRegex  *regexp.Regexp
	ExtraRegexes  map[string]*regexp.Regexp
	prepare       func(*BannerParser, string) string
	StartString   string
	SplitParts    int
	inited        bool
}

// init prepares the parser for its use.
// It should be automatically executed as first method of getVersion and getSoftware.
func (p *BannerParser) init() {
	if !p.inited {
		softRegex, err := regexp.Compile("[a-z0-9.\\-_ ]*")
		if err != nil {
		}
		verRegex, err := regexp.Compile("v?([0-9]+)([\\-+a-z]?\\.[\\-+a-z]?[0-9]+)+")
		if err != nil {
		}
		p.SoftwareRegex = softRegex
		p.VersionRegex = verRegex
		p.ExtraRegexes = make(map[string]*regexp.Regexp)
		if p.regexes != nil {
			for regexType, regex := range p.regexes {
				regex, err := regexp.Compile(regex)
				if err != nil {
					panic(panics.Info{
						Text: fmt.Sprintf("%s Regex is not well defined: %s", regexType, regex),
						Err:  err,
					})
				}
				p.ExtraRegexes[regexType] = regex
			}
		}
	}
	p.inited = true
}

// Returns true if the banner is valid
func (p *BannerParser) IsValid(banner string) bool {
	p.init()
	banner = strings.ToLower(strings.SplitN(banner, "\n", 2)[0])
	return len(banner) >= len(p.StartString) && strings.HasPrefix(banner, p.StartString)
}

func (p *BannerParser) Prepare(banner string) string {
	return p.prepare(p, banner)
}

// Returns the software name and version of the service obtained from the banner.
// If it can find a value, it returns an empty string.
func (p *BannerParser) GetSoftwareAndVersion(banner string) (software string, version string) {
	p.init()
	banner = strings.ToLower(banner)
	if p.prepare != nil {
		banner = p.Prepare(banner)
	}

	// parsing software
	software = strings.Trim(p.SoftwareRegex.FindString(banner), trimmable)

	// parsing version
	versionSlice := strings.SplitN(banner, " ", p.SplitParts)
	if versionSlice == nil || len(versionSlice) < p.SplitParts {
		return
	}
	version = strings.Trim(p.VersionRegex.FindString(versionSlice[p.SplitParts-1]), trimmable)

	// Removing version part from software
	if len(version) > 0 {
		software = strings.Trim(strings.Split(software, p.SoftwareRegex.FindString(version))[0], trimmable)
	}
	if len(software) == 0 {
		return "", ""
	}
	// removing v from versions (some servers have it, some servers not)
	version = strings.TrimPrefix(version, "v")
	return
}
