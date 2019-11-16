package protocols

var Parsers  = make(ParserMap)

func init() {
	Parsers.Add(
		&Parser{
			Name:       "http",
			regexes:    httpRegexes,
		},
		&Parser{
			Name:       "ssh",
			Ok:         "^ssh-\\d.\\d+-",
		},
		&Parser{
			Name:       "imap",
			Ok:         "^\\* ((ok)|(bye)|(no))",
			regexes:    imapRegexes,
		},
		&Parser{
			Name:       "ftp",
			Ok:         "^\\d{3}",
			regexes:    ftpRegexes,
		},
		&Parser{
			Name:       "pop3",
			Ok:         "^((\\+ok)|(-err))",
			regexes:    pop3Regexes,
		},
		&Parser{
			Name:       "smtp",
			Ok:         "^\\d{3}",
			regexes:    smtpRegexes,
		},
	)
	// parser from https is equal to parser from http
	Parsers["https"] = Parsers["http"]
}
