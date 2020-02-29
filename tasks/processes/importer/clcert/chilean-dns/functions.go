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
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
)

type accessibleMap map[models.RRType]map[string]struct{}

// Returns a list of the accessible IPs, obtained from Zmap port scans.
// The files are located in the folder "basepath" and have the following
// structure as name: %s_port_%d.txt, where %s is the name of the scan and
// %d is the port scanned.
func getAccessibleIPs(source sources.Source, args *tasks.Context) (accessibleMap, error) {
	var accessible = make(accessibleMap)
	var nameType models.RRType
	for {
		file := source.Next()
		if file == nil {
			break
		}
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

		accessible[nameType] = make(map[string]struct{})
		args.Log.WithFields(logrus.Fields{
			"name":   file.Name(),
			"prefix": splittedName,
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
			accessible[nameType][strings.TrimSpace(line)] = struct{}{}
		}
	}
	return accessible, nil
}

// This function reads a Mercury results file in JSON and
// returns a list of DnsRRs representing each scan.
// It ignores "error" scans.
func getScanEntries(source sources.Source, saver savers.Saver, accessible accessibleMap, args *tasks.Context) error {
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
		err := readFile(file, saver, accessible, args, privateNetworks)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"name": file.Name(),
			}).Error("Error reading file: %s", err)
		}
	}
	return nil
}

func readFile(file sources.Entry, saver savers.Saver, accessible accessibleMap, args *tasks.Context, privateNetworks utils.NetList) error {
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
			// TODO: log in debug mode?
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
	accessible, err := getAccessibleIPs(ipsSource, args)
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
	distinctIPs := db.Model(&models.DnsRR{}).ColumnExpr("distinct task_id, ip_value as ip").Where("task_id = ?", args.GetTaskID())

	taskIDASN, err := models.LatestModelTaskID(db, &models.SubnetASN{})
	if err != nil {
		return err
	}

	taskIDCountry, err := models.LatestModelTaskID(db, &models.SubnetCountry{})
	if err != nil {
		return err
	}

	joinSubnetAsn := db.Model(&models.SubnetASN{}).
		Column("subnet", "asn_id", "source_id").
		Join("JOIN asns as asn").
		JoinOn("subnet_asns.asn_id = asn.id").
		JoinOn("subnet_asns.task_id = ?", taskIDASN)

	joinSubnetCountry := db.Model(&models.SubnetCountry{}).
		Column("subnet", "country_geoname_id").
		Join("JOIN countries as country").
		JoinOn("subnet_countries.country_geoname_id = country.geoname_id").
		JoinOn("subnet_countries.task_id = ?", taskIDCountry)

	ipAsnCountryList := make([]*models.IpAsnCountry, 0)
	query := db.Model().TableExpr("(?) as dr", distinctIPs).
		Column("dr.task_id", "dr.ip", "asns.asn_id", "asns.source_id", "countries.country_geoname_id").
		Join("JOIN (?) as asns", joinSubnetAsn).
		JoinOn("dr.ip << asns.subnet").
		Join("JOIN (?) as countries", joinSubnetCountry).
		JoinOn("dr.ip << countries.subnet")
	if err := query.Select(&ipAsnCountryList); err != nil {
		args.Log.Error("Err with the query :(")
		return err
	}
	args.Log.Info("Information acquired! now saving...")

	for _, ipAsnCountry := range ipAsnCountryList {
		if err := args.Savers[0].Save(ipAsnCountry); err != nil {
			return err
		}
	}

	return nil
}
