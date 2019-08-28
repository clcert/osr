package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	chilean_dns "github.com/clcert/osr/tasks/processes/importer/clcert/chilean-dns"
	"github.com/clcert/osr/tasks/processes/importer/clcert/darknet"
	domain_categories "github.com/clcert/osr/tasks/processes/importer/clcert/domain-categories"
	port_scan "github.com/clcert/osr/tasks/processes/importer/clcert/port-scan"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "CLCERT Chilean DNS",
			Command:     "import/clcert-chilean-dns",
			Description: "Imports Mercury DNS Resource Records Scans",
			URL:         "",
			Source:      models.CLCERT,
			Execute:     chilean_dns.Execute,
			NumSources:  2,
			NumSavers:   1,
		},
		&tasks.Process{
			Name:        "CLCERT Darknet",
			Command:     "import/clcert-darknet",
			Description: "Imports information from a PCAP file captured by CLCERT's Darknet",
			URL:         "",
			Source:      models.CLCERT,
			Execute:     darknet.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
		&tasks.Process{
			Name:        "CLCERT Domain Categories Definition",
			Command:     "import/clcert-domain-categories",
			Description: "Imports Domain Categories classification by CLCERT.",
			URL:         "",
			Source:      models.CLCERT,
			Execute:     domain_categories.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
		&tasks.Process{
			Name:        "CLCERT Port Scan",
			Command:     "import/clcert-port-scan",
			Description: "Imports port scans locally made.",
			URL:         "",
			Source:      models.CLCERT,
			Execute:     port_scan.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
}
