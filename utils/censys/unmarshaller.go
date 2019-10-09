package censys

import (
	"encoding/json"
	"fmt"
	"github.com/clcert/osr/utils/protocols"
	"time"
)

type Unmarshaler func(line string) (protocols.Entry, error)

type ParserOptions struct {
	DefaultDate time.Time
	Port        uint16
	Protocol    string
}

func (o *ParserOptions) Unmarshal(line string) (protocols.Entry, error) {
	if o == nil {
		return nil, fmt.Errorf("option cannot be blank")
	}
	if _, ok := Unmarshalers[o.Protocol]; !ok {
		return nil, fmt.Errorf("parser not found for protocol %s", o.Protocol)
	}
	return Unmarshalers[o.Protocol](line)
}

const DateFormat = "2006-01-02T15:04:05-07:00"

var Unmarshalers = map[string]Unmarshaler{
	"http":  unmarshalHTTP,
	"https": unmarshalHTTP,
	"smtp":  unmarshalSMTP,
	"imap":  unmarshalIMAP,
	"pop3":  unmarshalPOP3,
	"ftp":   unmarshalFTP,
	"ssh":   unmarshalSSH,
}

func unmarshalHTTP(line string) (protocols.Entry, error) {
	var entry HTTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalSMTP(line string) (protocols.Entry, error) {
	var entry SMTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalIMAP(line string) (protocols.Entry, error) {
	var entry IMAPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalPOP3(line string) (protocols.Entry, error) {
	var entry POP3Entry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalFTP(line string) (protocols.Entry, error) {
	var entry FTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func unmarshalSSH(line string) (protocols.Entry, error) {
	var entry SSHEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
