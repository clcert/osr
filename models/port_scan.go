package models

import (
	"net"
	"time"
)

func init() {
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
		"CREATE INDEX IF NOT EXISTS port_scan_service_active ON ?TableName USING btree (service_active)",
		"CREATE INDEX IF NOT EXISTS port_scan_service_name ON ?TableName USING btree (service_name)",
		"CREATE INDEX IF NOT EXISTS port_scan_service_version ON ?TableName USING btree (service_version)",
	},
}

// PortScan represents an open protocol port on a machine with an specific IP in a specific time.
type PortScan struct {
	TaskID         int          `sql:",type:bigint"` // Protocol of the importer session
	Task           *Task        // Task structure
	SourceID       DataSourceID `sql:",pk,notnull,type:bigint"` // A listed source for the data.
	Source         *Source      // Source pointer
	Date           time.Time    `sql:",pk,notnull"`     // Date of the scan
	ScanIP         net.IP       `sql:",pk"`             // IP address used to scan the server
	IP             net.IP       `sql:",pk"`             // Address
	PortNumber     uint16       `sql:",pk,type:bigint"` // Protocol number scanned
	Protocol       PortProtocol `sql:",pk,type:smallint,notnull"`
	Port           *Port
	ServiceActive  bool `sql:",notnull,default:false"`
	ServiceName    string
	ServiceVersion string
	ServiceExtra   string
}
