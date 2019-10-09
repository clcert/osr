package models

import "time"

func init() {
	DefaultModels.Append(ContactModel)
}

var ContactModel = Model{
	Name:                "Contact",
	Description:         "Represents the contact information responsible of some internet resource",
	StructType:          &Contact{},
}

type Contact struct {
	Id               string	   `sql:",pk,type:varchar(32)"` // Contact ID
	Name             string    `sql:",notnull,type:varchar(255)"` // Contact name
	Position         string    `sql:",notnull,type:varchar(255)"` // Contact position
	LandlinePhone    string    `sql:",type:varchar(11)"`          // Contact Landline Phone
	MobilePhone      string    `sql:",type:varchar(11)"`          // Contact Mobile Phone
	Email            string    `sql:",type:varchar(512)"`         // Contact Email
	OrganizationName string    `sql:",type:varchar(512)"`         // Contact Org. Name
	OrganizationURL  string    `sql:",type:varchar(512)"`         // Contact Org. URL
	Description      string    // Source Description              // Contact Description
	UpdatedAt        time.Time `sql:"default:now()"`
}
