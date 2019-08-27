package models

import (
	"net"
)

func init() {
	DefaultModels.Append(SubnetCountryModel)
}

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
	TaskID           int `sql:",pk,type:bigint"`          // Number of the importer session
	Task             *Task                                // Task struct
	SourceID         DataSourceID `sql:",pk,type:bigint"` // A listed source for the data.
	Source           *Source                              // Source pointer
	Subnet           *net.IPNet `sql:",pk"`               // Subnet associated to this entry
	CountryGeonameId int        `sql:",pk,type:integer"`  // Geoname Number associated to this Subnet
	Country          *Country                             // Country struct
}