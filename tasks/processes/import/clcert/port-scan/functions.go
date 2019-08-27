package port_scan

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
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
	regex, err := NewRegex(args)
	if err != nil {
		args.Log.WithFields(logrus.Fields{
			"params": args.Params,
		}).Error("couldn't create regexes: %s", err)
	}
	for {
		file := source.Next()
		if file == nil {
			return nil
		}
		args.Log.WithFields(logrus.Fields{
			"files_read": filesRead,
		}).Info("Reading files...")
		err := parseFile(file, saver, args, srcIP, regex)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
			}).Errorf("Error reading file: %s", err)
		}
		filesRead++
	}
}

func parseFile(file sources.Entry, saver savers.Saver, args *tasks.Args, srcIP net.IP, regex *Regexes) error {
	date := parseDate(file.Dir())
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
		port, err = regex.GetPort(file.Name())
		if err != nil {
			file.Close()
			return fmt.Errorf("couldn't get port number from filename: %s", err)
		}
	}
	protocol := regex.GetProtocol(file.Path())
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

func parseDate(dir string) (date time.Time) {
	var err error
	dirSlice := strings.Split(dir, "/")
	for i := len(dirSlice) - 1; i >= 0; i-- {
		date, err = time.Parse(DateFolderFormat, dirSlice[i])
		if err == nil && !date.IsZero() {
			break
		}
	}
	if err != nil {
		date = time.Now()
	}
	return
}
