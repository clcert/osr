// Package darknet reads a pcap file (compressed or not) and saves the data to a database and to a csv (optional).
package darknet

import (
	"fmt"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"sort"
	"strconv"
	"sync"
)

// Task obtains the current configuration and reads and process every file.
// It creates a channel of jobs containing the filenames.
// It also uses a channel for writing to the db.
// It returns an error if it fails to importer the data.
func Execute(args *tasks.Args) error {
	source := args.Sources[0]
	saver := args.Savers[0]

	fileChannel := make(chan sources.Entry)

	numWorkers := 1
	if workersStr, ok := args.Params["numWorkers"]; ok {
		confWorkers, err := strconv.Atoi(workersStr)
		if err == nil && confWorkers > 0 {
			numWorkers = confWorkers
		}
	}

	args.Log.
		Info(fmt.Sprintf("Starting %d workers\n", numWorkers))

	go func() {
		args.Log.Info("Making channel of files.")
		files := make([]sources.Entry, 0)
		for {
			file := source.Next()
			if file == nil {
				break
			}
			msg := fmt.Sprintf("Pushing file %s\n", file.Path())
			args.Log.
				Info(msg)
			files = append(files, file)
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Path() < files[j].Path()
		})
		for _, file := range files {
			fileChannel <- file
		}
		close(fileChannel)
	}()

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		args.Log.
			Info(fmt.Sprintf("Worker %d started.", i))
		go worker(i+1, &wg, fileChannel, saver, args)
	}
	wg.Wait()
	return nil
}
