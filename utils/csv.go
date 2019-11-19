package utils

import (
	"encoding/csv"
	"fmt"
	"io"
)

type HeadedCSV struct {
	*HeadedCSVOptions
	Reader  *csv.Reader

}

type HeadedCSVOptions struct{
	Name string
	Headers []string
}

// Returns a new headed CSV from a reader source.
// if headers is not nil, it defines the headers from that list. If it is nil, the
// first line of the reader is the headers.
func NewHeadedCSV(source io.Reader, options *HeadedCSVOptions) (headedCSV *HeadedCSV, err error) {
	if options == nil {
		options = &HeadedCSVOptions{}
	}
	csvReader := csv.NewReader(source)
	if options.Headers == nil {
		headers, err := csvReader.Read()
		if err != nil {
			return nil, err
		}
		options.Headers = headers
	}
	if len(options.Name) == 0 {
		options.Name = GenerateRandomString(16)
	}
	headedCSV = &HeadedCSV{
		Reader:  csvReader,
		HeadedCSVOptions: options,
	}
	return
}

// Returns the next row of the CSV.
func (csv *HeadedCSV) NextRow() (map[string]string, error) {
	nextRow, err := csv.Reader.Read()
	if err != nil {
		return nil, err
	}
	if len(nextRow) != len(csv.Headers) {
		return nil, fmt.Errorf("number of fields on line is different than number of headers")
	}
	newMap := make(map[string]string)
	for i, header := range csv.Headers {
		newMap[header] = nextRow[i]
	}
	return newMap, nil
}

func (csv *HeadedCSV) HasHeader(str string) bool {
	if csv.Headers == nil {
		return false
	}
	for _, header := range csv.Headers {
		if header == str {
			return true
		}
	}
	return false
}

func (csv *HeadedCSV) ToArrayMap() (map[string][]string, error) {
	arrayMap := make(map[string][]string)
	for _, header := range csv.Headers {
		arrayMap[header] = make([]string, 0)
	}
	for {
		line, err := csv.NextRow()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		for key, _ := range arrayMap {
			arrayMap[key] = append(arrayMap[key], line[key])
		}
	}
	return arrayMap, nil
}

func (csv *HeadedCSV) ToMapArray() ([]map[string]string, error) {
	mapArray := make([]map[string]string, 0)
	for {
		line, err := csv.NextRow()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		mapArray = append(mapArray, line)
	}
	return mapArray, nil
}
