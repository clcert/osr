package models

// DomainRankingModel contains the metainformation related to the respective model.
var DomainRankingModel = Model{
	Name:        "Domain Ranking",
	Description: "Ranking based on Alexa Top 1M",
	StructType:  &DomainRanking{},
}

// Ranking of a domain, in a specified date, based on Alexa results.
type DomainRanking struct {
	TaskID          int          `pg:",pk"`                           // Number of the importer session
	Task            *Task        `pg:"rel:has-one"`                   // Task structure
	SourceID        DataSourceID `pg:",pk,type:bigint"`               // A listed source for the data.
	Source          *Source      `pg:"rel:has-one"`                   // Source pointer
	Domain          *Domain      `pg:"rel:has-one"`                   // DomainDomainCategory (Actually, FQDN) scanned
	DomainSubdomain string       `pg:",pk,notnull,type:varchar(255)"` //Subdomain(s) of the domain scanned
	DomainName      string       `pg:",pk,notnull,type:varchar(255)"` //name of the domain scanned
	DomainTLD       string       `pg:",pk,notnull,type:varchar(255)"` // TLD of the domain scanned
	Ranking         int64        `pg:",notnull,default:-1"`           // Ranking value of the result
}
