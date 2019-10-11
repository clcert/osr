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
	TaskID   int `sql:",pk,type:bigint"`                  // Task session Number
	Task     *Task                                        // Task session
	SourceID DataSourceID `sql:",pk,notnull,type:bigint"` // A listed source for the data.
	Source   *Source                                      // Source pointer
	Subnet   *net.IPNet `sql:",pk"`                       // Subnet associated to this entry
	AsnID    int        `sql:",pk,type:bigint"`           // Number of the ASN associated to the subnet
	ASN      *ASN                                         // ASN struct
}
