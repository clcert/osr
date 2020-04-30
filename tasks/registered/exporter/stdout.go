package exporter

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/exporter/stdout"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:            "Standard Output Data Exporter",
			Command:         "export/stdout",
			Description:     "Exports information received by sources to STDOUT",
			DefaultSourceID: models.CLCERT, // CLCERT transforms it
			Execute:         stdout.Execute,
			NumSources:      1,
			NumSavers:       0,
		},
	)
}

