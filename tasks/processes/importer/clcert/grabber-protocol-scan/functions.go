package grabber_protocol_scan

import (
	"bufio"
	"fmt"
	"github.com/clcert/osr/models"
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

func parseFile(file sources.Entry, args *tasks.Args, srcIP net.IP) error {
	saver := args.Savers[0]
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
	options := &grabber.ParserOptions{
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
		if entry.GetError() != nil {
			continue
		}
		cert, err := entry.GetCertificate()
		if err == nil {
			if err = saver.Save(&models.Certificate{
				TaskID:             args.Task.ID,
				SourceID:           args.Process.Source,
				Date:               entry.GetTime(censys.DateFormat, options.DefaultDate).Local(),
				ScanIP:             srcIP,
				IP:                 entry.GetIP(),
				PortNumber:         port,
				IsAutosigned:       cert.IsAutosigned(),
				KeySize:            cert.GetKeySize(),
				ExpirationDate:     cert.GetExpirationDate(),
				OrganizationName:   cert.GetOrganizationName(),
				OrganizationURL:    cert.GetOrganizationURL(),
				Authority:          cert.GetAuthority(),
				SignatureAlgorithm: cert.GetSigAlgorithm(),
				TLSProtocol:        cert.GetTLSProtocol(),
			}); err != nil {
				args.Log.WithFields(logrus.Fields{
					"file_path": file.Path(),
					"ip":        entry.GetIP(),
					"port":      port,
				}).Error("Couldn't save cert: %s", err)
				continue
			}
		}
		if strings.Contains(file.Path(), "certificate") {
			continue // File contains only certificates
		}
		software, version := parser.GetSoftwareAndVersion(entry.GetBanner())
		if err = saver.Save(&models.PortScan{
			TaskID:         args.Task.ID,
			SourceID:       args.Process.Source,
			Date:           entry.GetTime(censys.DateFormat, options.DefaultDate),
			ScanIP:         srcIP,
			IP:             entry.GetIP(),
			PortNumber:     port,
			Protocol:       protocols.GetTransport(port),
			ServiceActive:  parser.IsValid(entry.GetBanner()),
			ServiceName:    software,
			ServiceVersion: version,
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
	date, err = grabber.ParseDate("2006-01-02", file.Dir())
	if err != nil {
		return
	}
	date = 	date.In(time.Local)
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
