package models

import (
	"github.com/go-pg/pg"
	"net"
	"time"
)

func init() {
	DefaultModels.Append(PortModel)
	DefaultModels.Append(PortScanModel)
}

// PortScanModel contains the metainformation related to the respective model.
var PortScanModel = Model{
	Name:        "Port Scan",
	Description: "Port Scans",
	StructType:  &PortScan{},
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS port_scan_ip ON ?TableName USING gist (ip inet_ops)",
		"CREATE INDEX IF NOT EXISTS port_scan_scan_ip ON ?TableName USING gist (scan_ip inet_ops)",
		"CREATE INDEX IF NOT EXISTS port_scan_timestamp ON ?TableName USING btree (date)",
		"CREATE INDEX IF NOT EXISTS port_scan_port ON ?TableName USING btree (port_number)",
	},
}

// PortModel contains the metainformation related to the respective model.
var PortModel = Model{
	Name:                "Ports",
	Description:         "Ports definition",
	StructType:          &Port{},
	AfterCreateFunction: createPortDefinitions,
}
// PortProtocol represents the transport protocol checked in a port scan.
type PortProtocol int

const (
	UnknownProtocol PortProtocol = iota
	TCP
	UDP
)

// PortScan represents an open protocol port on a machine with an specific IP in a specific time.
type PortScan struct {
	TaskID     int `sql:",type:bigint"`                     // Number of the importer session
	Task       *Task                                        // Task structure
	SourceID   DataSourceID `sql:",pk,notnull,type:bigint"` // A listed source for the data.
	Source     *Source                                      // Source pointer
	Date       time.Time    `sql:",pk,notnull"`             // Date of the scan
	ScanIP     net.IP       `sql:",pk"`                     // IP address used to scan the server
	IP         net.IP       `sql:",pk"`                     // Address
	PortNumber uint16       `sql:",pk,type:smallint"`       // Port number scanned
	Protocol   PortProtocol `sql:",pk,type:smallint,notnull"`
	Port       *Port
}

// Port groups all the scanned ports and their meanings.
type Port struct {
	Number      uint16 `sql:",pk,type:smallint"`          // Port number
	Name        string `sql:",notnull,type:varchar(255)"` // Port service name
	Description string                                    // Port service description
}

// TODO: Extract this information from another source
// createPortDefinitions inserts the used port definitions when port table is created.
func createPortDefinitions(db *pg.DB) error {
	ports := []Port{
		{Number: 21, Name: "FTP", Description: "File Transfer Protocol"},
		{Number: 22, Name: "SSH", Description: "Secure Shell"},
		{Number: 23, Name: "Telnet", Description: "Telnet"},
		{Number: 25, Name: "SMTP", Description: "Simple Mail Transfer Protocol"},
		{Number: 53, Name: "DNS", Description: "Domain name System"},
		{Number: 80, Name: "HTTP", Description: "Hyper Text Transfer Protocol"},
		{Number: 110, Name: "POP3", Description: "Post Office Protocol v3"},
		{Number: 123, Name: "NTP", Description: "Network Time Protocol"},
		{Number: 143, Name: "IMAP", Description: "Internet Message Access Protocol"},
		{Number: 443, Name: "HTTPS", Description: "Hyper Text Transfer Protocol Secure"},
		{Number: 465, Name: "SMTP", Description: "Simple Mail Transfer Protocol (encrypted)"},
		{Number: 587, Name: "SMTP", Description: "Simple Mail Transfer Protocol (encrypted)"},
		{Number: 631, Name: "CUPS", Description: "Common UNIX Printing System"},
		{Number: 993, Name: "IMAP", Description: "Internet Message Access Protocol (encrypted)"},
		{Number: 995, Name: "POP3", Description: "Post Office Protocol v3 (encrypted)"},
		{Number: 1433, Name: "MSSQL", Description: "Microsoft SQL Server"},
		{Number: 3306, Name: "MySQL", Description: "MySQL Server & MariaDB Server"},
		{Number: 5432, Name: "Postgres", Description: "PostgreSQL Server"},
		{Number: 5800, Name: "VNC", Description: "Virtual Network Computing (Java)"},
		{Number: 5900, Name: "VNC", Description: "Virtual Network Computing"},
		{Number: 8080, Name: "HTTP", Description: "Hyper Text Transfer Protocol (deployment)"},
		{Number: 8443, Name: "HTTPS", Description: "Hyper Text Transfer Protocol Secure (deployment)"},
		{Number: 9200, Name: "ElasticSearch", Description: "ElasticSearch service"},
		{Number: 27018, Name: "MongoDB", Description: "MongoDB NOSQL Database"},
	}
	_, err := db.Model(&ports).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
