package models

import "github.com/go-pg/pg"

// DomainCategoryModel contains the metainformation related to the respective model.
var DomainCategoryModel = Model{
	Name:                "Domain Category",
	Description:         "Category for a Domain",
	StructType:          &DomainCategory{},
	AfterCreateFunction: createDomainCategories,
}

// This structure defines a category for a internet domain.
type DomainCategory struct {
	Slug        string `sql:",pk,type:varchar(255)"` // Numerical name of the category
	Name        string `sql:",type:varchar(255)"` // name of the category
	Description string // Short of the category
}

// createDomainCategories initializes the domain categories available on the system
func createDomainCategories(db *pg.DB) error {
	categories := []DomainCategory{
		{Slug: "gov", Name: "Government", Description: "Government Domains"},
		{Slug: "social", Name: "Social", Description: "Social Networks, forums, chatrooms, buy-sell portals, etc"},
		{Slug: "edu", Name: "Education", Description: "Websites of schools, universities or other educational institutions"},
		{Slug: "estore", Name: "Electronic Store", Description: "Electronic Store site. You can buy and products there"},
		{Slug: "service", Name: "General Services", Description: "Services like water, electricity, internet, medical centers, webhosting, etc"},
		{Slug: "finance", Name: "Financial sites", Description: "Sites of banks and other financial institutions (like credit card providers)"},
		{Slug: "media", Name: "Media", Description: "Media and news portals from TV, Radio or online."},
		{Slug: "other", Name: "Other", Description: "Other categories not covered before"},
	}
	_, err := db.Model(&categories).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
