package models

import (
	"github.com/go-pg/pg"
	"net"
	"time"
)

func init() {
	DefaultModels.Append(ReportEntryModel)
	DefaultModels.Append(ReportTypeModel)
}


// ReportTypeID represents the possible types of reports present in the system
type ReportTypeID int

const (
	UnknownReport ReportTypeID = iota
	BotReport
	BruteforceReport
	C2Report
	DarknetReport
	HoneypotReport
	DNSResolverReport
	PhishingReport
	ProxyReport
	SpamReport
)

// ReportEntryModel contains the metainformation related to the respective model.
var ReportEntryModel = Model{
	Name:        "Reported IPs",
	Description: "List of IPs reported as curious",
	StructType:  &ReportEntry{},
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS report_entry_timestamp ON ?TableName USING btree (date)",
		"CREATE INDEX IF NOT EXISTS report_entry_ip ON ?TableName USING gist (ip inet_ops)",
		"CREATE INDEX IF NOT EXISTS report_entry_source_id ON ?TableName USING btree (source_id)",
		"CREATE INDEX IF NOT EXISTS report_entry_report_type_id ON ?TableName USING btree (report_type_id)",
	},
}

// ReportTypeModel contains the metainformation related to the respective model.
var ReportTypeModel = Model{
	Name:                "Report Category",
	Description:         "A specific category for a report IP entry",
	StructType:          &ReportType{},
	AfterCreateFunction: CreateReportTypes,
}

// ReportEntry represents an entry of an IP report.
type ReportEntry struct {
	TaskID       int `sql:",type:bigint"`         // Number of the importer session
	Task         *Task                            // Task structure
	SourceID     DataSourceID `sql:",pk,notnull"` // A listed source for the data.
	Source       *Source                          // Source pointer
	ReportTypeID ReportTypeID `sql:",pk,notnull"` // Type of report
	ReportType   *ReportType                      // Type of report
	Date         time.Time `sql:",pk,notnull"`    // Date of the scan
	IP           net.IP    `sql:",pk,notnull"`    // Source Address (scanned device)
	Properties   map[string]string                // Report extra metadata
}

type ReportType struct {
	ID          ReportTypeID `sql:",pk"`                        // scanned Category ID
	Name        string       `sql:",notnull,type:varchar(255)"` // scanned Category name
	Description string                                          // scanned Category description
}

// TODO: Extract this information from another source
func CreateReportTypes(db *pg.DB) error {
	ports := []ReportType{
		{ID: BotReport, Name: "Bot", Description: "IP Related to a Botnet"},
		{ID: BruteforceReport, Name: "Bruteforce", Description: "IPs found executing Bruteforce on SSH/FTP/etc services"},
		{ID: C2Report, Name: "C2", Description: "IPs found used as Comand and Control servers"},
		{ID: DarknetReport, Name: "Darknet", Description: "IPs found by Darknets"},
		{ID: HoneypotReport, Name: "Honeypot", Description: "IPs found by Honeypots"},
		{ID: DNSResolverReport, Name: "OpenResolver", Description: "IPs found being used as Open DNS Resolvers"},
		{ID: PhishingReport, Name: "Phishing", Description: "IPs found being used for Phishing websites"},
		{ID: ProxyReport, Name: "Proxy", Description: "IPs found being used as public Proxies"},
		{ID: SpamReport, Name: "SPAM", Description: "IPs found being used as SPAM mail servers"},
	}
	_, err := db.Model(&ports).OnConflict("DO NOTHING").Insert()
	if err != nil {
	}
	return nil
}
