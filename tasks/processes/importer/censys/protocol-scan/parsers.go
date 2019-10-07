package port_scan

import (
	"encoding/json"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/utils/censys"
	"net"
	"time"
)

type Parser func(saver savers.Saver, line string, options *ParserOptions) error

type ParserOptions struct {
	DefaultDate time.Time
	Port        uint16
	Protocol    string
}

const DateFormat = "2006-01-02T15:04:05-07:00"

var parsers = map[string]Parser{
	"http":  httpParser,
	"https": httpParser,
	"smtp":  smtpParser,
	"imap":  imapParser,
	"pop3":  pop3Parser,
	"ftp":   ftpParser,
	"ssh":   sshParser,
}

func getIPAndDate(entry censys.BasicEntry, options *ParserOptions) (t time.Time, ip net.IP, err error) {
	t, err = time.Parse(entry.Timestamp, DateFormat)
	if err != nil {
		t = options.DefaultDate
	}
	return t, net.ParseIP(entry.IP), nil
}

func httpParser(saver savers.Saver, line string, options *ParserOptions) error {
	var entry censys.HTTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return err
	}
	ip, date, err := getIPAndDate(entry.BasicEntry, options)
	if err != nil {
		return err
	}
}

func smtpParser(saver savers.Saver, line string, options *ParserOptions) error {
	var entry censys.SMTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return err
	}
}

func imapParser(saver savers.Saver, line string, options *ParserOptions) error {
	var entry censys.IMAPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return err
	}
}

func pop3Parser(saver savers.Saver, line string, options *ParserOptions) error {
	var entry censys.POP3Entry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return err
	}
}

func ftpParser(saver savers.Saver, line string, options *ParserOptions) error {
	var entry censys.FTPEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return err
	}
}

func sshParser(saver savers.Saver, line string, options *ParserOptions) error {
	var entry censys.SSHEntry
	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		return err
	}
}
