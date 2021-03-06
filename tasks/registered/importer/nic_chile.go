package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/importer/nic-chile/domains"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:            "NIC Chile: New and Deleted Domains",
			Command:         "import/nic-chile-domains",
			Description:     "Imports information about new domains",
			URL:             "https://nic.cl/registry/Ultimos.do?t=1w&f=csv",
			DefaultSourceID: models.NICChile,
			Execute:         domains.Execute,
			NumSources:      tasks.InfiniteSources,
			NumSavers:       1,
		},
	)
}
