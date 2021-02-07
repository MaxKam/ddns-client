package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

	publicIP4Url := viper.GetString("app.ip4Url")
	publicIP6Url := viper.GetString("app.ip6Url")

	// End config setup

	ipV4, ipV6 := getPublicIP(publicIP4Url, publicIP6Url)
	fmt.Println(fmt.Sprintf("IPv4: %s \nIPv6: %s", ipV4, ipV6))
}
