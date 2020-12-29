package models

import (
	"net"
)

var SubnetASNModel = Model{
	Name:        "Subnet ASN",
	Description: "Subnet ASN Model",
	StructType:  &SubnetASN{},
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS subnet_asn_index ON ?TableName USING gist (subnet inet_ops)",
	},
}

// SubnetASN represents the association between a subnet and an ASN.
type SubnetASN struct {
	TaskID   int          `pg:",pk,type:bigint"`          // Task session Number
	Task     *Task        `pg:"rel:has-one"`              // Task session
	SourceID DataSourceID `pg:",pk,use_zero,type:bigint"` // A listed source for the data.
	Source   *Source      `pg:"rel:has-one"`              // Source pointer
	Subnet   *net.IPNet   `pg:",pk"`                      // Subnet associated to this entry
	AsnID    int          `pg:",pk,type:bigint"`          // Number of the ASN associated to the subnet
	ASN      *ASN         `pg:"rel:has-one"`              // ASN struct
}
