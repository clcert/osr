package sources

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"
)

// HTTPConfig defines a configuration for a HTTP Source
type HTTPConfig struct {
	URL      string              // Initial URL to download
	Filename string              // Name for initial URL
	Username string              // If defined, BASIC username used to crawl the file
	Password string              // If defined, BASIC password used to crawl the file
	Method   string              // Method to use when retrieving the first file. It can be GET or POST
	Body     map[string][]string // Body to send as body with the first file.
	Filter   *FilterConfig       // Filter configuration
}

// HTTPource defines a remote source of files, connected via HTTP.
type HTTPSource struct {
	*HTTPConfig           // Configuration
	name   string         // Source name
	filter *Filter        // Filter
	files  chan *HTTPFile // Channel of files
	log    *logs.OSRLog   // Source Log
	params utils.Params   // Process Parameters
	client *http.Client   // Client used in all requests
}

// HTTPFile represents a HTTP entry source.
type HTTPFile struct {
	source   *HTTPSource    // Source related to the entry
	name     string         // Name of file
	url      *url.URL       // Complete file URL
	body     url.Values     // Body to send when HTTPFile is retrieved
	response *http.Response // Response body
	buffer   io.Reader      // Buffered extReader
	method   string         // Method to use when retrieving file
}

// Creates a new SFTP source from a configuration
func (config *HTTPConfig) New(name string, params utils.Params) (source *HTTPSource, err error) {
	if config.Filter == nil {
		config.Filter = &FilterConfig{}
	}
	err = config.Format(params)
	if err != nil {
		return
	}
	var filter *Filter
	filter, err = config.Filter.New()
	if err != nil {
		return
	}
	log, err := logs.NewLog(name)
	if err != nil {
		return
	}

	client := &http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.Jar = jar
	source = &HTTPSource{
		name:       name,
		HTTPConfig: config,
		files:      make(chan *HTTPFile),
		filter:     filter,
		log:        log,
		params:     params,
		client:     client,
	}
	return
}

// Format formats the configuration, using the params defined in the task
func (config *HTTPConfig) Format(params utils.Params) error {
	if params == nil {
		// do nothing
		return nil
	}
	config.URL = params.FormatString(config.URL)
	config.Username = params.FormatString(config.Username)
	config.Password = params.FormatString(config.Password)
	config.Method = params.FormatString(config.Method)
	config.Filter = config.Filter.Format(params)
	return nil
}

func (source *HTTPSource) Init() error {
	id, err := source.GetID()
	if err != nil {
		id = fmt.Sprintf("cannot get id: %s", err)
	}
	source.log.WithFields(logrus.Fields{
		"type":   "HTTP",
		"id":     id,
		"filter": source.filter,
	}).Info("Starting HTTP Source...")
	if source.URL == "" {
		return fmt.Errorf("config mandatory fields not initialized")
	}
	aUrl, err := url.Parse(source.URL)
	if err != nil {
		return err
	}

	rootFile := &HTTPFile{
		source: source,
		name:   source.Filename,
		url:    aUrl,
		body:   source.Body,
		method: source.Method,
	}
	go func() {
		if source.filter == nil || !source.filter.Recursive {
			source.log.WithFields(logrus.Fields{
				"type": "HTTP",
				"id":   id,
			}).Info("Adding page to channel and exiting")
			source.files <- rootFile
		} else {
			source.log.WithFields(logrus.Fields{
				"type": "HTTP",
				"id":   id,
			}).Info("Recursive enabled. Adding files in links...")
			source.retrieveFiles(rootFile)
		}
		source.log.WithFields(logrus.Fields{
			"type": "HTTP",
			"id":   id,
		}).Info("Closing channel...")
		close(source.files)
	}()
	return nil
}

func (source *HTTPSource) GetID() (string, error) {
	return source.URL, nil
}

func (source *HTTPSource) Next() Entry {
	entry, ok := <-source.files
	if !ok {
		return nil
	}
	return entry
}

func (source *HTTPSource) Close() error {
	return nil
}

func (source *HTTPSource) GetName() string {
	return source.name
}

func (source *HTTPSource) GetAttachments() []string {
	return []string{source.log.Path}
}

func (source *HTTPSource) retrieveFiles(httpFile *HTTPFile) {
	// TODO: to this recursively some day
	// Open file and scan for links
	f, err := httpFile.Open()
	if err != nil {
		// TODO: log this
		return
	}
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		// TODO: log this
		return
	}
	visited := make(map[string]struct{})
	// We get the response location, in case we were redirected
	rootURL := httpFile.response.Request.URL
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, hasLink := s.Attr("href")
		if hasLink {
			innerText := s.Text()
			aURL, err := rootURL.Parse(href)
			if err != nil {
				// TODO log this
				return
			}
			// Skip if the url returns back into the tree
			if _, ok := visited[aURL.String()]; ok {
				// TODO: log this
				return
			}
			// mark url as visited
			visited[aURL.String()] = struct{}{}
			aFile := &HTTPFile{
				name:   innerText,
				source: source,
				url:    aURL,
				method: "GET",
			}
			if len(source.filter.Patterns) == 0 {
				source.log.WithFields(logrus.Fields{
					"path": aFile.url.String(),
				}).Info("Adding file...")
				source.files <- aFile
				return
			}
			for _, regex := range source.filter.Patterns {
				if regex.MatchString(aFile.Name()) {
					source.log.WithFields(logrus.Fields{
						"path": aFile.Path(),
					}).Info("Adding file...")
					source.files <- aFile
					return
				}
			}
		}
	})
}

func (srcFile *HTTPFile) Open() (reader io.Reader, err error) {
	if srcFile.response == nil {
		var req *http.Request
		req, err = http.NewRequest(
			srcFile.source.Method,
			srcFile.url.String(),
			strings.NewReader(srcFile.body.Encode()),
		)
		if err != nil {
			return
		}
		if srcFile.body != nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		if srcFile.source.Username != "" && srcFile.source.Password != "" {
			req.SetBasicAuth(srcFile.source.Username, srcFile.source.Password)
		}
		var file *http.Response
		file, err = srcFile.source.client.Do(req)
		if err != nil {
			return
		}
		srcFile.response = file
	}
	if srcFile.buffer == nil {
		srcFile.buffer = bufio.NewReader(srcFile.response.Body)
	}
	reader = srcFile.buffer
	return

}

func (srcFile *HTTPFile) Name() string {
	if srcFile.name == "" {
		return path.Base(srcFile.url.String())
	} else {
		return srcFile.name
	}
}

func (srcFile *HTTPFile) Dir() string {
	return path.Dir(srcFile.url.String())
}

func (srcFile *HTTPFile) Path() string {
	return srcFile.url.String()
}

func (srcFile *HTTPFile) Close() error {
	if srcFile.response != nil {
		err := srcFile.response.Body.Close()
		srcFile.response = nil
		return err
	}
	return nil
}
