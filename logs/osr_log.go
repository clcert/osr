package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

// OSRLog defines a custom logger, and adds a path to it.
type OSRLog struct {
	*logrus.Logger // extends from logrus.Logger
	Path string // logfile local path
}

// NewLog creates a new logger with the name and folder path defined in the arguments.
func NewLog(name string) (*OSRLog, error) {
	log := logrus.New()
	logsPath, err := getLogsPath()
	if err != nil {
		return nil, err
	}
	fileName := filepath.Join(logsPath, createLogName(name))
	err = os.MkdirAll(path.Dir(fileName), 0744)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	log.Out = io.MultiWriter(os.Stderr, f)
	return &OSRLog{
		Logger: log,
		Path:   fileName,
	}, nil
}

// GetAttachments returns a list of paths to attach if the object is mailed.
func (log *OSRLog) GetAttachments() []string {
	return []string{log.Path}
}

// createLogName creates a log name, with the date and time of the log in the name.
func createLogName(name string) string {
	dateLayout := "2006-01-02"
	timeLayout := "150405.000"
	return fmt.Sprintf("%s/%s_%s.log", time.Now().Format(dateLayout), time.Now().Format(timeLayout), name)
}
