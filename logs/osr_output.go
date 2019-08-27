package logs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// OSROutput defines a file and a path where it is.
type OSROutput struct {
	io.Writer
	Path string
}

// Creates a new output with the name and folder path defined in the arguments.
// TODO: Join with NewLog
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

// Returns a list of paths to attach if the object is mailed.
func (output *OSROutput) GetAttachments() []string {
	return []string{output.Path}
}

// Prints a line into the file.
func (output *OSROutput) Println(text string) {
	_, _ = fmt.Fprintln(output.Writer, text)
}

// Creates a log name, with the date and time of the log in the name.
func createOutputName(name string) string {
	dateLayout := "2006-01-02_150405"
	return name + "_" + time.Now().Format(dateLayout) + ".txt"
}
