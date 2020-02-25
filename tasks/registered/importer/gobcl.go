package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
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
			NumSources:      2,
			NumSavers:       1,
		},
	)
}
