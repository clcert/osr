	package models

func init() {
	DefaultModels.Append(ASNModel)
}

var ASNModel = Model{
	Name:        "ASN",
	Description: "Autonomous System Number",
	StructType:  &ASN{},
}

// ASN Represents an Autonomous System number.
type ASN struct {
	ID            int    `sql:",notnull"`                   // Official assigned Number for the ASN
	Name          string `sql:",notnull,type:varchar(255)"` // Official name of the ASN
	CountryAlpha2 string `sql:"type:varchar(2)"`            // 2 letter code of the ASN Country
	Subnets       *[]SubnetASN                              // List of asns associated to the ASN
}
