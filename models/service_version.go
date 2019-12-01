package models

import "time"

// ServiceVersionModel contains the metainformation related to the respective model.
var ServiceVersionModel = Model{
	Name:        "Service Version",
	Description: "Information about a software version",
	StructType:  &ServiceVersion{},
}

// Service Version Timelines
type ServiceVersion struct {
	ServiceID string    `sql:",pk,notnull,type:varchar(255)"`
	Version   string    `sql:",pk,notnull,type:varchar(255)"`
	Date      time.Time `sql:",pk,notnull"`
	Comments  string
}
