// Sources are abstractions of providers of raw information for scans.
// Currently, we have three types of sources: Commands, SFTP and HTTP
//
// The Script source connects to a server already registered and executes a command, associating
// its standard output to a virtual file.
//
// The SFTP source connects to an already registered server in SFTP mode.
//
// The HTTP source downloads (recursively or not) an HTTP webpage, the page could be
// authenticated vía BASIC Auth.
package sources

import (
	"fmt"
	"github.com/clcert/osr/mailer"
	"github.com/clcert/osr/utils"
	"io"
	"regexp"
	"strings"
)

// It defines a source config. It can define only one from [SFTP, HTTP, Script, ...]
type Config struct {
	SFTP     *SFTPConfig      // Config if type is sftp
	HTTP     *HTTPConfig      // Config if type is http
	Script   *ScriptConfig    // Config if type is command
	Query    *QueryListConfig // Config if type is query
}

// Source defines a stream of entries, related to an import.
// The files should be similar on each source, but it's not obligatory
// When a file is asked for, it disappears, so if a process needs to access them sequentially,
// it needs to store them somewhere (an array of entries, for example).
type Source interface {
	mailer.Attachable // Attachments to mailer
	// Init inits the source, connecting to it and starting the retrieval of entries
	Init() error
	// GetName returns a given name of the source, based on task and process association. Two instances with the same source conf
	// could have different IDs.
	GetName() string
	// GetID returns an ID for the source. Two instances with the same source conf have the same ID
	GetID() (string, error)
	// Next returns the next entry element, or nil if the stream is over.
	Next() Entry
	// Close closes and cleans the stream.
	Close() error
}

// An entry represents a single file/object in a source.
// It has a name, a Path, a dir and it's openable and closeable.
type Entry interface {
	// Open tries to open the entry, returning a extReader if it is successful.
	Open() (io.Reader, error)
	// Name returns the name of the file in the entry.
	Name() string
	// Path returns the Path of the file in the entry.
	Path() string
	// Dir returns the dir of the file in the entry.
	Dir() string
	// Close closes the entry.
	Close() error
}

// A FilterConfig defines specific filter options, like regex matching and recursiveness.
type FilterConfig struct {
	Patterns  []string // List of regexes
	Recursive bool     // True if recursive
}

// A filter defines a set of rules to add a file to the entries list.
type Filter struct {
	Patterns  []*regexp.Regexp // A list of regexes. If a file has a name that matches them, it's added to the list.
	Recursive bool             // If true, some sources look onto internal folders of a root Path.
}

// New creates a new Source based on a specific configuration.
// It also sets a name to the source, for logging purposes.
func (source *Config) New(name string, params utils.Params) (Source, error) {
	switch {
	case source.SFTP != nil:
		return source.SFTP.New(name, params)
	case source.HTTP != nil:
		return source.HTTP.New(name, params)
	case source.Script != nil:
		return source.Script.New(name, params)
	case source.Query != nil:
		return source.Query.New(name, params)
	default:
		return nil, fmt.Errorf("invalid source: no config defined (you need to define sftp, http, query or command config)")
	}
}

func (filter *FilterConfig) Format(params utils.Params) *FilterConfig {
	newFilterConfig := &FilterConfig{
		Patterns:  make([]string, len(filter.Patterns)),
		Recursive: filter.Recursive,
	}
	if filter.Patterns != nil {
		for i, pattern := range filter.Patterns {
			newFilterConfig.Patterns[i] = params.FormatString(pattern)
		}
	}
	return newFilterConfig
}

// New creates a new Filter based on a filter configuration, or returns an error if the configuration is invalid.
func (filter *FilterConfig) New() (*Filter, error) {
	patterns := make([]*regexp.Regexp, 0)
	for i, pattern := range filter.Patterns {
		regexpPattern, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid filter: couldn't compile regexp rule number %d", i)
		}
		patterns = append(patterns, regexpPattern)
	}
	return &Filter{
		Recursive: filter.Recursive,
		Patterns:  patterns,
	}, nil
}

// Matches checks if an entry complies with a simple pattern matching technique:
// We match a string that could have an asterisk "*", and match the part before the asterisk as prefix and the part after
// the asterisk as suffix. If it matches, we return true, and if not, we return false.
func Matches(entry Entry, name string) bool {
	if entry == nil {
		return false
	}
	splittedName := strings.Split(name, "*")
	if len(splittedName) == 1 {
		if entry.Name() == name {
			return true
		}
	} else if len(splittedName) == 2 {
		if strings.HasPrefix(entry.Name(), splittedName[0]) && strings.HasSuffix(entry.Name(), splittedName[1]) {
			return true
		}
	}
	return false
}

// Get consumes a source until it gets a file with a name that matches to the provided pattern.
func Get(source Source, name string) Entry {
	for {
		next := source.Next()
		if next == nil {
			return nil
		} else if Matches(next, name) {
			return next
		}
	}
}
