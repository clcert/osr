package domain_categories

import (
	"github.com/clcert/osr/tasks"
)

// Imports new domains from NIC webpage
func Execute(args *tasks.Context) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	err := ImportCategories(source, saver, args)
	if err != nil {
		return err
	}
	return nil
}
