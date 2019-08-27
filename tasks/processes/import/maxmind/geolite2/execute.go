package geolite2

import (
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
)

type FuncMap map[string]func(entry sources.Entry, saver savers.Saver, args *tasks.Args) error

var nameToFunc = FuncMap{
	"GeoLite2-Country-Locations-es.csv": saveCountries,
	"GeoLite2-Country-Blocks-IPv4.csv":  saveCountrySubnets,
	"GeoLite2-ASN-Blocks-IPv4.csv":      saveASNSubnets,
}

func Execute(args *tasks.Args) error {
	saver := args.Savers[0]
	// We parse all the sources
	for _, source := range args.Sources {
		parseSource(source, saver, args)
	}
	args.Log.Info("Done parsing all sources, exiting...")
	return nil
}
