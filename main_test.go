package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	const (
		ip4Address string = "192.168.1.1"
		ip6Address string = "2001:0db8:85a3:0000:0000:8a2e:0370:7334" // Sample IPv6 address found on Wikipedia
	)

	// Create a new Serv Mux so that we can attach multiple routes to the httptest server
	mux := http.NewServeMux()

	mux.HandleFunc("/ipv4", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, ip4Address)
	})

	mux.HandleFunc("/ipv6", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, ip6Address)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Use stub server to test function
	ans1 := getPublicIP(ts.URL + "/ipv4")

	if ans1 != ip4Address {
		t.Errorf("Did not receive expected IPv4")
	}

	// Use stub server to test function
	ans2 := getPublicIP(ts.URL + "/ipv6")

	if ans2 != ip6Address {
		t.Errorf("Did not receive expected IPv6")
	}

	ts.Close()

}

func TestCheckIPsMatch(t *testing.T) {
	publicIPv4, domainIPv4 := "192.168.1.1", "192.168.1.1"
	publicIPv6, domainIPv6 := "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"

	resultIPv4 := checkIPsMatch(publicIPv4, domainIPv4)
	resultIPv6 := checkIPsMatch(publicIPv6, domainIPv6)

	if resultIPv4 != true && resultIPv6 != true {
		t.Errorf("Failed to detect that public and domain IPs match")
	}

}

func TestCheckIPsMatchFalse(t *testing.T) {
	publicIPv4, domainIPv4 := "192.168.1.1", "192.168.10.10"
	publicIPv6, domainIPv6 := "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "0000:0db8:85a3:0000:0000:8a2e:0370:0000"

	resultIPv4 := checkIPsMatch(publicIPv4, domainIPv4)
	resultIPv6 := checkIPsMatch(publicIPv6, domainIPv6)

	if resultIPv4 != false && resultIPv6 != false {
		t.Errorf("Failed to detect that public and domain IPs don't match")
	}
}
