package logs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// OSROutput defines a file and a path where the file is.
type OSROutput struct {
	io.Writer // object of type writer
	Path string
}

// NewOutput creates a new output with the name and folder path defined in the arguments.
// TODO: Compare with NewLog
func NewOutput(folder, name string) (*OSROutput, error) {
	logsPath, err := getLogsPath()
	if err != nil {
		return nil, err
	}
	folderPath := filepath.Join(logsPath, folder)
	err = os.MkdirAll(folderPath, 0744)
	if err != nil {
		return nil, err
	}
	fileName := filepath.Join(folderPath, createOutputName(name))
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	out := io.MultiWriter(os.Stderr, f)
	return &OSROutput{
		Writer: out,
		Path:   fileName,
	}, nil
}

// GetAttachments returns a list of paths to attach if the object is mailed.
func (output *OSROutput) GetAttachments() []string {
	return []string{output.Path}
}

// Println prints a line into the file.
func (output *OSROutput) Println(text string) {
	_, _ = fmt.Fprintln(output.Writer, text)
}

// createOutputName creates a log name, with the date and time of the log in the name.
func createOutputName(name string) string {
	dateLayout := "2006-01-02_150405"
	return name + "_" + time.Now().Format(dateLayout) + ".txt"
}
