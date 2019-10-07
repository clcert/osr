package port_scan

import (
	"encoding/json"
	"fmt"
	"github.com/clcert/osr/utils/censys"
	"time"
)

type Unmarshaler func(line string) (censys.Entry, error)

type ParserOptions struct {
	DefaultDate time.Time
	Port        uint16
	Protocol    string
}

func (o *ParserOptions) Unmarshal(line string) (censys.Entry, error) {
	if o == nil {
		return nil, fmt.Errorf("option cannot be blank")
	}
	if _, ok := unmarshalers[o.Protocol]; !ok {
		return nil, fmt.Errorf("parser not found for protocol %s", options.Protocol)
	}
	return unmarshalers[o.Protocol](line)
}

const DateFormat = "2006-01-02T15:04:05-07:00"

var unmarshalers = map[string]Unmarshaler{
	"http":  unmarshalHTTP,
	"https": unmarshalHTTP,
	"smtp":  unmarshalSMTP,
	"imap":  unmarshalIMAP,
	"pop3":  unmarshalPOP3,
	"ftp":   unmarshalFTP,
	"ssh":   unmarshalSSH,
}

func unmarshalHTTP(line string) (censys.Entry, error) {
	var entry censys.HTTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalSMTP(line string) (censys.Entry, error) {
	var entry censys.SMTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalIMAP(line string) (censys.Entry, error) {
	var entry censys.IMAPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalPOP3(line string) (censys.Entry, error) {
	var entry censys.POP3Entry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalFTP(line string) (censys.Entry, error) {
	var entry censys.FTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalSSH(line string) (censys.Entry, error) {
	var entry censys.SSHEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
