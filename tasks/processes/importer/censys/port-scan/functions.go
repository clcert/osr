package port_scan

import (
	"bufio"
	"encoding/json"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/protocols"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

func parseFiles(source sources.Source, saver savers.Saver, args *tasks.Args) error {
	srcIP := net.ParseIP("216.239.34.21") // Censys.io IP. NOT SCANNING IP!
	filesRead := 0
	for {
		file := source.Next()
		if file == nil {
			return nil
		}
		args.Log.WithFields(logrus.Fields{
			"files_read": filesRead,
		}).Info("Reading files...")
		err := parseFile(file, saver, args, srcIP)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
			}).Errorf("Error reading file: %s", err)
		}
		filesRead++
	}
}

func parseFile(file sources.Entry, saver savers.Saver, args *tasks.Args, srcIP net.IP) error {
	date, err := parseDate(file.Name())
	if err != nil {
		date = time.Now()
		args.Log.WithFields(logrus.Fields{
			"file_name": file.Name(),
			"current_date": date,
		}).Error("Couldn't determine date. Using current date and time...")
	}
	reader, err := file.Open()
	if err != nil {
		return err
	}
	args.Log.WithFields(logrus.Fields{
		"file_path": file.Path(),
		"date":      date,
	}).Info("File opened")
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := Entry{}
		err := json.Unmarshal([]byte(scanner.Text()), &line)
		if err != nil {
			// Line is not an IP. Logging this event and continuing
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
				"line":      scanner.Text(),
			}).Error("line is not a valid object. Skipping...")
			continue
		}
		ip := net.ParseIP(line.IP)
		for _, port := range line.Ports {
			if err = saver.Save(&models.PortScan{
				TaskID:     args.Task.ID,
				PortNumber: port,
				SourceID:   args.Process.Source,
				ScanIP:     srcIP,
				IP:         ip,
				Date:       date,
				Protocol:   protocols.GetTransport(port),
			}); err != nil {
				args.Log.WithFields(logrus.Fields{
					"file_path": file.Path(),
					"ip": ip,
					"port": port,
				}).Error("Couldn't save entry: %s", err)
				continue
			}
		}
	}
	return file.Close()
}

func parseDate(dir string) (date time.Time, err error) {
	dirSlice := strings.Split(dir, ".")
	return time.Parse(DateFormat, dirSlice[0])
}

