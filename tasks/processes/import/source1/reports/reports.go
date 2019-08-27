package reports

import (
	"fmt"
	"github.com/clcert/osr/models"
	"regexp"
	"strings"
)

var botsExtraStructure = "(?:srcport (\\d+) )?mwtype ([^\\s]+) (?:destaddr ([0-9.]+))?(?: dsthost: ([\\w.]+))?(?:(.*))?"
var botsExtraRegexp *regexp.Regexp

func newReport(extra string) map[string]string {
	props := make(map[string]string)
	if len(extra) > 0 {
		vals := strings.Split(extra, ";")
		for _, val := range vals {
			if len(val) > 0 {
				keyValueArr := strings.Split(val, ": ")
				if len(keyValueArr) == 2 {
					props[keyValueArr[0]] = keyValueArr[1]
				}
			}
		}
	}
	return props
}

func newBots(extra string) (models.ReportTypeID, map[string]string, error) {
	tags := []string{"port", "family", "dstaddr", "dsthost", "extra"}
	props := make(map[string]string)
	if len(extra) == 0 { // empty extra
		return models.BotReport, props, nil
	}
	// TODO make this thread-safe
	if botsExtraRegexp == nil {
		var err error
		botsExtraRegexp, err = regexp.Compile(botsExtraStructure)
		if err != nil {
			return models.UnknownReport, nil, err
		}
	}
	vals := botsExtraRegexp.FindStringSubmatch(extra)
	if len(vals) != len(tags)+1 { // first match is entire string
		return models.UnknownReport, nil, fmt.Errorf("tags length is not equal to acquired values")
	}
	for i, tag := range tags {
		val := vals[i+1]
		if len(val) > 0 {
			props[strings.ToLower(tag)] = strings.ToLower(val)
		}
	}
	return models.BotReport, props, nil
}

func newBot(extra string) (models.ReportTypeID, map[string]string, error) {
	return models.BotReport, newReport(extra), nil
}

func newBruteForce(extra string) (models.ReportTypeID, map[string]string, error) {
	return models.BruteforceReport, map[string]string{"protocol": extra}, nil
}

func newC2(extra string) (models.ReportTypeID, map[string]string, error) {
	return models.C2Report, newReport(extra), nil

}

func newDarknet(extra string) (models.ReportTypeID, map[string]string, error) {
	return models.DarknetReport, newReport(extra), nil
}

func newHoneypot(extra string) (models.ReportTypeID, map[string]string, error) {
	return models.HoneypotReport, newReport(extra), nil
}

func newDNSResolver(extra string) (models.ReportTypeID, map[string]string, error) {
	return models.DNSResolverReport, newReport(extra), nil
}

func newPhishing(extra string) (models.ReportTypeID, map[string]string, error) {
	return models.PhishingReport, map[string]string{"url": extra}, nil

}

func newProxy(extra string) (models.ReportTypeID, map[string]string, error) {
	// three notations:
	// - <PROTOCOL> (<PORT>)
	// - <PROTOCOL>-<PORT>;
	// - proxy_type: <PROTOCOL>[-<PORT>];

	extraArr := strings.Split(extra, " ")
	protocolStr := strings.Join(extraArr[:len(extraArr)-1], " ")
	portStr := extraArr[len(extraArr)-1]

	// last one is the port. It has parenthesis rounding it
	props := map[string]string{
		"protocol": strings.ToLower(protocolStr),
		"port":     strings.Trim(portStr, "();"),
	}
	return models.ProxyReport, props, nil
}

func newSpam(extra string) (models.ReportTypeID, map[string]string, error) {
	return models.SpamReport, newReport(extra), nil
}
