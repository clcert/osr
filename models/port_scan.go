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
	Name:        "Protocol Scan",
	Description: "Protocol Scans",
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
	StructType:          &Protocol{},
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
	TaskID     int          `sql:",type:bigint"` // Protocol of the importer session
	Task       *Task        // Task structure
	SourceID   DataSourceID `sql:",pk,notnull,type:bigint"` // A listed source for the data.
	Source     *Source      // Source pointer
	Date       time.Time    `sql:",pk,notnull"`       // Date of the scan
	ScanIP     net.IP       `sql:",pk"`               // IP address used to scan the server
	IP         net.IP       `sql:",pk"`               // Address
	PortNumber uint16       `sql:",pk,type:smallint"` // Protocol number scanned
	Protocol   PortProtocol `sql:",pk,type:smallint,notnull"`
	Port       *Protocol
}

// Protocol groups all the scanned ports and their meanings.
type Protocol struct {
	Port        uint16 `sql:",pk,type:smallint"`          // Protocol number
	Name        string `sql:",notnull,type:varchar(255)"` // Protocol service name
	Description string // Protocol service description
}

// TODO: Extract this information from another source
// createPortDefinitions inserts the used port definitions when port table is created.
func createPortDefinitions(db *pg.DB) error {
	ports := []Protocol{
		{Port: 21, Name: "FTP", Description: "File Transfer Protocol"},
		{Port: 22, Name: "SSH", Description: "Secure Shell"},
		{Port: 23, Name: "Telnet", Description: "Telnet"},
		{Port: 25, Name: "SMTP", Description: "Simple Mail Transfer Protocol"},
		{Port: 53, Name: "DNS", Description: "Domain name System"},
		{Port: 80, Name: "HTTP", Description: "Hyper Text Transfer Protocol"},
		{Port: 102, Name: "S7", Description: "Siemens S7"},
		{Port: 110, Name: "POP3", Description: "Post Office Protocol v3"},
		{Port: 123, Name: "NTP", Description: "Network Time Protocol"},
		{Port: 143, Name: "IMAP", Description: "Internet Message Access Protocol"},
		{Port: 443, Name: "HTTPS", Description: "Hyper Text Transfer Protocol Secure"},
		{Port: 445, Name: "SMB", Description: "SAMBA"},
		{Port: 465, Name: "SMTP", Description: "Simple Mail Transfer Protocol (encrypted)"},
		{Port: 502, Name: "Modbus", Description: "Modicon Industrial Protocol"},
		{Port: 587, Name: "SMTP", Description: "Simple Mail Transfer Protocol (encrypted)"},
		{Port: 623, Name: "IPMI", Description: "Intelligent Platform Managment Interface"},
		{Port: 631, Name: "IPP", Description: "Common UNIX Printing System"},
		{Port: 993, Name: "IMAP", Description: "Internet Message Access Protocol (encrypted)"},
		{Port: 995, Name: "POP3", Description: "Post Office Protocol v3 (encrypted)"},
		{Port: 1433, Name: "MSSQL", Description: "Microsoft SQL Server"},
		{Port: 1521, Name: "Oracle", Description: "Oracle Server"},
		{Port: 1883, Name: "MQTT", Description: "Message Queueing Telemetry Transport"},
		{Port: 1900, Name: "UPnP", Description: "Universal Plug and Play"},
		{Port: 1911, Name: "Fox", Description: "Fox Protocol"},
		{Port: 3306, Name: "MySQL", Description: "MySQL Server & MariaDB Server"},
		{Port: 5432, Name: "Postgres", Description: "PostgreSQL Server"},
		{Port: 5632, Name: "pcAnywhere", Description: "Symantec pcAnywhere software"},
		{Port: 5672, Name: "AMQP", Description: "Advanced Message Queueing Protocol"},
		{Port: 5800, Name: "VNC", Description: "Virtual Network Computing (Java)"},
		{Port: 5900, Name: "VNC", Description: "Virtual Network Computing"},
		{Port: 5901, Name: "VNC", Description: "Virtual Network Computing"},
		{Port: 5902, Name: "VNC", Description: "Virtual Network Computing"},
		{Port: 5903, Name: "VNC", Description: "Virtual Network Computing"},
		{Port: 6443, Name: "Kubernetes", Description: "Open source Container Orchestration System"},
		{Port: 7547, Name: "CWMP", Description: "CPE WAN Management Protocol"},
		{Port: 8080, Name: "HTTP", Description: "Hyper Text Transfer Protocol (deployment)"},
		{Port: 8443, Name: "HTTPS", Description: "Hyper Text Transfer Protocol Secure (deployment)"},
		{Port: 8883, Name: "MQTT", Description: "Message Queueing Telemetry Transport"},
		{Port: 9090, Name: "Prometheus", Description: "Prometheus Monitoring System"},
		{Port: 9200, Name: "ElasticSearch", Description: "ElasticSearch service"},
		{Port: 27017, Name: "MongoDB", Description: "MongoDB NOSQL Database"},
		{Port: 27018, Name: "MongoDB", Description: "MongoDB NOSQL Database"},
		{Port: 47808, Name: "BACnet", Description: "ASHRAE building automation and control networking protocol"},
	}
	_, err := db.Model(&ports).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
