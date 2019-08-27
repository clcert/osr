package models

func init() {
	DefaultModels.Append(CountryModel)
}

var CountryModel = Model{
	Name:        "Country",
	Description: "Country Model",
	StructType:  &Country{},
}

// Country represents a country in the world, duh.
type Country struct {
	Alpha2    string `sql:",unique,type:varchar(2)"` // 2 letter representation of the country
	Name      string `sql:",type:varchar(255)"`      // name of the country in spanish
	GeonameId int    `sql:",pk,type:integer"`        // Geoname Number of the country, as noted by Geolite database.
	Subnets   *[]SubnetCountry                       // List of asns associated to the country.
}