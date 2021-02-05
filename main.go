package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		fmt.Fprintf(os.Stderr, "Could not read IPv4 response: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "Could not read IPv6 response: %v\n", err)
		os.Exit(1)
	}

	return string(ipV4), string(ipV6)
}

func main() {
	ipV4, ipV6 := getPublicIP("https://api.ipify.org", "https://api6.ipify.org")
	fmt.Println(fmt.Sprintf("IPv4: %s \nIPv6: %s", ipV4, ipV6))
}
