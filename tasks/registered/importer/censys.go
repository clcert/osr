package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	port_scan "github.com/clcert/osr/tasks/processes/importer/censys/port-scan"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "Censys Port Scan",
			Command:     "import/censys-port-scan",
			Description: "Imports port scans made by Censys.",
			URL:         "",
			Source:      models.Censys,
			Execute:     port_scan.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
}
