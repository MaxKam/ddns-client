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
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Fatal error reading config file: default \n", err)
		os.Exit(1)
	}

	publicIP4Url := viper.GetString("app.publicIP4Url")
	publicIP6Url := viper.GetString("app.publicIP6Url")
	domainName := viper.GetString("app.domainName")

	// End config setup

	ipV4, ipV6 := getPublicIP(publicIP4Url, publicIP6Url)
	fmt.Println(fmt.Sprintf("IPv4: %s \nIPv6: %s", ipV4, ipV6))
	domainIPv4, domainIPv6 := getDomainIP(domainName)
	fmt.Println(compareIPs(ipV4, domainIPv4, ipV6, domainIPv6))
}
