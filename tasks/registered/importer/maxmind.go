package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/importer/maxmind/geolite2"
)

//TODO: join all of this in an only maxmind command (Import this data separatedly is purposeless)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "Countries and Subnet-to-ASN and Subnet-to-Country information",
			Command:     "import/maxmind-geolite2",
			Description: "Updates information about Countries and their Subnets, and ASNS and their ips.",
			URL:         "https://dev.maxmind.com/geoip/geoip2/maxmind/",
			Execute:     geolite2.Execute,
			Source:      models.MaxMind,
			NumSources:  2,
			NumSavers:   1,
		},
	)
}
