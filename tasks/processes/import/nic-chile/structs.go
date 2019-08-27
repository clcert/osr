package nic_chile

import "io"

type NamedReader interface {
	io.Reader
	Name() string
}

type GenericNamedReader struct {
	io.Reader
	NameString string
}

func (namedReader GenericNamedReader) Name() string {
	return namedReader.NameString
}
