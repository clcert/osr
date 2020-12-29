package models

import "time"

var ContactModel = Model{
	Name:        "Contact",
	Description: "Represents the contact information responsible of some internet resource",
	StructType:  &Contact{},
}

type Contact struct {
	ID               string    `pg:",pk,type:varchar(32)"`        // Contact ID
	Name             string    `pg:",use_zero,type:varchar(255)"` // Contact name
	Position         string    `pg:",use_zero,type:varchar(255)"` // Contact position
	LandlinePhone    string    `pg:",type:varchar(11)"`           // Contact Landline Phone
	MobilePhone      string    `pg:",type:varchar(11)"`           // Contact Mobile Phone
	Email            string    `pg:",type:varchar(512)"`          // Contact Email
	OrganizationName string    `pg:",type:varchar(512)"`          // Contact Org. Name
	OrganizationURL  string    `pg:",type:varchar(512)"`          // Contact Org. URL
	Description      string    // Contact Description
	UpdatedAt        time.Time `pg:"default:now()"`
}
