package chilean_subnets

import (
	"encoding/json"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
)

func saveRIPE(entry sources.Entry, saver savers.Saver, args *tasks.Args, country *models.Country) error {
	resources, err := parseRIPE(entry)
	if err != nil {
		args.Log.WithFields(logrus.Fields{
			"entry": entry.Name(),
			"error": err,
		}).Error("cannot parse as JSON.")
		return err
	}
	ips, ok := resources["ipv4"]
	if !ok {
		args.Log.WithFields(logrus.Fields{
			"entry": entry.Name(),
		}).Error("data.resources.ipv4 missing")
		return err
	}
	args.Log.Info("Parsing RIPE country Subnets...")
	for _, subnetStr := range ips {
		_, subnet, err := net.ParseCIDR(subnetStr)
		if err != nil {
			args.Log.WithFields(logrus.Fields{
				"subnet": subnetStr,
			}).Error("error parsing this subnet")
		}
		if err := saver.Save(&models.SubnetCountry{
			TaskID:           args.GetTaskID(),
			SourceID:         args.GetSourceID(),
			Subnet:           subnet,
			CountryGeonameId: country.GeonameId,
		}); err != nil {

		}
	}
	return nil
}

// Parses RIPE json and returns a map with resources (ipv4, asn, ipv6). Each resource is a list
// with string values.
func parseRIPE(entry sources.Entry) (map[string][]string, error) {

	reader, err := entry.Open()
	if err != nil {
		return nil, err
	}

	rawJson, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	allResources := make(map[string][]string)

	root := make(map[string]interface{})
	err = json.Unmarshal(rawJson, &root)
	if err != nil {
		return nil, err
	}
	data, ok := root["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data property not found")
	}

	resources, ok := data["resources"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data.resources property not found")
	}

	for key, resource := range resources {
		resourceList, ok := resource.([]interface{})
		if !ok {
			continue
		}
		allResources[key] = make([]string, len(resourceList))
		for i, resourceElem := range resourceList {
			allResources[key][i] = resourceElem.(string)
		}
	}

	return allResources, nil

}
