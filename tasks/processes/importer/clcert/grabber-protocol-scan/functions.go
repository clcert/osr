package grabber_protocol_scan

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/censys"
	"github.com/clcert/osr/utils/grabber"
	"github.com/clcert/osr/utils/protocols"
	"github.com/sirupsen/logrus"
	"net"
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
	date, port, protocol, err := parseMeta(file)
	if err != nil {
		args.Log.WithFields(logrus.Fields{
			"file_name":    file.Name(),
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
	options := &censys.ParserOptions{
		DefaultDate: date,
		Port:        port,
		Protocol:    protocol,
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
		if err = saver.Save(&models.PortScan{
			TaskID:         args.Task.ID,
			SourceID:       args.Process.Source,
			Date:           entry.GetTime(censys.DateFormat, options.DefaultDate),
			ScanIP:         srcIP,
			IP:             entry.GetIP(),
			PortNumber:     port,
			Protocol:       getTransport(port),
			ServiceActive:  parser.IsValid(entry.GetBanner()),
			ServiceName:    parser.GetSoftware(entry.GetBanner()),
			ServiceVersion: parser.GetVersion(entry.GetBanner()),
			ServiceExtra:   entry.GetBanner(),
		}); err != nil {
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

func parseMeta(file sources.Entry) (date time.Time, port uint16, protocol string, err error) {
	date, err = grabber.ParseDate(file.Path(), "20060102")
	if err != nil {
		return
	}
	port, err = grabber.ParsePort(file.Name())
	if err != nil {
		return
	}
	protocol, ok := protocols.PortToProtocol[port]
	if !ok {
		err = fmt.Errorf("unknown protocol for port %s", port)
		return
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
