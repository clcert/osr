package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/importer/gobcl/gob_domains"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:            "Chilean Government Subdomains",
			Command:         "import/gobcl-subdomains",
			Description:     "Imports subdomains based on gob.cl information.",
			URL:             "https://www.gob.cl/instituciones/",
			Execute:         gob_domains.Execute,
			DefaultSourceID: models.MaxMind,
			NumSources:      1,
			NumSavers:       1,
		},
	)
}
