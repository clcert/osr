package models

func init() {
	DefaultModels.Append(DomainRankingModel)
}

var DomainRankingModel = Model{
	Name:        "Domain Ranking",
	Description: "Ranking based on Alexa Top 1M",
	StructType:  &DomainRanking{},
}

// Ranking of a domain, in a specified date, based on Alexa results.
type DomainRanking struct {
	TaskID          int `sql:",pk"`                              // Number of the importer session
	Task            *Task                                        // Task structure
	SourceID        DataSourceID `sql:",pk,type:bigint"`         // A listed source for the data.
	Source          *Source                                      // Source pointer
	Domain          *Domain                                      // DomainDomainCategory (Actually, FQDN) scanned
	DomainSubdomain string `sql:",pk,notnull,type:varchar(255)"` //Subdomain(s) of the domain scanned
	DomainName      string `sql:",pk,notnull,type:varchar(255)"` //name of the domain scanned
	DomainTLD       string `sql:",pk,notnull,type:varchar(255)"` // TLD of the domain scanned
	Ranking         int64  `sql:",notnull,default:-1"`           // Ranking value of the result
}
