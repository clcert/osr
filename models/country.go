package models

// CountryModel contains the metainformation related to the respective model.
var CountryModel = Model{
	Name:        "Country",
	Description: "Country Model",
	StructType:  &Country{},
}

// Country represents a country in the world.
type Country struct {
	Alpha2    string           `pg:",unique,type:varchar(2)"` // 2 letter representation of the country
	Name      string           `pg:",type:varchar(255)"`      // name of the country in spanish
	GeonameId int              `pg:",pk,type:integer"`        // Geoname Number of the country, as noted by Geolite database.
	Subnets   []*SubnetCountry `pg:"rel:has-many"`            // List of asns associated to the country.
}
