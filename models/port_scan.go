package models

import (
	"net"
	"time"
)

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
		"CREATE INDEX IF NOT EXISTS port_scan_service_active ON ?TableName USING btree (service_active)",
		"CREATE INDEX IF NOT EXISTS port_scan_service_name ON ?TableName USING btree (service_name)",
		"CREATE INDEX IF NOT EXISTS port_scan_service_version ON ?TableName USING btree (service_version)",
		"CREATE INDEX IF NOT EXISTS port_scan_source_id ON ?TableName USING btree (source_id)",
		"SELECT partman.create_parent('public.port_scans', 'date', 'native', 'weekly');",
	},
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
	tableName      struct{}     `pg:"port_scans,partition_by:RANGE(date)"` // Partitioning
	TaskID         int          `pg:",type:bigint"`                        // Protocol of the importer session
	Task           *Task        `pg:"rel:has-one"`                         // Task structure
	SourceID       DataSourceID `pg:",pk,use_zero,type:bigint"`            // A listed source for the data.
	Source         *Source      `pg:"rel:has-one"`                         // Source pointer
	Date           time.Time    `pg:",pk,use_zero"`                        // Date of the scan
	ScanIP         net.IP       `pg:",pk"`                                 // IP address used to scan the server
	IP             net.IP       `pg:",pk"`                                 // Address
	PortNumber     uint16       `pg:",pk,type:bigint"`                     // Protocol number scanned
	Protocol       PortProtocol `pg:",pk,type:smallint,use_zero"`
	Port           *Port        `pg:"rel:has-one"`
	ServiceActive  bool         `pg:",use_zero,default:false"`
	ServiceName    string
	ServiceVersion string
	ServiceExtra   string
}
