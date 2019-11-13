package protocols

var ftpRegexes = map[string]FromTo{
	"brackets":   {" ?\\[.*\\]", ""},
	"welcome":    {" ?welcome to", ""},
	"none":       {" ?\\(none\\)", ""}, // Weird thing from GNU FTP server
	"ftpService": {" ?ftp serv((er)|(ice)).*", ""},
	"dow":        {" ?(mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun),? .*", ""},
	"version":    {"version", ""},
	"server":    {"serv((er)|(ice))", ""},
}

var httpRegexes = map[string]FromTo{
	"brackets": {"/", " "}, // Generally slash separes software from version
}

var imapRegexes = map[string]FromTo{
	"ready":    {"( ?is )?ready.*", ""}, // Some IMAP Servers write "ready at DATE"
	"server":   {" ?serv((er)|(ice))", ""},
	"brackets": {" ?\\[.*\\]", ""},
	"lessmore": {" ?<.*>", ""},
	"dow":      {" ?(mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun),? .*", ""},
	"monthSpanish":   {" ?(ene)|(feb)|(mar)|(abr)|(may)|(jun)|(jul)|(ago)|(sep)|(oct)|(nov)|(dic)? .*", ""},
	"monthEngilish":  {" ?(ene)|(feb)|(mar)|(apr)|(may)|(jun)|(jul)|(aug)|(sep)|(oct)|(nov)|(dec)? .*", ""},
}

var smtpRegexes = map[string]FromTo{
	"ready":          {" ?ready.*", ""}, // Some SMTP Servers write "ready at DATE"
	"randomhex":      {" ?\\([0-9a-f]*\\)", ""},
	"dow":            {" ?(mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun),? .*", ""},
	"monthSpanish":   {" ?(ene)|(feb)|(mar)|(abr)|(may)|(jun)|(jul)|(ago)|(sep)|(oct)|(nov)|(dic)? .*", ""},
	"monthEngilish":  {" ?(ene)|(feb)|(mar)|(apr)|(may)|(jun)|(jul)|(aug)|(sep)|(oct)|(nov)|(dec)? .*", ""},
	"esmtpAndBefore": {".* esmtp ", ""},
}

var pop3Regexes = map[string]FromTo{
	"ready":    {"( ?is )?ready.*", ""}, // Some POP3 Servers write "ready at DATE"
	"server":   {" ?serv(er|ice)", ""},
	"pop3":     {" ?pop3", ""},
	"lessmore": {" ?<.*>", ""},
	"dow":      {" ?(mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun),? .*", ""},
	"monthSpanish":   {" ?(ene)|(feb)|(mar)|(abr)|(may)|(jun)|(jul)|(ago)|(sep)|(oct)|(nov)|(dic)? .*", ""},
	"monthEngilish":  {" ?(ene)|(feb)|(mar)|(apr)|(may)|(jun)|(jul)|(aug)|(sep)|(oct)|(nov)|(dec)? .*", ""},
}
