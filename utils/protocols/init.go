package protocols

var Parsers map[string]*BannerParser

func init() {
	Parsers = map[string]*BannerParser{
		"http": {
			prepare:    prepareHTTP,
			SplitParts: 2,
		},
		"ssh": {
			prepare:     prepareSSH,
			StartString: "ssh",
			SplitParts:  2, // ssh <something> rest... (rest could contain a version)
		},
		"ftp": {
			regexes:     ftpRegexes,
			prepare:     prepareFTP,
			StartString: "220",
			SplitParts:  2, // 220 <something> rest... (rest could contain a version)
		},
		"pop3": {
			regexes:     pop3Regexes,
			prepare:     preparePOP3,
			StartString: "+ok",
			SplitParts:  2, // +ok <something> rest... (rest could contain a version)
		},
		"imap": {
			regexes:     imapRegexes,
			prepare:     prepareIMAP,
			StartString: "* ok",
			SplitParts:  3, // * ok <something> rest... (rest could contain a version)
		},
		"smtp": {
			regexes:     smtpRegexes,
			prepare:     prepareSMTP,
			StartString: "220",
			SplitParts:  2, // 220 <something> rest... (rest could contain a version)
		},
	}
	// parser from https is equal to parser from http
	Parsers["https"] = Parsers["http"]
}
