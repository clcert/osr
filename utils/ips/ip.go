package ips

import (
	"net"
)

func CopyIP(ip net.IP) (ipCopy net.IP) {
	newIP := make(net.IP,len(ip))
	copy(newIP, ip)
	return newIP
}


// SubToIP substracts a specific integer to the IP Address.
func SubIP(ip net.IP, i int) net.IP {
	return AddIP(ip, -i)
}

// AddToIP adds a specific integer to the IP address.
func AddIP(ip net.IP, i int) net.IP {
	if ip == nil {
		return nil
	}
	val := i;
	newIP := make(net.IP, len(ip))
	copy(newIP, ip)
	for i := len(ip) - 1; i >= 0; i-- {
		newIP[i] = byte((int(ip[i]) + val) % 256)
		val = (int(ip[i]) + val) / 256
		if val == 0 {
			break
		}
	}
	return newIP
}
