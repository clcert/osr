package chilean_dns

import (
	"github.com/clcert/osr/tasks"
	"strconv"
)

func Execute(args *tasks.Context) error {
	skipImport := false
	onlyIPASN, ok := args.Params["onlyIPASN"]
	if ok {
		skipImport, _ = strconv.ParseBool(onlyIPASN)
	}

	if !skipImport {
		err := parseScan(args)
		if err != nil {
			return err
		}
	}
	return GetIpAsnCountries(args)
}
