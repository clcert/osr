package protocol_scan

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/censys"
	"github.com/clcert/osr/utils/protocols"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
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
	port, protocol, err := parseMeta(file)
	if err != nil {
		args.Log.WithFields(logrus.Fields{
			"file_name": file.Name(),
		}).Error("%s. Skipping file...", err)
		return err
	}
	reader, err := file.Open()
	if err != nil {
		return err
	}
	args.Log.WithFields(logrus.Fields{
		"file_path": file.Path(),
	}).Info("File opened")
	scanner := bufio.NewScanner(reader)
	options := &censys.ParserOptions{
		Port:     port,
		Protocol: protocol,
	}
	parser, ok := protocols.Parsers[protocol]
	if !ok {
		return fmt.Errorf("couldn't find parser for this protocol: %s", protocol)
	}
	for scanner.Scan() {
		entry, err := options.Unmarshal(scanner.Text())
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
				"line":      scanner.Text(),
			}).Error("Error unmarshaling line, skipping...")
			continue
		}
		if err := entry.GetError(); err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
				"line":      scanner.Text(),
				"error":     err,
			}).Error("This line is not well formed, skipping...")
			continue
		}
		scanDate := entry.GetTime(censys.DateFormat, time.Time{})
		portScan := &models.PortScan{
			TaskID:     args.Task.ID,
			SourceID:   args.Process.Source,
			Date:       scanDate,
			ScanIP:     srcIP,
			IP:         entry.GetIP(),
			PortNumber: port,
			Protocol:   protocols.GetTransport(port),
		}
		if parser.IsValid(entry.GetBanner()) {
			software, version := parser.GetSoftwareAndVersion(entry.GetBanner())
			portScan.ServiceActive = true
			portScan.ServiceName = software
			portScan.ServiceVersion = version
			portScan.ServiceExtra = entry.GetBanner()
		}
		if err = saver.Save(portScan); err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
				"ip":        entry.GetIP(),
				"port":      port,
			}).Error("Couldn't save entry: %s", err)
			continue
		}
	}
	return file.Close()
}

func parseMeta(source sources.Entry) (port uint16, protocol string, err error) {
	slice1 := strings.Split(source.Name(), ".")
	slice2 := strings.Split(slice1[0], "_")
	if len(slice2) >= 2 {
		var port64 int64
		port64, err = strconv.ParseInt(slice2[1], 10, 32)
		if err != nil {
			return
		}
		port = uint16(port64)
	}
	if len(slice2) >= 3 {
		protocol = slice2[2]
	}
	return
}
