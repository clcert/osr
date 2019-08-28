package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/importer/source2/reports"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "Source#2 Reports",
			Command:     "import/source2-reports",
			Description: "Imports reports from source#2",
			URL:         "",
			Source:      models.Source2,
			Execute:     reports.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
}
