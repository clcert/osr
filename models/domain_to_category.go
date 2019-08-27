package models

import "github.com/go-pg/pg/orm"

func init() {
	DefaultModels.Append(DomainToCategoryModel)
}

var DomainToCategoryModel = Model{
	Name:        "Domain To Category",
	Description: "A join between a domain and a category",
	StructType:  &DomainToCategory{},
}

func init() {
	orm.RegisterTable((*DomainToCategory)(nil))
}

// This structure defines a category for a internet domain.
type DomainToCategory struct {
	TaskID             int                                          // Number of the importer session
	Task               *Task                                        // Task structure
	SourceID           DataSourceID                                 // A listed source for the data.
	Source             *Source                                      // Source pointer
	DomainSubdomain    string `sql:",pk,notnull,type:varchar(255)"` // Subdomain of the domain being categorized
	DomainName         string `sql:",pk,notnull,type:varchar(255)"` // Domain name of the domain being categorized
	DomainTLD          string `sql:",pk,notnull,type:varchar(255)"` // Domain tld of the domain being categorized
	Domain             Domain                                       // Domain (Actually, FQDN) scanned
	DomainCategorySlug string `sql:",pk,type:varchar(255)"`         // Numerical Number of the category
	DomainCategory     *DomainCategory
}
