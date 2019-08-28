package chilean_dns

import (
	"encoding/json"
	"fmt"
	"github.com/clcert/osr/models"
	"net"
	"strings"
	"time"
)

// DNSScanEntry represents a raw scan line.
type DNSScanEntry struct {
	Url       string
	Status    ScanStatus
	Timestamp time.Time
	Resolver  net.IP
	Error     error
	RRs       []*RREntry `json:"answer_section"`
}

// RREntry represents a value in the RRs array of the scan entry.
type RREntry struct {
	Type     models.RRType
	Priority int // Only for MX
	Value    string
}

// Unmarshals a line in Mercury JSON.
func (r *DNSScanEntry) UnmarshalJSON(data []byte) error {
	var rawScanEntry map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawScanEntry); err != nil {
		return err
	}
	// Results["53"].dns_resolver
	var rawResults map[string]map[string]map[string]json.RawMessage
	if err := json.Unmarshal(rawScanEntry["Results"], &rawResults); err != nil {
		return err
	}
	dnsResolver := rawResults["53"]["dns_resolver"]
	var scanStatus string
	if err := json.Unmarshal(dnsResolver["status"], &scanStatus); err != nil {
		return err
	}
	r.Status = StringToScanStatus[scanStatus]
	if err := json.Unmarshal(dnsResolver["timestamp"], &r.Timestamp); err != nil {
		return err
	}
	answerSection, ok := dnsResolver["answer_section"]
	if ok {
		var url string
		if err := json.Unmarshal(rawScanEntry["url"], &url); err != nil {
			return err
		}
		r.Url = strings.Trim(strings.ToLower(url), ".")
		var rawAnswerSection map[string]json.RawMessage
		if err := json.Unmarshal(answerSection, &rawAnswerSection); err != nil {
			return err
		}
		var hostPort string
		if err := json.Unmarshal(dnsResolver["resolver"], &hostPort); err != nil {
			return err
		}
		host, _, err := net.SplitHostPort(hostPort)
		if err != nil {
			return err
		}
		r.Resolver = net.ParseIP(host)
		for ansType, ansValue := range rawAnswerSection {
			ansRRType := models.StringToRRType(ansType)
			switch ansRRType {
			case models.A, models.NS, models.CNAME:
				var rawAnswer []string
				if err := json.Unmarshal(ansValue, &rawAnswer); err != nil {
					return err
				}
				for _, value := range rawAnswer {
					r.RRs = append(r.RRs, &RREntry{
						Type:  ansRRType,
						Value: strings.Trim(strings.ToLower(value), "."),
					})
				}
			case models.MX:
				var rawAnswer []map[string]json.RawMessage
				if err := json.Unmarshal(ansValue, &rawAnswer); err != nil {
					return err
				}
				for _, value := range rawAnswer {
					var preference int
					json.Unmarshal(value["preference"], &preference) // If it doesn't exist, it's 0
					var domain string
					if err := json.Unmarshal(value["domain"], &domain); err != nil {
						return err
					}
					r.RRs = append(r.RRs, &RREntry{
						Type:     ansRRType,
						Priority: preference,
						Value:    strings.ToLower(strings.Trim(domain, ".")),
					})
				}
			}
		}
	} else {
		var errorMessage string
		if err := json.Unmarshal(dnsResolver["error"], &errorMessage); err != nil {
			return err
		}
		r.Error = fmt.Errorf("%s", errorMessage)
	}
	return nil
}

func GetDerivedType(filename string) models.RRType {
	nameExt := strings.Split(filename, ".")
	name := nameExt[0]
	splitTypes := strings.Split(name, "_")
	if len(splitTypes) == 2 {
		return models.StringToRRType(splitTypes[1])
	}
	return models.NORR
}

// ScanStatus represents the status of the DNSScanEntry
// It could be ERROR or SUCCESS
type ScanStatus int

const (
	ERROR ScanStatus = iota
	SUCCESS
)

// Transforms an DNSScanEntry status screen for the associated CONST.
var StringToScanStatus = map[string]ScanStatus{
	"error":   ERROR,
	"success": SUCCESS,
}
