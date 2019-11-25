package importer

import (
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/tasks/processes/importer/alexa/rankings"
)

func init() {
	tasks.Registered.Register(
		&tasks.Process{
			Name:            "Alexa Chilean Domain Rankings",
			Command:         "import/alexa-rankings",
			Description:     "Updates information about chilean domains rankings.",
			URL:             "https://www.alexa.com/",
			DefaultSourceID: models.Alexa, // Alexa in sources
			Execute:         rankings.Execute,
			NumSources:      1,
			NumSavers:       1,
		},
	)
}
