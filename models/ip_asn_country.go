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
		"CREATE INDEX ip_asn_countries_task_id ON ?TableName USING btree (task_id)",
	},
}

// IpAsnCountry represents a valid scanned Address (That is, properly formed and not in a reserved range).
type IpAsnCountry struct {
	TaskID           int          `pg:",pk,type:bigint"` // Number of the importer session related to this value
	Task             *Task        `pg:"rel:has-one"`     // The importer session
	SourceID         DataSourceID `pg:",pk,type:bigint"` // A listed source for the data.
	Source           *Source      `pg:"rel:has-one"`     // Source pointer
	IP               net.IP       `pg:",pk"`             // Address of the relation
	AsnID            int          // ASN associated to the Address
	ASN              ASN          `pg:"rel:has-one"`
	CountryGeonameId int          `pg:",type:integer"` // Country associated to the Address
	Country          *Country     `pg:"rel:has-one"`   // Country Struct
}
