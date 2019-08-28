package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/importer/cidr-report/asns"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "ASNS",
			Command:     "import/cidr-report-asns",
			Description: "Imports information about AS numbers and names",
			URL:         "http://www.cidr-report.org/as2.0/autnums.html",
			Source:      models.CIDRReport,
			Execute:     asns.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
}
