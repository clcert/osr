package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	chilean_dns "github.com/clcert/osr/tasks/processes/importer/clcert/chilean-dns"
	fix_accessible "github.com/clcert/osr/tasks/processes/importer/clcert/chilean-dns/fix-accessible"
	"github.com/clcert/osr/tasks/processes/importer/clcert/darknet"
	domain_categories "github.com/clcert/osr/tasks/processes/importer/clcert/domain-categories"
	grabber_protocol_scan "github.com/clcert/osr/tasks/processes/importer/clcert/grabber-protocol-scan"
	port_scan "github.com/clcert/osr/tasks/processes/importer/clcert/port-scan"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:            "CLCERT Chilean DNS",
			Command:         "import/clcert-chilean-dns",
			Description:     "Imports Mercury DNS Resource Records Scans",
			URL:             "",
			DefaultSourceID: models.CLCERT,
			Execute:         chilean_dns.Execute,
			NumSources:      2,
			NumSavers:       1,
		},
		&tasks.Process{
			Name:            "Fix accessible in Chilean DNS scan",
			Command:         "import/clcert-chilean-dns/fix-accessible",
			Description:     "Fixes currently defined Accessible flag using a custom directory with ZMAP scans",
			URL:             "",
			DefaultSourceID: models.CLCERT,
			Execute:         fix_accessible.Execute,
			NumSources:      1,
			NumSavers:       0,
		},
		&tasks.Process{
			Name:            "CLCERT Darknet",
			Command:         "import/clcert-darknet",
			Description:     "Imports information from a PCAP file captured by CLCERT's Darknet",
			URL:             "",
			DefaultSourceID: models.CLCERT,
			Execute:         darknet.Execute,
			NumSources:      1,
			NumSavers:       1,
		},
		&tasks.Process{
			Name:            "CLCERT Domain Categories Definition",
			Command:         "import/clcert-domain-categories",
			Description:     "Imports Domain Categories classification by CLCERT.",
			URL:             "",
			DefaultSourceID: models.CLCERT,
			Execute:         domain_categories.Execute,
			NumSources:      1,
			NumSavers:       1,
		},
		&tasks.Process{
			Name:            "CLCERT Port Scan",
			Command:         "import/clcert-port-scan",
			Description:     "Imports port scans locally made.",
			URL:             "",
			DefaultSourceID: models.CLCERT,
			Execute:         port_scan.Execute,
			NumSources:      1,
			NumSavers:       1,
		},
		&tasks.Process{
			Name:            "CLCERT Grabber Protocol Scan",
			Command:         "import/clcert-protocol-scan",
			Description:     "Imports protocol scans locally made.",
			URL:             "",
			DefaultSourceID: models.CLCERT,
			Execute:         grabber_protocol_scan.Execute,
			NumSources:      1,
			NumSavers:       1,
		},
	)
}
