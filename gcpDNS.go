package main

import (
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
)

type gcpData struct {
	projectName string
	zoneName    string
	ttlValue    int64
}

// UpdateDNSRecord will update the DNS record in a Google Cloud DNS Managed Zone.
func UpdateDNSRecord(projectName string, zoneName string, domainName string, previousIP string, newPublicIP string, ttlValue int64) {
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, dns.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	dnsService, err := dns.New(c)
	if err != nil {
		log.Fatal(err)
	}

	ipType := getIPType(newPublicIP)
	// GCP requires that the domain name be fully qualified, i.e. includes a period at the end for the root zone
	domainName = domainName + "."

	addResource := &dns.ResourceRecordSet{
		Kind: "dns#resourceRecordSet",
		Name: domainName,
		Rrdatas: []string{
			newPublicIP,
		},
		Ttl:  ttlValue,
		Type: ipType,
	}

	deleteResource := &dns.ResourceRecordSet{
		Kind: "dns#resourceRecordSet",
		Name: domainName,
		Rrdatas: []string{
			previousIP,
		},
		Ttl:  ttlValue,
		Type: ipType,
	}

	rb := &dns.Change{
		Additions: []*dns.ResourceRecordSet{
			addResource,
		},
		Deletions: []*dns.ResourceRecordSet{
			deleteResource,
		},
		IsServing: true,
		Kind:      "dns#change",
	}

	resp, err := dnsService.Changes.Create(projectName, zoneName, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	log.Printf("Request to update DNS record %s (%s) with IP %s sent. Status: %s", domainName, ipType, newPublicIP, resp.Status)

	time.Sleep(10 * time.Second)

	getStatus, err := dnsService.Changes.Get(projectName, zoneName, resp.Id).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Status of request to update DNS record %s (%s): %s", domainName, ipType, getStatus.Status)

}

func getIPType(inputIP string) string {
	ipType := ""
	if strings.Count(inputIP, ":") < 2 {
		ipType = "A"
	} else if strings.Count(inputIP, ":") >= 2 {
		ipType = "AAAA"
	}
	return ipType
}
