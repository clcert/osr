package _import

import (
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/transform/ips/compare"
	delta_dates "github.com/clcert/osr/tasks/processes/transform/ips/compare-historic"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "Compare IPs",
			Command:     "transform/compare-ips",
			Description: "Compares two ip sets, providing the common IPs and the missing IPs on each set.",
			Source:      1, // Data aggregated by CLCERT
			Execute:     compare.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "Historical comparison of two sources",
			Command:     "transform/compare-historic",
			Description: "Compares two IP sources, providing the number of common and distinct IPs on each set by date.",
			Source:      1, // Data aggregated by CLCERT
			Execute:     delta_dates.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
}
