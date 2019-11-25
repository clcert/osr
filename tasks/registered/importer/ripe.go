package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	chilean_subnets "github.com/clcert/osr/tasks/processes/importer/ripe/chilean-subnets"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:            "RIPE List of chilean Subnets",
			Command:         "import/ripe-chilean-ips",
			Description:     "Updates information about ASNS and its Subnets from RIPE",
			URL:             "https://stat.ripe.net/data/country-resource-list/data.json?resource=cl",
			DefaultSourceID: models.RIPE,
			Execute:         chilean_subnets.Execute,
			NumSources:      1,
			NumSavers:       1,
		},
	)
}
