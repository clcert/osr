package models

import "time"

// SoftwareVersionModel contains the metainformation related to the respective model.
var SoftwareVersionModel = Model{
	Name:        "Software Version",
	Description: "Information about a software version",
	StructType:  &SoftwareVersion{},
}

// Software Version Timelines
type SoftwareVersion struct {
	Vendor       string    `sql:",notnull,type:varchar(255)"`
	SoftwareName string    `sql:",notnull,type:varchar(255)"`
	SoftwareTag  string    `sql:",pk,notnull,type:varchar(255)"`
	Version      string    `sql:",pk,notnull,type:varchar(255)"`
	LaunchDate   time.Time `sql:",pk,notnull,type:varchar(255)"`
	Comments string
}
