package chilean_dns

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/clcert/osr/utils/geoasn"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
	"time"
)

type AccessibleMap map[models.RRType]map[string]struct{}
type AccessibleDateMap map[string]AccessibleMap

const DayFormat = "2006-01-02"

// Returns a list of the accessible IPs, obtained from Zmap port scans.
// The files are located in the folder "basepath" and have the following
// structure as name: %s_port_%d.txt, where %s is the name of the scan and
// %d is the port scanned.
func GetAccessibleIPs(source sources.Source, args *tasks.Context) (AccessibleDateMap, error) {
	var accessible = make(AccessibleDateMap)
	var nameType models.RRType
	for {
		file := source.Next()
		if file == nil {
			break
		}
		date := getFileDate(file)
		splittedName := strings.Split(file.Name(), "_port_")
		splittedPrefix := strings.Split(splittedName[0], "_")
		nameType = models.StringToRRType(splittedPrefix[len(splittedPrefix)-1])
		if nameType == models.NORR {
			args.Log.WithFields(logrus.Fields{
				"name":   file.Name(),
				"prefix": splittedName,
			}).Info("File without desired structure, skipping...")

			continue
		}
		if _, ok := accessible[date]; !ok {
			accessible[date] = make(map[models.RRType]map[string]struct{})
		}
		if _, ok := accessible[date][nameType]; !ok {
			accessible[date][nameType] = make(map[string]struct{})
		}
		args.Log.WithFields(logrus.Fields{
			"path":   file.Path(),
			"date": date,
		}).Info("Importing file")
		reader, err := file.Open()
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"name": file.Name(),
			}).Error("Failed reading file: %s", err)
			return accessible, err
		}
		bufReader := bufio.NewReader(reader)
		for {
			line, err := bufReader.ReadString('\n')
			if err != nil { // end of file
				if err == io.EOF {
					break
				} else {
					args.Log.WithFields(logrus.Fields{
						"error": err,
					}).Error("Err reading file")
					return accessible, err
				}
			}
			accessible[date][nameType][strings.TrimSpace(line)] = struct{}{}
		}
	}
	return accessible, nil
}

// This function reads a Mercury results file in JSON and
// returns a list of DnsRRs representing each scan.
// It ignores "error" scans.
func getScanEntries(source sources.Source, saver savers.Saver, accessibleDate AccessibleDateMap, args *tasks.Context) error {
	privateNetworks, err := utils.GetPrivateNetworks()
	if err != nil {
		args.Log.Error("Error getting private networks")
		return err
	}
	// for each file
	for {
		file := source.Next()
		if file == nil {
			break
		}
		date := getFileDate(file)
		accessible, ok := accessibleDate[date]
		if !ok {
			continue
		}
		err := readFile(file, saver, accessible , args, privateNetworks)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"name": file.Name(),
			}).Error("Error reading file: %s", err)
		}
	}
	return nil
}

