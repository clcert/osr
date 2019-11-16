package protocols

import (
	"fmt"
	"github.com/clcert/osr/panics"
	"regexp"
	"strings"
)

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
	Name             string
	regexes          []FromTo
	StartFormatRegex *regexp.Regexp
	SoftwareRegex    *regexp.Regexp
	VersionRegex     *regexp.Regexp
	ExtraRegexes     []*regexp.Regexp
	Ok               string
	inited           bool
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
		softRegex, err := regexp.Compile(softwareRegex)
		if err != nil {
			panic(panics.Info{
				Text: fmt.Sprintf("Protocol Software Regex is not well defined."),
				Err:  err,
			})
		}
		verRegex, err := regexp.Compile(versionRegex)
		if err != nil {
			panic(panics.Info{
				Text: fmt.Sprintf("Protocol Version Regex is not well defined."),
				Err:  err,
			})
		}
		p.StartFormatRegex = okRegex
		p.SoftwareRegex = softRegex
		p.VersionRegex = verRegex
		p.ExtraRegexes = make([]*regexp.Regexp, len(p.regexes))
		if p.regexes != nil {
			for pos, regex := range p.regexes {
				regex, err := regexp.Compile(regex.From)
				if err != nil {
					panic(panics.Info{
						Text: fmt.Sprintf("Regex in position %d is not well defined: %s", pos, regex),
						Err:  err,
					})
				}
				p.ExtraRegexes[pos] = regex
			}
		}
		p.inited = true
	}
}

// Returns true if the banner is valid
func (p *Parser) IsValid(banner string) bool {
	p.init()
	banner = strings.ToLower(strings.SplitN(banner, "\n", 2)[0])
	return p.StartFormatRegex.MatchString(banner)
}

// Returns the software name and version of the service obtained from the banner.
// If it can find a value, it returns an empty string.
func (p *Parser) GetSoftwareAndVersion(banner string) (software string, version string) {
	p.init()
	// Removing ok string
	banner = strings.ToLower(strings.SplitN(banner, "\n", 2)[0])
	banner = strings.Trim(p.StartFormatRegex.ReplaceAllLiteralString(banner, ""), trimmable)

	for pos, regex := range p.ExtraRegexes {
		banner = regex.ReplaceAllString(banner, p.regexes[pos].To) // length of both arrays is the same
	}
	// parsing software and version
	software = strings.Trim(p.SoftwareRegex.FindString(banner), trimmable)
	version = strings.Trim(p.VersionRegex.FindString(banner), trimmable)

	// If software is null, don't return a version.
	if len(software) == 0 {
		return "", ""
	} else if len(version) > 0 { // remove version from software string
		// split will never be empty because software len is larger than 0
		software = strings.Trim(strings.Split(software, version)[0], trimmable)
	}

	// trim "v" from version code
	version = strings.TrimLeft(version, "v")
	return
}
