package models

// ASNModel contains the metainformation related to the respective model.
var ASNModel = Model{
	Name:        "ASN",
	Description: "Autonomous System Number",
	StructType:  &ASN{},
}

// ASN Represents an Autonomous System number.
type ASN struct {
	ID            int          `pg:",use_zero"`                   // Official assigned Number for the ASN
	Name          string       `pg:",use_zero,type:varchar(255)"` // Official name of the ASN
	CountryAlpha2 string       `pg:"type:varchar(2)"`             // 2 letter code of the ASN Country
	Subnets       []*SubnetASN `pg:"rel:has-many"`                // List of asns associated to the ASN
}
