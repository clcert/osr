package exporter

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/exporter/simple"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:            "Simple CSV Data Exporter",
			Command:         "export/simple",
			Description:     "Exports information from csv-like sources to a SFTP location",
			DefaultSourceID: models.CLCERT, // CLCERT transforms it
			Execute:         simple.Execute,
			NumSources:      -1,
			NumSavers:       1,
		},
	)
}
