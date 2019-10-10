package port_scan

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/grabber"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"time"
)

func parseFiles(source sources.Source, saver savers.Saver, args *tasks.Args) error {
	srcAddr, err := source.GetID()
	if err != nil {
		return err
	}
	srcIPStr := strings.Split(srcAddr, ":")[0]
	srcIP := net.ParseIP(srcIPStr)
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
	date, err := grabber.ParseDate(DateFormat, file.Dir())
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
	var port uint16
	if portStr, ok := args.Params["port"]; ok {
		port64, err := strconv.ParseUint(portStr, 10, 16)
		if err == nil {
			port = uint16(port64)
		}
	}
	if port == 0 {
		port, err = grabber.ParsePort(file.Name())
		if err != nil {
			file.Close()
			return fmt.Errorf("couldn't get port number from filename: %s", err)
		}
	}
	protocol := grabber.ParseProtocol(file.Path())
	for scanner.Scan() {
		ip := net.ParseIP(scanner.Text())
		if ip == nil {
			// Line is not an IP. Logging this event and continuing
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
				"line":      scanner.Text(),
			}).Error("line is not an IP. Skipping...")
			continue
		}
		if err = saver.Save(&models.PortScan{
			TaskID:     args.Task.ID,
			PortNumber: port,
			SourceID:   args.Process.Source,
			ScanIP:     srcIP,
			IP:         ip,
			Date:       date,
			Protocol:   protocol,
		}); err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
				"line":      scanner.Text(),
			}).Error("Couldn't save line: %s", err)
			continue
		}
	}
	return file.Close()
}

