package domains

import (
	"github.com/clcert/osr/tasks"
)

const NicTimeLayout = "2006-01-02 15:04:05"

// Export imports NIC Chile data.
func Execute(args *tasks.Context) error {
	saver := args.Savers[0]
	for _, source := range args.Sources {
		if err := processSource(source, saver, args); err != nil {
			return err
		}
	}
	return nil
}
