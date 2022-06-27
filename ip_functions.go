package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

// getPublicIP returns public IPv4 and IPv6 IPs of host.
func getPublicIP(ipApiUrl string) string {
	// Get public IP address of host
	reqAddr, err := http.Get(ipApiUrl)
	if err != nil {
		log.Fatalf("Could not fetch public IP address: %v\n", err)
	}

	publicIPAddress, err := ioutil.ReadAll(reqAddr.Body)
	if err != nil {
		log.Fatalf("Could not read public IP response: %v\n", err)
	}

	return string(publicIPAddress)
}

// getDomainIP returns the A or AAAA records for a provided domain.
// recordType can be 'A' for IPv4 or 'AAAA' for IPv6
func getDomainIP(domain, recordType string) string {
	var ipV4, ipV6, dnsRecord string

	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Fatalf("Could not resolve IPs: %v\n", err)
	}

	for _, ip := range ips {
		ipString := ip.String()
		if strings.Count(ipString, ":") < 2 {
			ipV4 = ipString
		} else if strings.Count(ipString, ":") >= 2 {
			ipV6 = ipString
		}
	}

	if recordType == "A" {
		dnsRecord = ipV4
	} else if recordType == "AAAA" {
		dnsRecord = ipV6
	}

	return dnsRecord

}

// checkIPsMatch returns if two IP addresses match.
func checkIPsMatch(publicIPAddress, domainIPAddress string) bool {
	ipAddressesMatch := publicIPAddress == domainIPAddress

	return ipAddressesMatch

}
