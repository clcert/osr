package stdout

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
)

func Execute(args *tasks.Context) error {
	source := args.Sources[0]
	for {
		file := source.Next()
		if file == nil {
			break
		}
		logs.Log.WithFields(logrus.Fields{
			"file": file.Path(),
		}).Info("opening file")
		reader, err := file.Open()
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"file": file.Path(),
			}).Error("couldn't open file: %s", err)
			file.Close()
			return err
		}
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}
