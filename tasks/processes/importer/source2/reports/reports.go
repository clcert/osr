package reports

import (
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils/filters"
	"net"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

type ReportList map[string]*Report

var fileToReport = make(ReportList)

type Report struct {
	Filename       string
	Type           models.ReportTypeID
	DeletedFields  []string
	RemappedFields map[string]string
}

func (list ReportList) Add(report *Report) {
	list[report.Filename] = report
}

func (list ReportList) Exists(name string) bool {
	_, ok := list[name]
	return ok
}

func (list ReportList) Parse(filename string, entry map[string]string, args *tasks.Args, conf *filters.DateConfig) (*models.ReportEntry, error) {
	report, ok := list[filename]
	if !ok {
		return nil, fmt.Errorf("cannot get report type")
	}
	// Get Date
	timestampString, ok := entry["timestamp"]
	if !ok {
		return nil, fmt.Errorf("cannot get timestamp")
	}
	date, err := time.Parse(timeLayout, timestampString)
	if err != nil {
		return nil, fmt.Errorf("cannot parse timestamp")
	}
	if !conf.IsDateInRange(date) {
		return nil, fmt.Errorf("line is not in date range")
	}
	// get IP
	ipString, ok := entry["ip"]
	if !ok {
		return nil, fmt.Errorf("cannot get IP")
	}
	ip := net.ParseIP(ipString)

	// we remove the unused general fields from entry and marshal them
	report.getExtra(entry)
	return &models.ReportEntry{
		ReportTypeID: report.Type,
		SourceID:     args.GetSourceID(),
		Date:         date,
		IP:           ip,
		TaskID:       args.GetTaskID(),
		Properties:   entry,
	}, nil
}

func (report *Report) getExtra(entry map[string]string) {
	// Deleting all repeated fields
	for _, field := range report.DeletedFields {
		delete(entry, field)
	}
	for orig, modified := range report.RemappedFields {
		value, ok := entry[orig]
		if ok {
			delete(entry, orig)
			entry[modified] = value
		}
	}
}

func init() {
	fileToReport.Add(botnet)
	fileToReport.Add(bruteForce)
	fileToReport.Add(sinkhole)
	fileToReport.Add(microsoftSinkhole)
	fileToReport.Add(darknet)
}

var botnet = &Report{
	Filename: "botnet_drone-chile-geo.csv",
	Type:     models.BotReport,
	DeletedFields: []string{
		"timestamp",
		"ip",
		"asn",
		"geo",
		"region",
		"city",
		"cc_asn",
		"cc_geo",
		"naics",
		"sic",
		"sector",
		"cc_naics",
		"cc_sic",
		"cc_sector",
	},
	RemappedFields: map[string]string{
		"type":      "protocol",
		"infection": "family",
	},
}

var bruteForce = &Report{
	Filename: "drone_brute_force-chile-geo.csv",
	Type:     models.BruteforceReport,
	DeletedFields: []string{
		"timestamp",
		"ip",
		"asn",
		"geo",
		"region",
		"city",
		"dest_asn",
		"dest_geo",
		"naics",
		"dest_naics",
		"sic",
		"dest_sic",
		"sector",
		"dest_sector",
	},
	RemappedFields: map[string]string{

	},
}

var sinkhole = &Report{
	Filename: "sinkhole_http_drone-chile-geo.csv",
	Type:     models.BotReport,
	DeletedFields: []string{
		"timestamp",
		"ip",
		"asn",
		"geo",
		"http_referer_asn",
		"http_referer_geo",
		"dst_asn",
		"dst_geo",
	},
	RemappedFields: map[string]string{
		"type":     "drone_type",
		"src_port": "port",
		"dst_ip":   "dest_ip",
	},
}

var microsoftSinkhole = &Report{
	Filename: "microsoft_sinkhole-chile-geo.csv",
	Type:     models.BotReport,
	DeletedFields: []string{
		"timestamp",
		"ip",
		"asn",
		"geo",
		"http_referer_asn",
		"http_referer_geo",
		"dst_asn",
		"dst_geo",
	},
	RemappedFields: map[string]string{
		"type":     "drone_type",
		"src_port": "port",
		"dst_ip":   "dest_ip",
	},
}

var darknet = &Report{
	Filename: "darknet-chile-geo.csv",
	Type:     models.DarknetReport,
	DeletedFields: []string{
		"timestamp",
		"ip",
		"asn",
		"dst_asn",
		"geo",
		"dst_geo",
		"region",
		"city",
		"naics",
		"dst_naics",
		"sic",
		"dst_sic",
		"sector",
		"dst_sector",
	},
	RemappedFields: map[string]string{
		"dst_ip":   "dest_ip",
		"dst_port": "dest_port",
	},
}
