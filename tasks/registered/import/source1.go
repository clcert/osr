package _import

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/import/source1/reports"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "Source#1 Reports",
			Command:     "import/source1-reports",
			Description: "Imports reports from source#1",
			URL:         "",
			Source:      models.Source1,
			Execute:     reports.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
}
