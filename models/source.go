package models

import "github.com/go-pg/pg"

type DataSourceID uint64

const (
	UnknownSource DataSourceID  = iota
	CLCERT
	NICChile
	CIDRReport
	RIPE
	MaxMind
	Alexa
	Source1
	Source2
)

var SourceModel = Model{
	Name:                "Source",
	Description:         "Source Model represents the different source of the data contained by the OSR",
	StructType:          &Source{},
	AfterCreateFunction: createSources,
}

type Source struct {
	ID          DataSourceID `sql:",pk"`                        // Source Number
	Name        string       `sql:",notnull,type:varchar(255)"` // Source name
	URL         string       `sql:",type:varchar(512)"`         // Source URL
	Description string                                          // Source Description
}

// TODO: Extract this information from another source
func createSources(db *pg.DB) error {
	sources := []Source{
		{
			ID:          CLCERT,
			Name:        "CLCERT",
			URL:         "https://www.clcert.cl",
			Description: "CLCERT Chilean CERT. Providers of Darknet and DNS Scan data",
		},
		{
			ID:          NICChile,
			Name:        "NIC Chile",
			URL:         "https://www.nic.cl",
			Description: "Chilean Domain Administrator. Providers of the chilean domains list.",
		},
		{
			ID:          CIDRReport,
			Name:        "CIDR Report",
			URL:         "https://www.cidr-report.org",
			Description: "Daily reports of CIDR information. Providers of ASN names and declared countries.",
		},
		{
			ID:          RIPE,
			Name:        "RIPE",
			URL:         "https://ripe.net",
			Description: "Regional Internet Registry for Europe, the Middle East and parts of Central Asia. Providers of a source of chilean-assigned IPs",
		},
		{
			ID:          MaxMind,
			Name:        "MaxMind",
			URL:         "https://maxmind.com",
			Description: "Mantainers of Geolite2 Database. Providers of a source of chilean-assigned IPs and ASNs.",
		},
		{
			ID:          Alexa,
			Name:        "Alexa",
			URL:         "https://alexa.com",
			Description: "Providers of Alexa Top 1M web ranking",
		},
		{
			ID:          Source1,
			Name:        "Reserved Source #1",
			Description: "Closed Source, providers of some daily scans.",
		},
		{
			ID:          Source2,
			Name:        "Reserved Source #2",
			Description: "Closed Source, providers of some daily scans.",
		},
	}
	_, err := db.Model(&sources).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
