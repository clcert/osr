package sources

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// TODO: More loging on this source
// ZipConfig defines the configuration that a Zip Source uses
type ZipConfig struct {
	Path   string        // Local Path to the zip file
	Filter *FilterConfig // Filter configuration
}

// ZipSource represents a Zip file.
type ZipSource struct {
	*ZipConfig             // Script Configuration
	name      string       // Source name
	extReader io.Reader    // External reader to use instead of config file
	zipReader *zip.Reader  // Zip reader
	filter    *Filter      // Source filter
	files     chan Entry   // Files channel
	log       *logs.OSRLog // Source log
	params    utils.Params // Process Parameters
}

// ZipEntry represents a file inside the Zip
type ZipEntry struct {
	source *ZipSource
	file   *zip.File
	buffer io.Reader
}

// New creates a new ZipSource from a ZipConfig.
func (config *ZipConfig) New(name string, params utils.Params) (source *ZipSource, err error) {
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
	source = &ZipSource{
		name:      name,
		ZipConfig: config,
		files:     make(chan Entry),
		filter:    filter,
		log:       log,
		params:    params,
	}
	return
}

// Format formats the configuration, using the params defined in the task
func (config *ZipConfig) Format(params utils.Params) error {
	if params == nil {
		// do nothing
		return nil
	}
	config.Path = params.FormatString(config.Path)
	config.Filter = config.Filter.Format(params)
	return nil
}

// NewFrom creates a new ZipSource, but using an specific reader and closer.
// The reader points to a zip content.
// This form of creation is useful if we want to create a ZipSource from an entry.
func (config *ZipConfig) NewFrom(name string, reader io.Reader, params utils.Params) (source *ZipSource, err error) {
	source, err = config.New(name, params)
	if err != nil {
		return
	}
	source.extReader = reader
	return
}

// newReaderAt transforms a Reader into a ReaderAt, and returns it and its length.
// It returns an error if it can read all the data of the reader.
// This reads the file completely so use with caution and small files.
func (source ZipSource) newReaderAt(reader io.Reader) (io.ReaderAt, int64, error) {
	fLen, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, -1, err
	}
	return bytes.NewReader(fLen), int64(len(fLen)), nil
}

func (source *ZipSource) Init() (err error) {
	if source.Path == "" && source.extReader == nil {
		err = fmt.Errorf("mandatory config fields not initialized: you must initialize Path in config or extReader with NewFromEntry()")
		return
	}
	var reader io.Reader
	if source.extReader == nil {
		var f *os.File
		f, err = os.Open(source.Path)
		if err != nil {
			return err
		}
		reader = f
	} else {
		reader = source.extReader
	}
	readerAt, fLen, err := source.newReaderAt(reader)
	if err != nil {
		return err
	}
	zipReader, err := zip.NewReader(readerAt, fLen)
	if err != nil {
		return err
	}
	source.zipReader = zipReader
	go func() {
		for _, f := range zipReader.File {
			if strings.Contains(f.Name, "/") && !source.filter.Recursive {
				continue
			}
			aFile := &ZipEntry{
				source: source,
				file:   f,
			}
			if len(source.filter.Patterns) == 0 {
				source.files <- aFile
			} else {
				for _, regex := range source.filter.Patterns {
					if regex.MatchString(f.Name) {
						source.files <- aFile
						break
					}
				}
			}
		}
		close(source.files)
	}()

	return nil
}

func (source *ZipSource) Next() Entry {
	entry, open := <-source.files
	if !open {
		return nil
	}
	return entry
}

func (source *ZipSource) Close() error {
	return nil
}

func (source *ZipSource) GetID() (string, error) {
	if source.zipReader == nil {
		return "", fmt.Errorf("source not initialized")
	}

	id := source.Path
	return id, nil
}

func (source *ZipSource) GetName() string {
	return source.name
}

func (source *ZipSource) GetAttachments() []string {
	return []string{source.log.Path}
}

func (srcFile *ZipEntry) Open() (reader io.Reader, err error) {
	if srcFile.buffer == nil {
		var fileReader io.Reader
		fileReader, err = srcFile.file.Open()
		if err != nil {
			return
		}
		srcFile.buffer = bufio.NewReader(fileReader)
	}
	return srcFile.buffer, nil
}

func (srcFile *ZipEntry) Name() string {
	return path.Base(srcFile.file.Name)
}

func (srcFile *ZipEntry) Dir() string {
	return path.Dir(srcFile.file.Name)
}

func (srcFile *ZipEntry) Path() string {
	return srcFile.file.Name
}

func (srcFile *ZipEntry) Close() error {
	// Nothing to do in this case
	return nil
}
