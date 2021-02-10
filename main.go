package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ipData struct {
	publicIP4Api string
	publicIP6Api string
	domainName   string
	publicIPv4   string
	publicIPv6   string
	domainIPv4   string
	domainIPv6   string
}

// Function to get public IPv4 and IPv6 IPs of host
func getPublicIP(ip4Url, ip6Url string) (string, string) {
	// Get public IPv4 address of host
	reqV4, err := http.Get(ip4Url)
	if err != nil {
		log.Fatalf("Could not get public IPv4 address: %v\n", err)
	}

	ipV4, err := ioutil.ReadAll(reqV4.Body)
	if err != nil {
		log.Fatalf("Could not read public IPv4 response: %v\n", err)
	}
	// Get public IPv6 address of host
	reqV6, err := http.Get(ip6Url)
	if err != nil {
		log.Fatalf("Could not get public IPv6 address: %v\n", err)
	}

	ipV6, err := ioutil.ReadAll(reqV6.Body)
	if err != nil {
		log.Fatalf("Could not read public IPv6 response: %v\n", err)
	}

	return string(ipV4), string(ipV6)
}

func getDomainIP(domain string) (string, string) {
	// Will be used to return the resolved IPs of the domain. If domain not resolved and user wants to create new record,
	// will return empty string so that public IPs of machine will be used for creating the record.
	var ipV4, ipV6 string

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

	return ipV4, ipV6

}

func checkIPsMatch(publicIPv4, domainIPv4, publicIPv6, domainIPv6 string) (bool, bool) {
	ip4AddressesMatch := publicIPv4 == domainIPv4

	ip6AddressesMatch := publicIPv6 == domainIPv6

	return ip4AddressesMatch, ip6AddressesMatch

}

func main() {
	// Config setup
	viper.SetConfigName("ddns_client_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/ddnsclient/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Fatal error reading config file: default \n", err)
		os.Exit(1)
	}

	var ipInfo ipData
	ipInfo.publicIP4Api = viper.GetString("app.publicIP4Api")
	ipInfo.publicIP6Api = viper.GetString("app.publicIP6Api")
	ipInfo.domainName = viper.GetString("app.domainName")

	gcpInfo := &gcpData{
		projectName: viper.GetString("gcpDNS.projectName"),
		zoneName:    viper.GetString("gcpDNS.zoneName"),
		ttlValue:    viper.GetInt64("gcpDNS.ttlValue"),
	}

	// End config setup

	log.Println("Dynamic DNS client - Starting check of public IPs")

	// Get public IPs of host
	ipInfo.publicIPv4, ipInfo.publicIPv6 = getPublicIP(ipInfo.publicIP4Api, ipInfo.publicIP6Api)
	log.Printf("Hosts public IPv4: %s and\nIPv6: %s", ipInfo.publicIPv4, ipInfo.publicIPv6)

	// Resolve IPs of provided domain
	ipInfo.domainIPv4, ipInfo.domainIPv6 = getDomainIP(ipInfo.domainName)

	// Check if public and resolved IPs are the same
	IPv4Same, IPv6Same := checkIPsMatch(ipInfo.publicIPv4, ipInfo.domainIPv4, ipInfo.publicIPv6, ipInfo.domainIPv6)

	if IPv4Same == false {
		UpdateDNSRecord(gcpInfo.projectName, gcpInfo.zoneName, ipInfo.domainName, ipInfo.domainIPv4, ipInfo.publicIPv4, gcpInfo.ttlValue)
	} else {
		log.Println("Public IPv4 address has not changed.")
	}

	if IPv6Same == false {
		UpdateDNSRecord(gcpInfo.projectName, gcpInfo.zoneName, ipInfo.domainName, ipInfo.domainIPv6, ipInfo.publicIPv6, gcpInfo.ttlValue)
	} else {
		log.Println("Public IPv6 address has not changed.")
	}

	os.Exit(0)

}
