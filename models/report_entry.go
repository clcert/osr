package models

import (
	"net"
	"time"
)

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

// ReportEntry represents an entry of an IP report.
type ReportEntry struct {
	TaskID       int               `pg:",type:bigint"` // Number of the importer session
	Task         *Task             `pg:"rel:has-one"`  // Task structure
	SourceID     DataSourceID      `pg:",pk,notnull"`  // A listed source for the data.
	Source       *Source           `pg:"rel:has-one"`  // Source pointer
	ReportTypeID ReportTypeID      `pg:",pk,notnull"`  // Type of report
	ReportType   *ReportType       `pg:"rel:has-one"`  // Type of report
	Date         time.Time         `pg:",pk,notnull"`  // Date of the scan
	IP           net.IP            `pg:",pk,notnull"`  // Source Address (scanned device)
	Properties   map[string]string // Report extra metadata
}
