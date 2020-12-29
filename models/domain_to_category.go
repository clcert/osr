package models

import "github.com/go-pg/pg/v10/orm"

func init() {
	orm.RegisterTable((*DomainToCategory)(nil))
}

// DomainToCategoryModel contains the metainformation related to the respective model.
var DomainToCategoryModel = Model{
	Name:        "Domain To Category",
	Description: "A join between a domain and a category",
	StructType:  &DomainToCategory{},
}

// This structure defines a category for a internet domain.
type DomainToCategory struct {
	TaskID             int             // Number of the importer session
	Task               *Task           `pg:"rel:has-one"` // Task structure
	SourceID           DataSourceID    // A listed source for the data.
	Source             *Source         `pg:"rel:has-one"`                   // Source pointer
	DomainSubdomain    string          `pg:",pk,notnull,type:varchar(255)"` // Subdomain of the domain being categorized
	DomainName         string          `pg:",pk,notnull,type:varchar(255)"` // Domain name of the domain being categorized
	DomainTLD          string          `pg:",pk,notnull,type:varchar(255)"` // Domain tld of the domain being categorized
	Domain             Domain          `pg:"rel:has-one"`                   // Domain (Actually, FQDN) scanned
	DomainCategorySlug string          `pg:",pk,type:varchar(255)"`         // Numerical name of the category
	DomainCategory     *DomainCategory `pg:"rel:has-one"`
}
