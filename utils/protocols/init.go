package protocols
var Parsers ParserMap

func init() {
	Parsers.Add(
		&Parser{
			Name:       "http",
			regexes:    httpRegexes,
			SplitParts: 2,
		},
		&Parser{
			Name:       "ssh",
			Ok:         "^ssh",
			SplitParts: 2,
		},
		&Parser{
			Name:       "imap",
			Ok:         "^\\* ((ok)|(bye)|(no))",
			regexes:    imapRegexes,
			SplitParts: 3,
		},
		&Parser{
			Name:       "ftp",
			Ok:         "^\\d{3}",
			regexes:    ftpRegexes,
			SplitParts: 2,
		},
		&Parser{
			Name:       "pop3",
			Ok:         "^((\\+ok)|(-err))",
			regexes:    pop3Regexes,
			SplitParts: 2,
		},
		&Parser{
			Name:       "smtp",
			Ok:         "^\\d{3}",
			regexes:    smtpRegexes,
			SplitParts: 2,
		},
	)
	// parser from https is equal to parser from http
	Parsers["https"] = Parsers["http"]
}
