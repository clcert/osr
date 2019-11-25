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
	"github.com/clcert/osr/utils/filters"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func parseFile(file sources.Entry, args *tasks.Args, conf *filters.ScanConfig) error {
	saver := args.Savers[0]
	date, port, protocol, err := parseMeta(file)
	if err != nil {
		args.Log.WithFields(logrus.Fields{
			"file_name":    file.Name(),
			"current_date": date,
		}).Error("Couldn't determine date, port and protocol. Skipping file...")
		return err
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
	if !conf.IsNotInBlacklist(port) {
		args.Log.WithFields(logrus.Fields{
			"file_path": file.Path(),
			"port":      port,
		}).Info("Skipping file because port is on blacklist")
		return nil
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
		if err := entry.GetError(); err != nil {
			continue
		}
		cert, err := entry.GetCertificate()
		if err == nil {
			if err = saver.Save(&models.Certificate{
				TaskID:             args.GetTaskID(),
				SourceID:           args.GetSourceID(),
				Date:               entry.GetTime(censys.DateFormat, options.DefaultDate).Local(),
				ScanIP:             conf.SourceIP,
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
		portScan := &models.PortScan{
			TaskID:     args.GetTaskID(),
			SourceID:   args.GetSourceID(),
			Date:       entry.GetTime(censys.DateFormat, options.DefaultDate),
			ScanIP:     conf.SourceIP,
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

func parseMeta(file sources.Entry) (date time.Time, port uint16, protocol string, err error) {
	date, err = grabber.ParseDate("2006-01-02", file.Dir())
	if err != nil {
		return
	}
	date = date.In(time.Local)
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
