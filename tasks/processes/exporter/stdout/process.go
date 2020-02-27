package stdout

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/sources"
	"github.com/sirupsen/logrus"
	"sync"
)

func process(entryChan chan sources.Entry, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range entryChan {
		logs.Log.WithFields(logrus.Fields{
			"file": file.Path(),
		}).Info("opening file")
		reader, err := file.Open()
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"file": file.Path(),
			}).Error("couldn't open file: %s", err)
			continue
		}
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}
	}
	logs.Log.WithFields(logrus.Fields{
	}).Info("thread done!")
}
