package protocols

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)

const trimmable = " .,;\n\t\r-_="

// FromTo Defines a pair of strings.
// The first string represents a regexp to find and the second one the replacement.
type FromTo struct {
	From, To string
}

// ParserMap is a map of parsers
type ParserMap map[string]*Parser

func (m ParserMap) Add(parsers ...*Parser) {
	for _, parser := range parsers {
		m[parser.Name] = parser
	}
}

// Parser defines a struct capable to parse version and software from a Banner.
type Parser struct {
	Name          string
	regexes       map[string]FromTo
	OkRegex       *regexp.Regexp
	SoftwareRegex *regexp.Regexp
	VersionRegex  *regexp.Regexp
	ExtraRegexes  map[string]*regexp.Regexp
	Ok            string
	SplitParts    int
	inited        bool
}

// init prepares the parser for its use.
// It should be automatically executed as first method of getVersion and getSoftware.
func (p *Parser) init() {
	if !p.inited {
		okRegex, err := regexp.Compile(p.Ok)
		if err != nil {
			panic(panics.Info{
				Text: fmt.Sprintf("Protocol OK Regex is not well defined: %s", p.Ok),
				Err:  err,
			})
		}
		softRegex, err := regexp.Compile("[a-z0-9.\\-_ ]*")
		if err != nil {
			panic(panics.Info{
				Text: fmt.Sprintf("Protocol Software Regex is not well defined."),
				Err:  err,
			})
		}
		verRegex, err := regexp.Compile("v?([0-9]+)([\\-+a-z]?\\.[\\-+a-z]?[0-9]+)+")
		if err != nil {
			panic(panics.Info{
				Text: fmt.Sprintf("Protocol Version Regex is not well defined."),
				Err:  err,
			})
		}
		p.OkRegex = okRegex
		p.SoftwareRegex = softRegex
		p.VersionRegex = verRegex
		p.ExtraRegexes = make(map[string]*regexp.Regexp)
		if p.regexes != nil {
			for regexType, regex := range p.regexes {
				regex, err := regexp.Compile(regex.From)
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
func (p *Parser) IsValid(banner string) bool {
	p.init()
	banner = strings.ToLower(strings.SplitN(banner, "\n", 2)[0])
	return p.OkRegex.MatchString(banner)
}
// Returns the software name and version of the service obtained from the banner.
// If it can find a value, it returns an empty string.
func (p *Parser) GetSoftwareAndVersion(banner string) (software string, version string) {
	p.init()
	banner = strings.ToLower(banner)
	// Removing OK string
	banner = strings.Split(banner, "\n")[0]
	banner = p.OkRegex.ReplaceAllLiteralString(banner, "")
	for regexKey, regex := range p.ExtraRegexes {
		banner = regex.ReplaceAllString(banner, p.regexes[regexKey].To)
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
