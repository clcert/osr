package models

import (
	"net"
	"time"
)

// DnsRRModel contains the metainformation related to the respective model.
var DnsRRModel = Model{
	Name:        "DNS RR",
	Description: "DNS Resource Record",
	StructType:  &DnsRR{},
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS dns_rr_index ON ?TableName USING gist (ip_value inet_ops)",
		"CREATE INDEX IF NOT EXISTS dns_rr_timestamp ON ?TableName USING btree (date)",
		"SELECT partman.create_parent('public.dns_rrs', 'date', 'native', 'weekly');",
	}}

// RRType maps a number to an RR type.
type RRType int

// The RRs scanned at the time are A, MX, NS and CNAME.
const (
	NORR  RRType = iota // Undefined RR
	A                   // A RR
	MX                  // MX RR
	NS                  // NS RR
	CNAME               // CNAME RR
)

var stringToRRType = map[string]RRType{
	"a":     A,
	"mx":    MX,
	"ns":    NS,
	"cname": CNAME,
}

// A simple map which transforms a string RR multiType to a const.
func StringToRRType(s string) RRType {
	if t, ok := stringToRRType[s]; ok {
		return t
	}
	return NORR
}

// A simple map which transforms a string RR multiType to a const.
var rrTypeToString = map[RRType]string{
	A:     "a",
	MX:    "mx",
	NS:    "ns",
	CNAME: "cname",
}

// A simple map which transforms a string RR multiType to a const.
func RRTypeToString(rr RRType) string {
	if t, ok := rrTypeToString[rr]; ok {
		return t
	}
	return "no_rr"
}

// DnsRR represents the result of a DNS scan of a domain.
type DnsRR struct {
	tableName       struct{}     `pg:"dns_rrs,partition_by:RANGE(date)"` // Partitioning
	TaskID          int          // Number of the task set
	Task            *Task        `pg:"rel:has-one"`                 // Task structure
	SourceID        DataSourceID `pg:",type:bigint"`                // A listed source for the data.
	Source          *Source      `pg:"rel:has-one"`                 // Source pointer
	Date            time.Time    `pg:",use_zero,default:now()"`     // Date of the scan
	Domain          *Domain      `pg:"rel:has-one"`                 // DomainDomainCategory (Actually, FQDN) scanned
	DomainSubdomain string       `pg:",use_zero,type:varchar(255)"` //Subdomain(s) of the domain scanned
	DomainName      string       `pg:",use_zero,type:varchar(255)"` //name of the domain scanned
	DomainTLD       string       `pg:",use_zero,type:varchar(255)"` // TLD of the domain scanned
	ScanType        RRType       `pg:",use_zero,default:0"`         // Scan type of the result. See the const for more details.
	DerivedType     RRType       `pg:",use_zero,default:0"`         // When ahother RR requires to make a specific scan (like the IPs pointed by the domain as value of MX scan), its type appears here.
	Index           int          `pg:",use_zero,default:0"`
	IPValue         net.IP       // Address value in A RRs
	ValueSubdomain  string       `pg:",use_zero,type:varchar(255)"` // Subdomain value in MX, NS and CNAME RRs
	ValueName       string       `pg:",use_zero,type:varchar(255)"` // name value in MX, NS and CNAME RRs
	ValueTLD        string       `pg:",use_zero,type:varchar(255)"` // TLD value in MX, NS and CNAME RRs
	Priority        int          `pg:",use_zero,default:0"`         // Priority value in MX RRs.
	Accessible      bool         `pg:",use_zero"`                   // Accessibility value in A RRs.
	Valid           bool         `pg:",use_zero"`                   // Determines if a value is well written or in a valid range of values.
}
