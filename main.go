package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getPublicIp() (string, string) {
	// Get public IPv4 address of host
	reqV4, _ := http.Get("https://api.ipify.org")
	ipV4, _ := ioutil.ReadAll(reqV4.Body)

	// Get public IPv6 address of host
	reqV6, _ := http.Get("https://api6.ipify.org")
	ipV6, _ := ioutil.ReadAll(reqV6.Body)

	return string(ipV4), string(ipV6)
}

func main() {
	ipV4, ipV6 := getPublicIp()
	fmt.Println(fmt.Sprintf("IPv4: %s \nIPv6: %s", ipV4, ipV6))

}
