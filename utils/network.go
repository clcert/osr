package utils

import (
	"fmt"
	"net"
	"strings"
)

// Defines a list of IPs.
type NetList []*net.IPNet

// Defines a list of CIDR values that represent the invalid public networks
// in IPv4 protocol.
var PrivateNetworks = []string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"192.88.99.0/24",
	"192.168.0.0/16",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"255.255.255.255/32",
}

// Returns a net list with all the private asns in IPv4 protocol.
func GetPrivateNetworks() (PrivateIPNet NetList, err error) {
	for _, ip := range PrivateNetworks {
		_, newNet, err := net.ParseCIDR(ip)
		if err != nil {
			return PrivateIPNet, err
		}
		PrivateIPNet = append(PrivateIPNet, newNet)
	}
	return PrivateIPNet, nil
}

// Checks if an Address is contained by any of the values in a NetList
func (nl NetList) Contains(ip net.IP) bool {
	for _, ipnet := range nl {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

// Splits a domain into tld, domain and a string with zero or more subdomains.
func SplitDomain(url string) (string, string, string, error) {
	spUrl := strings.Split(url, ".")
	var subdomain, domain, tld string
	if len(spUrl) < 2 {
		return "", "", "", fmt.Errorf("malformed URL")
	} else if len(spUrl) == 2 {
		domain, tld = spUrl[0], spUrl[1]
	} else {
		lenSpUrl := len(spUrl)
		subdomain, domain, tld = strings.Join(spUrl[:lenSpUrl-2], "."), spUrl[lenSpUrl-2], spUrl[lenSpUrl-1]
	}
	return subdomain, domain, tld, nil
}
