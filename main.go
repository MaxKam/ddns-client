package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type ipData struct {
	publicIPApi string
	domainName  string
	publicIP    string
	domainIP    string
}

func main() {
	// Config setup
	var err error

	viper.SetConfigName("ddns-client-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Fatal error reading config file: default \n", err)
		os.Exit(1)
	}

	// setup logfile if logOutput is set to logfile, otherwise by default with use Journald.
	if viper.GetString("app.logOutput") == "logfile" {
		logFile, err := os.OpenFile(viper.GetString("app.logLocation"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logFile)
	}
	// end log setup

	var ipv4Info ipData
	ipv4Info.publicIPApi = viper.GetString("app.publicIP4Api")
	ipv4Info.domainName = viper.GetString("app.domainName")

	gcpInfo := &gcpData{
		projectName: viper.GetString("gcpDNS.projectName"),
		zoneName:    viper.GetString("gcpDNS.zoneName"),
		ttlValue:    viper.GetInt64("gcpDNS.ttlValue"),
	}

	// End config setup

	log.Println("Dynamic DNS client - Starting check of public IPv4 address")

	// Get public IPv4 of host
	ipv4Info.publicIP = getPublicIP(ipv4Info.publicIPApi)
	log.Printf("Hosts public IPv4: %s", ipv4Info.publicIP)

	// Resolve IPv4 of provided domain
	ipv4Info.domainIP = getDomainIP(ipv4Info.domainName, "A")

	// Check if public and resolved IPs are the same
	IPv4Same := checkIPsMatch(ipv4Info.publicIP, ipv4Info.domainIP)

	if !IPv4Same {
		log.Println("Public IPv4 address has changed. Updating DNS record")
		UpdateDNSRecord(&ipv4Info, gcpInfo, "A")
	} else {
		log.Println("Public IPv4 address has not changed.")
	}

	// Check IPV6 if enabled
	if viper.GetBool("app.ipv6Enabled") == true {
		log.Println("Dynamic DNS client - Starting check of public IPv6")

		var ipv6Info ipData
		ipv6Info.publicIPApi = viper.GetString("app.publicIPApi")
		ipv6Info.domainName = viper.GetString("app.domainName")

		// Get public IPv4 of host
		ipv6Info.publicIP = getPublicIP(ipv6Info.publicIPApi)
		log.Printf("Hosts public IPv6: %s", ipv6Info.publicIP)

		// Resolve IPv6 of provided domain
		ipv6Info.domainIP = getDomainIP(ipv6Info.domainName, "AAAA")

		// Check if public and resolved IPs are the same
		IPv6Same := checkIPsMatch(ipv6Info.publicIP, ipv6Info.domainIP)

		if !IPv6Same {
			log.Println("Public IPv6 address has changed. Updating DNS record")
			UpdateDNSRecord(&ipv6Info, gcpInfo, "AAAA")
		} else {
			log.Println("Public IPv6 address has not changed.")
		}
	} else {
		log.Println("IPv6 is disabled. Skipping...")
	}

	log.Println("Dynamic DNS client finished run. Exiting.")

	os.Exit(0)

}
