package main

import (
	"fmt"
	"io/ioutil"
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
		fmt.Fprintf(os.Stderr, "Could not get public IPv4 address: %v\n", err)
		os.Exit(1)
	}

	ipV4, err := ioutil.ReadAll(reqV4.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read public IPv4 response: %v\n", err)
		os.Exit(1)
	}
	// Get public IPv6 address of host
	reqV6, err := http.Get(ip6Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get public IPv6 address: %v\n", err)
		os.Exit(1)
	}

	ipV6, err := ioutil.ReadAll(reqV6.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read public IPv6 response: %v\n", err)
		os.Exit(1)
	}

	return string(ipV4), string(ipV6)
}

func getDomainIP(domain string) (string, string) {
	// Will be used to return the resolved IPs of the domain. If domain not resolved and user wants to create new record,
	// will return empty string so that public IPs of machine will be used for creating the record.
	var ipV4, ipV6 string

	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("Could not resolve IPs. Do you want to create DNS records for the domain: %s?\n(Yes/No)\n", domain)
		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "Yes" {
			return ipV4, ipV6
		}

		fmt.Fprintf(os.Stderr, "Could not resolve IPs: %v\n", err)
		os.Exit(1)

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

func compareIPs(publicIPv4, domainIPv4, publicIPv6, domainIPv6 string) (bool, bool) {
	ip4AddressesMatch := publicIPv4 == domainIPv4

	ip6AddressesMatch := publicIPv6 == domainIPv6

	return ip4AddressesMatch, ip6AddressesMatch

}

func main() {
	// Config setup
	viper.SetConfigName("default") // config file name without extension
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Fatal error reading config file: default \n", err)
		os.Exit(1)
	}

	var ipInfo ipData
	ipInfo.publicIP4Api = viper.GetString("app.publicIP4Api")
	ipInfo.publicIP6Api = viper.GetString("app.publicIP6Api")
	ipInfo.domainName = viper.GetString("app.domainName")

	// End config setup

	// Get public IPs of host
	ipInfo.publicIPv4, ipInfo.publicIPv6 = getPublicIP(ipInfo.publicIP4Api, ipInfo.publicIP6Api)
	fmt.Println(fmt.Sprintf("IPv4: %s \nIPv6: %s", ipInfo.publicIPv4, ipInfo.publicIPv6))

	// Resolve IPs of provided domain
	ipInfo.domainIPv4, ipInfo.domainIPv6 = getDomainIP(ipInfo.domainName)

	// Check if public and resolved IPs are the same
	fmt.Println(compareIPs(ipInfo.publicIPv4, ipInfo.domainIPv4, ipInfo.publicIPv6, ipInfo.domainIPv6))
}
