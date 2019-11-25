package chilean_dns

import (
	"github.com/clcert/osr/tasks"
	"strconv"
)

// This function reads a remote server and savers all the
// results contained in folders with a signal empty file.
func Execute(args *tasks.Context) error {
	skipImport, err := strconv.ParseBool(args.Params["onlyIPASN"])
	if err != nil {
		skipImport = false
	}

	if !skipImport {
		err := parseScan(args)
		if err != nil {
			return err
		}
	}
	return GetIpAsnCountries(args)
}