func readFile(file sources.Entry, saver savers.Saver, accessible AccessibleMap, args *tasks.Context, privateNetworks utils.NetList) error {
	args.Log.WithFields(logrus.Fields{
		"name": file.Name(),
	}).Info("Importing file")
	reader, err := file.Open()
	if err != nil {
		return err
	}
	mercuryFile := utils.MercuryFromReader(reader, file.Name())
	args.Log.WithFields(logrus.Fields{
		"name": file.Name(),
	}).Info("Importing scan file")
	for {
		line, err := mercuryFile.NextResult()
		if err != nil {
			if err == io.EOF {
				args.Log.WithFields(logrus.Fields{
					"name": file.Name(),
				}).Info("Done reading file")
				return nil
			} else {
				args.Log.WithFields(logrus.Fields{
					"name": file.Name(),
				}).Error("Error reading file")
				continue
			}
		}
		var result DNSScanEntry
		if err = json.Unmarshal([]byte(line), &result); err != nil {
			err = fmt.Errorf("error unmarshaling line. File is probably not a Mercury DNS log:  %s", err)
			return err
		}
		if result.Status == ERROR {
			args.Log.WithFields(logrus.Fields{
				"message": result.Error,
				"domain":  result.Url,
			}).Error("Error when scanning")
			continue
		}
		subdomain, domain, tld, err := utils.SplitDomain(result.Url)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"result": result,
			}).Error("Malformed URL")
			continue
		}
		derivedType := GetDerivedType(file.Name())
		for i, rr := range result.RRs {
			dnsRR := &models.DnsRR{
				TaskID:          args.GetTaskID(),
				SourceID:        args.GetSourceID(),
				Date:            result.Timestamp,
				DomainSubdomain: subdomain,
				DomainName:      domain,
				DomainTLD:       tld,
				ScanType:        rr.Type,
				DerivedType:     derivedType,
				Index:           i,
			}
			switch rr.Type {
			case models.A:
				ip := net.ParseIP(rr.Value)
				if ip == nil {
					args.Log.WithFields(logrus.Fields{
						"Url":          result.Url,
						"Value":        rr.Value,
						"error":        err,
						"scan_type":    models.RRTypeToString(dnsRR.ScanType),
						"derived_type": models.RRTypeToString(dnsRR.DerivedType),
					}).Error("Malformed Address")
					dnsRR.Valid = false
				} else {
					dnsRR.IPValue = ip
					dnsRR.Valid = !privateNetworks.Contains(ip)
					if derivedType == models.NORR {
						_, dnsRR.Accessible = accessible[rr.Type][ip.String()]
					} else {
						_, dnsRR.Accessible = accessible[derivedType][ip.String()]
					}
				}
			case models.MX, models.NS, models.CNAME:
				valid := true
				vSubdomain, vDomain, vTLD, err := utils.SplitDomain(rr.Value)
				if err != nil {
					args.Log.WithFields(logrus.Fields{
						"Url":          result.Url,
						"Value":        rr.Value,
						"error":        err,
						"type":         models.RRTypeToString(dnsRR.ScanType),
						"derived_type": models.RRTypeToString(dnsRR.DerivedType),
					}).Error("Malformed DomainDomainCategory Value")
					valid = false
				}
				dnsRR.ValueSubdomain, dnsRR.ValueName, dnsRR.ValueTLD = vSubdomain, vDomain, vTLD
				dnsRR.Valid = valid
				dnsRR.Priority = rr.Priority
			}
			if err := saver.Save(dnsRR); err != nil {
				return err
			}
		}
	}
}

// Imports a folder to the scan service.
func parseScan(args *tasks.Context) error {
	mercurySource := args.Sources[0]
	ipsSource := args.Sources[1]
	args.Log.Info("Getting accessible IPs from scan...")
	accessible, err := GetAccessibleIPs(ipsSource, args)
	if err != nil {
		args.Log.Error("Problems with getting accessible IPs")
		return err
	}
	args.Log.Info("Getting RRs values and IPs involved in this scan...")
	err = getScanEntries(mercurySource, args.Savers[0], accessible, args)
	if err != nil {
		args.Log.Error("Problems with getting Scan Entries")
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

// Returns a list of IpAsnCountry, with the IPs of the last imports and the ASN and Country associated to them.
func GetIpAsnCountries(args *tasks.Context) error {
	args.Log.Info("Getting ASN and Country information of recently scanned IPs...")

	db, err := databases.GetPostgresWriter()
	if err != nil {
		return err
	}
	defer db.Close()
	dnsRRs := make([]*models.DnsRR, 0)
	taskID, err := models.LatestModelTaskID(db, &models.DnsRR{})
	if err != nil {
		return err
	}
	err = db.Model(&models.DnsRR{}).
		ColumnExpr("distinct ip_value").
		Where("task_id = ?", taskID).Select(&dnsRRs)
	if err != nil {
		return err
	}
	classifier, err := geoasn.NewClassifier(db, models.MaxMind)
	if err != nil {
		return err
	}

	args.Log.Info("Information acquired! now saving...")

	for _, dnsRR := range dnsRRs {
		countryID, asnID, err := classifier.GetGeoASN(dnsRR.IPValue)
		ipAsnCountry := &models.IpAsnCountry{
			TaskID:           args.GetTaskID(),
			SourceID:         args.GetSourceID(),
			IP:               dnsRR.IPValue,
			AsnID:            asnID,
			CountryGeonameId: countryID,
		}
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"ip": dnsRR.IPValue,
			}).Errorf("Cannot get ASN and Geoinfo from ip: %s", err)
			continue
		}
		if err := args.Savers[0].Save(ipAsnCountry); err != nil {
			return err
		}
	}

	return nil
}

func getFileDate(file sources.Entry) string {
	date := time.Now()
	splitPath := strings.Split(file.Path(), "/")
	if len(splitPath) >= 2 {
		fileDate, err := time.Parse(DayFormat, splitPath[len(splitPath)-2])
		if err == nil {
			date = fileDate
		}
	}
	// Now we know we can parse the value as a date, we return it as string
	return date.Format(DayFormat)
}