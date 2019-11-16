package protocols

// Regexes

const trimmable = " .,;\n\t\r-_="
const softwareRegex = "[a-z][a-z0-9.\\-_ ]+" // starts with letter
const versionRegex = "v?([0-9]+)([\\-+a-z]?\\.[a-z0-9]+)+"

// Common Regexes
var (
	dowEnglish     = FromTo{" (mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun),? .*", ""}
	dowSpanish     = FromTo{" (lun)|(mar)|(mie)|(jue)|(vie)|(s[aá]b)|(dom),? .*", ""}
	monthEnglish   = FromTo{" (jan)|(feb)|(mar)|(apr)|(may)|(jun)|(jul)|(aug)|(sep)|(oct)|(nov)|(dec),? .*", ""}
	monthSpanish   = FromTo{" (ene)|(feb)|(mar)|(abr)|(may)|(jun)|(jul)|(ago)|(sep)|(oct)|(nov)|(dic),? .*", ""}
	localhost      = FromTo{" ?localhost", ""}
	domain         = FromTo{"^ ?([a-z0-9áéíóúñ\\-]*)(\\.[a-z0-9áéíóúñ\\-]+)+\\.?", ""}
	serviceServer  = FromTo{" ?serv((er)|(ice)|(idor)|(icio))", ""}
	welcomeTo      = FromTo{" ?((welcome to)|(bienvenido al?))", ""}
	version        = FromTo{" ?((version)|(release))", ""}
	ready          = FromTo{" ?(is )?ready.*", ""}
	allInBrackets  = FromTo{" ?\\[.*\\]", ""}
	hexIPBrackets  = FromTo{" ?\\[[0-9a-f.:]*\\]", ""}
	hexParenthesis = FromTo{" ?\\([0-9a-f]*\\)", ""}
	lessmore       = FromTo{" ?<.*>", ""}
)

var ftpRegexes = []FromTo{
	{" ?\\(none\\)", ""},           // Weird thing from GNU FTP server
	{"\\(", ""},                    // Sometimes the software and version is inside parenthesis
	{"\\)", ""},                    // Sometimes the software and version is inside parenthesis
	{" ?ftp serv((er)|(ice))", ""}, // XXX ftp server or XXX ftp service (and the rest is useless)
	allInBrackets,
	version,
	welcomeTo,
	serviceServer,
	dowEnglish,
	dowSpanish,
	monthSpanish,
	monthEnglish,
	domain,
	localhost,
}

var httpRegexes = []FromTo{
	{"/", " "}, // Generally slash separes software from version
}

var imapRegexes = []FromTo{
	{" ?imap\\d?([a-z0-9]*)", ""}, // IMAP4
	allInBrackets,
	lessmore,
	ready,
	version,
	welcomeTo,
	serviceServer,
	dowEnglish,
	dowSpanish,
	monthSpanish,
	monthEnglish,
	domain,
	localhost,
}

var smtpRegexes = []FromTo{
	{" ?esmtp", ""}, // esmtp
	hexIPBrackets,
	hexParenthesis,
	lessmore,
	ready,
	version,
	welcomeTo,
	serviceServer,
	dowEnglish,
	dowSpanish,
	monthSpanish,
	monthEnglish,
	domain,
	localhost,
}

var pop3Regexes = []FromTo{
	{" ?pop3", ""},
	hexIPBrackets,
	hexParenthesis,
	lessmore,
	ready,
	version,
	welcomeTo,
	serviceServer,
	dowEnglish,
	dowSpanish,
	monthSpanish,
	monthEnglish,
	domain,
	localhost,
}
