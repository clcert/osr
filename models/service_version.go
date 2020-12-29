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
	ServiceID string    `pg:",pk,use_zero,type:varchar(255)"`
	Version   string    `pg:",pk,use_zero,type:varchar(255)"`
	Date      time.Time `pg:",pk,use_zero"`
	Comments  string
}
