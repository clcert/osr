package chilean_subnets

import (
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
)

// Downloads the webpage and uploads its data to the
// Subnets database.
func Execute(args *tasks.Context) error {
	source := args.Sources[0]
	saver := args.Savers[0]

	country := &models.Country{}

	db, err := databases.GetPostgresReader()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Model(country).
		Where("alpha2 = ?", "CL").
		Select()
	if err != nil {
		args.Log.WithFields(logrus.Fields{
			"country_alpha2": "cl",
		}).Error("sql error: %s", err)
		return err
	}

	for {
		entry := source.Next()
		if entry == nil {
			break
		}
		if err := saveRIPE(entry, saver, args, country); err != nil {
			return err
		}
	}
	return nil
}
