package utils

import (
	"bufio"
	"io"
	"strings"
)

type MercuryFile struct {
	Name   string
	Reader *bufio.Reader
}

func MercuryFromReader(reader io.Reader, name string) *MercuryFile {
	bufReader := bufio.NewReader(reader)
	return &MercuryFile{
		Name:   name,
		Reader: bufReader,
	}
}

func (f MercuryFile) NextResult() (string, error) {
	result, err := f.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}
