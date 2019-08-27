package _import

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/import/alexa/rankings"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:        "Alexa Chilean Domain Rankings",
			Command:     "import/alexa-rankings",
			Description: "Updates information about chilean domains rankings.",
			URL:         "https://www.alexa.com/",
			Source:      models.Alexa, // Alexa in sources
			Execute:     rankings.Execute,
			NumSources:  1,
			NumSavers:   1,
		},
	)
}
