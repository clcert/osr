package simple

import (
	"fmt"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"strconv"
	"sync"
)

func Execute(args *tasks.Args) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	threadsStr := args.Params.Get("threads", "1")
	threads, err := strconv.Atoi(threadsStr)
	if err != nil {
		return fmt.Errorf("cannot parse threads arg: %s", err)
	}
	wg := new(sync.WaitGroup)
	entries := make(chan sources.Entry)
	for i:=0; i<threads; i++ {
		wg.Add(1)
		go process(entries, saver, wg, args)
	}

	go func() {
		for {
			entry := source.Next()
			if entry == nil {
				close(entries)
				break
			}
			entries <- entry
		}
	}()
	wg.Wait()
	return nil
}
