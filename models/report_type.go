package models

import "github.com/go-pg/pg/v10"

// ReportTypeModel contains the metainformation related to the respective model.
var ReportTypeModel = Model{
	Name:                "Report Category",
	Description:         "A specific category for a report IP entry",
	StructType:          &ReportType{},
	AfterCreateFunction: CreateReportTypes,
}

type ReportType struct {
	ID          ReportTypeID `pg:",pk"`                        // scanned Category ID
	Name        string       `pg:",notnull,type:varchar(255)"` // scanned Category name
	Description string       // scanned Category description
}

// TODO: Extract this information from another source
func CreateReportTypes(db *pg.DB) error {
	ports := []ReportType{
		{ID: BotReport, Name: "Bot", Description: "IP Related to a Botnet"},
		{ID: BruteforceReport, Name: "Bruteforce", Description: "IPs found executing Bruteforce on SSH/FTP/etc services"},
		{ID: C2Report, Name: "C2", Description: "IPs found used as Comand and Control servers"},
		{ID: DarknetReport, Name: "Darknet", Description: "IPs found by Darknets"},
		{ID: HoneypotReport, Name: "Honeypot", Description: "IPs found by Honeypots"},
		{ID: DNSResolverReport, Name: "OpenResolver", Description: "IPs found being used as Open DNS Resolvers"},
		{ID: PhishingReport, Name: "Phishing", Description: "IPs found being used for Phishing websites"},
		{ID: ProxyReport, Name: "Proxy", Description: "IPs found being used as public Proxies"},
		{ID: SpamReport, Name: "SPAM", Description: "IPs found being used as SPAM mail servers"},
	}
	_, err := db.Model(&ports).OnConflict("DO NOTHING").Insert()
	if err != nil {
	}
	return nil
}
