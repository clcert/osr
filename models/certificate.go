package models

import (
	"crypto/x509"
	"net"
	"time"
)

// PortScanModel contains the metainformation related to the respective model.
var CertificateModel = Model{
	Name:        "Protocol Certificates",
	Description: "Protocol Certificates",
	StructType:  &Certificate{},
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS certificate_ip ON ?TableName USING gist (ip inet_ops)",
		"CREATE INDEX IF NOT EXISTS certificate_scan_ip ON ?TableName USING gist (scan_ip inet_ops)",
		"CREATE INDEX IF NOT EXISTS certificate_timestamp ON ?TableName USING btree (date)",
		"CREATE INDEX IF NOT EXISTS certificate_port ON ?TableName USING btree (port_number)",
		"CREATE INDEX IF NOT EXISTS certificate_key_size ON ?TableName USING btree (key_size)",
		"CREATE INDEX IF NOT EXISTS certificate_status ON ?TableName USING btree (status)",
		"CREATE INDEX IF NOT EXISTS certificate_signature_algorithm ON ?TableName USING btree (signature_algorithm)",
		"CREATE INDEX IF NOT EXISTS certificate_authority ON ?TableName USING btree (authority)",
		"CREATE INDEX IF NOT EXISTS certificate_expiration_date ON ?TableName USING btree (expiration_date)",
	},
}

type CertStatus int

const (
	CertUnknownError CertStatus = iota
	CertValid
	CertExpired
	CertSelfSigned
	CertUnparseable
	CertEmptyChain
	CertUnknownAuthority
	CertNotAuthorizedToSign
)

type TLSProto int

const (
	UnknownTLSPRoto TLSProto = iota
	SSL30
	TLS10
	TLS11
	TLS12
	TLS13
)

// PortScan represents an open protocol port on a machine with an specific IP in a specific time.
type Certificate struct {
	TaskID             int                     `pg:",use_zero,type:bigint"`    // Protocol of the importer session
	Task               *Task                   `pg:"rel:has-one"`              // Task structure
	SourceID           DataSourceID            `pg:",pk,use_zero,type:bigint"` // A listed source for the data.
	Source             *Source                 `pg:"rel:has-one"`              // Source pointer
	Date               time.Time               `pg:",pk,use_zero"`             // Date of the scan
	ScanIP             net.IP                  `pg:",pk"`                      // IP address used to scan the server
	IP                 net.IP                  `pg:",pk"`                      // Address
	PortNumber         uint16                  `pg:",pk,type:bigint"`          // Protocol number scanned
	Port               *Port                   `pg:"rel:has-one"`              // Port object
	Status             CertStatus              `pg:",use_zero"`
	KeySize            int                     `pg:",use_zero"` // Key Size
	ExpirationDate     time.Time               `pg:",use_zero"` // Expiration Date
	OrganizationName   string                  `pg:",use_zero"` // Organization Name
	OrganizationURL    string                  `pg:",use_zero"` // Organization URL
	Authority          string                  `pg:",use_zero"` // Certificate Authority
	SignatureAlgorithm x509.SignatureAlgorithm `pg:",use_zero"` // Signature Algorithm
	TLSProtocol        TLSProto                `pg:",use_zero"` // TLS Protocol
}
