package _import

import (
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/transform/ips/compare"
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
			Name:        "Distinct IPs over time",
			Command:     "transform/compare-ips",
			Description: "Compares two subnet sets, providing the common IPs and the missing IPs on each set.",
			Source:      1, // Data aggregated by CLCERT
			Execute:     compare.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
}
