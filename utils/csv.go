package utils

import (
	"encoding/csv"
	"fmt"
	"io"
)

type HeadedCSV struct {
	Reader  *csv.Reader
	Headers []string
}

func NewHeadedCSV(source io.Reader, headers []string) (headedCSV *HeadedCSV, err error) {
	if headers == nil {
		headers, err = GetCSVHeader(source)
		if err != nil {
			return
		}
	}
	csvReader := csv.NewReader(source)
	csvReader.ReuseRecord = true
	headedCSV = &HeadedCSV{
		Reader:  csvReader,
		Headers: headers,
	}
	return
}

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


func GetCSVHeader(file io.Reader) ([]string, error) {
	csvReader := csv.NewReader(file)
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	return header, nil
}
