package models

import (
	"github.com/go-pg/pg"
)

func init() {
	DefaultModels.Append(PortModel)
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
// Protocol groups all the scanned ports and their meanings.
type Port struct {
	Number      uint16 `sql:",pk,type:bigint"`          // Protocol number
	Name        string `sql:",notnull,type:varchar(255)"` // Protocol service name
	Description string // Protocol service description
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
		{Number: 102, Name: "S7", Description: "Siemens S7"},
		{Number: 110, Name: "POP3", Description: "Post Office Protocol v3"},
		{Number: 123, Name: "NTP", Description: "Network Time Protocol"},
		{Number: 143, Name: "IMAP", Description: "Internet Message Access Protocol"},
		{Number: 443, Name: "HTTPS", Description: "Hyper Text Transfer Protocol Secure"},
		{Number: 445, Name: "SMB", Description: "SAMBA"},
		{Number: 465, Name: "SMTP", Description: "Simple Mail Transfer Protocol (encrypted)"},
		{Number: 502, Name: "Modbus", Description: "Modicon Industrial Protocol"},
		{Number: 587, Name: "SMTP", Description: "Simple Mail Transfer Protocol (encrypted)"},
		{Number: 623, Name: "IPMI", Description: "Intelligent Platform Managment Interface"},
		{Number: 631, Name: "IPP", Description: "Common UNIX Printing System"},
		{Number: 993, Name: "IMAP", Description: "Internet Message Access Protocol (encrypted)"},
		{Number: 995, Name: "POP3", Description: "Post Office Protocol v3 (encrypted)"},
		{Number: 1433, Name: "MSSQL", Description: "Microsoft SQL Server"},
		{Number: 1521, Name: "Oracle", Description: "Oracle Server"},
		{Number: 1883, Name: "MQTT", Description: "Message Queueing Telemetry Transport"},
		{Number: 1900, Name: "UPnP", Description: "Universal Plug and Play"},
		{Number: 1911, Name: "Fox", Description: "Fox Protocol"},
		{Number: 3306, Name: "MySQL", Description: "MySQL Server & MariaDB Server"},
		{Number: 5432, Name: "Postgres", Description: "PostgreSQL Server"},
		{Number: 5632, Name: "pcAnywhere", Description: "Symantec pcAnywhere software"},
		{Number: 5672, Name: "AMQP", Description: "Advanced Message Queueing Protocol"},
		{Number: 5800, Name: "VNC", Description: "Virtual Network Computing (Java)"},
		{Number: 5900, Name: "VNC", Description: "Virtual Network Computing"},
		{Number: 5901, Name: "VNC", Description: "Virtual Network Computing"},
		{Number: 5902, Name: "VNC", Description: "Virtual Network Computing"},
		{Number: 5903, Name: "VNC", Description: "Virtual Network Computing"},
		{Number: 6443, Name: "Kubernetes", Description: "Open source Container Orchestration System"},
		{Number: 7547, Name: "CWMP", Description: "CPE WAN Management Protocol"},
		{Number: 8080, Name: "HTTP", Description: "Hyper Text Transfer Protocol (deployment)"},
		{Number: 8443, Name: "HTTPS", Description: "Hyper Text Transfer Protocol Secure (deployment)"},
		{Number: 8883, Name: "MQTT", Description: "Message Queueing Telemetry Transport"},
		{Number: 9090, Name: "Prometheus", Description: "Prometheus Monitoring System"},
		{Number: 9200, Name: "ElasticSearch", Description: "ElasticSearch service"},
		{Number: 27017, Name: "MongoDB", Description: "MongoDB NOSQL Database"},
		{Number: 27018, Name: "MongoDB", Description: "MongoDB NOSQL Database"},
		{Number: 47808, Name: "BACnet", Description: "ASHRAE building automation and control networking protocol"},
	}
	_, err := db.Model(&ports).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
