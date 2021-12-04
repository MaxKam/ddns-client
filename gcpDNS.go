package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/dns/v1"
)

type gcpData struct {
	projectName string
	zoneName    string
	ttlValue    int64
}

// UpdateDNSRecord will update the DNS record in a Google Cloud DNS Managed Zone.
func UpdateDNSRecord(ipInfo *ipData, gcpInfo *gcpData, ipType string) {
	var newPublicIP string
	var previousIP string

	if ipType == "A" {
		newPublicIP = ipInfo.publicIPv4
		previousIP = ipInfo.domainIPv4
	} else if ipType == "AAAA" {
		newPublicIP = ipInfo.publicIPv6
		previousIP = ipInfo.domainIPv6
	}

	ctx := context.Background()

	dnsService, err := dns.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// GCP requires that the domain name be fully qualified, i.e. includes a period at the end for the root zone
	domainName := ipInfo.domainName + "."

	addResource := &dns.ResourceRecordSet{
		Kind: "dns#resourceRecordSet",
		Name: domainName,
		Rrdatas: []string{
			newPublicIP,
		},
		Ttl:  gcpInfo.ttlValue,
		Type: ipType,
	}

	deleteResource := &dns.ResourceRecordSet{
		Kind: "dns#resourceRecordSet",
		Name: domainName,
		Rrdatas: []string{
			previousIP,
		},
		Ttl:  gcpInfo.ttlValue,
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

	resp, err := dnsService.Changes.Create(gcpInfo.projectName, gcpInfo.zoneName, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Request to update DNS record %s (%s) with IP %s sent. Status: %s", domainName, ipType, newPublicIP, resp.Status)

	time.Sleep(10 * time.Second)

	getStatus, err := dnsService.Changes.Get(gcpInfo.projectName, gcpInfo.zoneName, resp.Id).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Status of request to update DNS record %s (%s): %s", domainName, ipType, getStatus.Status)

}
