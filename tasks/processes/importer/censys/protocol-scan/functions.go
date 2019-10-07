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
	date, port, protocol, err := parseMetaFromName(file.Name())
	if err != nil {
		args.Log.WithFields(logrus.Fields{
			"file_name": file.Name(),
			"current_date": date,
		}).Error("Couldn't determine date, port and protocol. Skipping file...")
		return err
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
	options := &ParserOptions{
		DefaultDate: date,
		Port:        port,
		Protocol:    protocol,
	}
	for scanner.Scan() {
		if err := insertWithParser(saver, scanner.Text(), options); err != nil {
			args.Log.WithFields(logrus.Fields{
				"file_path": file.Path(),
				"line":      scanner.Text(),
			}).Error("Error parsing line, skipping...")
		}
	}
	return file.Close()
}

func parseMetaFromName(name string) (date time.Time, port uint16, protocol string, err error) {
	slice1 := strings.Split(name, ".")
	slice2 := strings.Split(slice1[0], "_")
	date, err = time.Parse(slice2[0],"20060102")
	if err != nil {
		return
	}
	if len(slice2) >= 2 {
		port64, err := strconv.ParseInt(slice2[1], 10, 32)
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

// returns UDP if the port scanned is related to an UDP protocol.
func getTransport(port uint16) models.PortProtocol {
	switch port {
	case 53, 623, 1900, 20000, 47808:
		return models.UDP
	default:
		return models.TCP
	}
}


func insertWithParser(saver savers.Saver, line string, options *ParserOptions) error {
	if options == nil {
		return fmt.Errorf("option cannot be blank")
	}
	if _, ok := parsers[options.Protocol]; !ok {
		return fmt.Errorf("parser not found for protocol %s", options.Protocol)
	}
	return parsers[options.Protocol](saver, line, options)
}