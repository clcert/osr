package models

import (
	"net"
)

// IpAsnCountryModel contains the metainformation related to the respective model.
var IpAsnCountryModel = Model{
	Name:        "Address ASN Country",
	Description: "Scanned Address with info about its ASN and Country associated in a given moment",
	StructType:  &IpAsnCountry{},
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS ip_asn_country_index ON ?TableName USING gist (ip inet_ops)",
	},
}

// IpAsnCountry represents a valid scanned Address (That is, properly formed and not in a reserved range).
type IpAsnCountry struct {
	TaskID           int `sql:",pk,type:bigint"`          // Number of the importer session related to this value
	Task             *Task                                // The importer session
	SourceID         DataSourceID `sql:",pk,type:bigint"` // A listed source for the data.
	Source           *Source                              // Source pointer
	IP               *net.IP `sql:",pk"`                  // Address of the relation
	AsnID            int                                  // ASN associated to the Address
	ASN              ASN
	CountryGeonameId int     // Country associated to the Address
	Country          Country // Country Struct
}
