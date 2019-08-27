package logs

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

// OSRLog defines a custom logger, and adds a path to it.
type OSRLog struct {
	*logrus.Logger
	Path string
}

// Creates a new logger with the name and folder path defined in the arguments.
func NewLog(name string) (*OSRLog, error) {
	log := logrus.New()
	logsPath, err := getLogsPath()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(logsPath, 0744)
	if err != nil {
		return nil, err
	}
	fileName := filepath.Join(logsPath, createLogName(name))
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

// Returns a list of paths to attach if the object is mailed.
func (log *OSRLog) GetAttachments() []string {
	return []string{log.Path}
}

// Creates a log name, with the date and time of the log in the name.
func createLogName(name string) string {
	dateLayout := "2006-01-02_150405.000"
	return name + "_" + time.Now().Format(dateLayout) + ".log"
}
