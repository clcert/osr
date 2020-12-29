package models

// ServiceModel contains the metainformation related to the respective model.
var ServiceModel = Model{
	Name:        "Service",
	Description: "Information about software",
	StructType:  &Service{},
}

// Service Info
type Service struct {
	ID       string `pg:",pk,use_zero,type:varchar(255)"`
	Vendor   string `pg:",pk,use_zero,type:varchar(255)"`
	Name     string `pg:",pk,use_zero,type:varchar(255)"`
	Comments string
}
