package port_scan

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/grabber"
	"github.com/clcert/osr/utils/scans"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"time"
)

func parseFiles(source sources.Source, saver savers.Saver, args *tasks.Args) error {
	srcAddr, err := source.GetID()
	srcIPStr := strings.Split(srcAddr, ":")[0]
	srcIP := net.ParseIP(srcIPStr)
	if err != nil {
		return err
	}
	conf, errs := scans.ParseConf(args.Params, srcIP)
	for err := range errs {
		args.Log.Errorf("Error parsing config: %s", err)
	}
	filesRead := 0
	for {
		file := source.Next()
		if file == nil {
			return nil
		}
		args.Log.WithFields(logrus.Fields{
			"files_read": filesRead,
		}).Info("Reading files...")
		err := parseFile(file, saver, args, conf)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
			}).Errorf("Error reading file: %s", err)
		}
		filesRead++
	}
}

func parseFile(file sources.Entry, saver savers.Saver, args *tasks.Args, conf *scans.ScanConfig) error {
	date, err := grabber.ParseDate(DateFormat, file.Dir())
	if err != nil {
		date = time.Now()
		args.Log.WithFields(logrus.Fields{
			"file_name":    file.Name(),
			"current_date": date,
		}).Error("Couldn't determine date. Using current date and time...")
	}
	if !conf.IsDateInRange(date) {
		args.Log.WithFields(logrus.Fields{
			"since":   conf.Since,
			"until":   conf.Until,
			"curDate": date,
			"file":    file.Name(),
		}).Error("Ignoring file because is outside the date interval")
		return nil
	}
	port, err := getPort(args, file.Name())
	if err != nil {
		file.Close()
		return err
	}
	if !conf.IsPortAllowed(port) {
		args.Log.WithFields(logrus.Fields{
			"file_path": file.Path(),
			"port":      port,
		}).Info("Skipping file because port is on blacklist")
		return nil
	}
	protocol := grabber.ParseProtocol(file.Path())
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
			ScanIP:     conf.SourceIP,
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

func getPort(args *tasks.Args, filename string) (port uint16, err error) {
	if portStr, ok := args.Params["port"]; ok {
		port64, err := strconv.ParseUint(portStr, 10, 16)
		if err == nil {
			port = uint16(port64)
		}
	}
	if port == 0 {
		port, err = grabber.ParsePort(filename)
		if err != nil {
			err = fmt.Errorf("couldn't get port number from filename: %s", err)
			return
		}
	}
	return
}
