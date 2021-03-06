package models

import (
	"time"
)

// DomainModel contains the metainformation related to the respective model.
var DomainModel = Model{
	Name:        "Domain",
	Description: "Domains model",
	StructType:  &Domain{},
}

// This structure defines an internet domain.
type Domain struct {
	TaskID           int              // Id of the task set
	Task             *Task            `pg:"rel:has-one"` // Task structure
	SourceID         DataSourceID     // A listed source for the data.
	Source           *Source          `pg:"rel:has-one"`                    // Source pointer
	Subdomain        string           `pg:",use_zero,pk,type:varchar(255)"` // Subdomain(s) of the entry
	Name             string           `pg:",pk,use_zero,type:varchar(255)"` // name of the entry
	TLD              string           `pg:",pk,use_zero,type:varchar(255)"` // TLD of the entry
	RegistrationDate time.Time        // Date of registration of the domain
	DeletionDate     time.Time        // Date of deletion of the domain.
	Categories       []DomainCategory `pg:"many2many:domain_to_categories"` // Categories associated to a domain
}

// Returns the canonical form of a FQDN.
func (d Domain) String() string {
	if len(d.Subdomain) > 0 {
		return d.Subdomain + "." + d.Name + "." + d.TLD
	}
	return d.Name + "." + d.TLD
}
