package models

import (
	"net"
)

var SubnetCountryModel = Model{
	Name:        "Subnet Country",
	Description: "Subnet Country Model",
	StructType:  &SubnetCountry{},
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS subnet_country_index ON ?TableName USING gist (subnet inet_ops)",
	},
}

// SubnetASN represents the association between a subnet and a Country.
type SubnetCountry struct {
	TaskID           int          `pg:",pk,type:bigint"`  // Number of the importer session
	Task             *Task        `pg:"rel:has-one"`      // Task struct
	SourceID         DataSourceID `pg:",pk,type:bigint"`  // A listed source for the data.
	Source           *Source      `pg:"rel:has-one"`      // Source pointer
	Subnet           *net.IPNet   `pg:",pk"`              // Subnet associated to this entry
	CountryGeonameId int          `pg:",pk,type:integer"` // Geoname Number associated to this Subnet
	Country          *Country     `pg:"rel:has-one"`      // Country struct
}
