// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary dyndns updates configured DNS records with the
// current public IPv4 address (of network interface uplink0).
package main

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/cellebyte/go-ddns/internal/discovery"
	"github.com/cellebyte/go-ddns/internal/doh"
	"github.com/cellebyte/go-ddns/internal/dyndns"
	"github.com/libdns/libdns"

	"golang.org/x/net/dns/dnsmessage"
)

var update = dyndns.Update

type DynDNS struct {
	// TODO: multiple providers support
	DynDNSAPIToken string          `json:"dyndns_api_token"`
	DynDNSProvider dyndns.Provider `json:"dyndns_provider"`
	Zone           string          `json:"zone"`
	RecordName     string          `json:"record_name"`
	RecordType     string          `json:"record_type"`
	// TODO: make RecordType customizable if non-A is ever desired
	RecordTTLSeconds int `json:"record_ttl_seconds"`

	DOHProvider doh.Provider `json:"doh_provider"`
	DOHEndpoint string       `json:"doh_endpoint"`
}

func getDiscoveredIPs() {
	wtfIPClient, err := discovery.NewAddressTxtClient("https://myip.wtf/text")
	if err != nil {
		panic(err)
	}
	fiIPClient, err := discovery.NewAddressTxtClient("https://my.ip.fi")
	if err != nil {
		panic(err)
	}
	publicV4, _ := wtfIPClient.IPv4()
	publicV6, _ := wtfIPClient.IPv6()

	fiPublicV4, _ := fiIPClient.IPv4()
	fiPublicV6, _ := fiIPClient.IPv6()

	fmt.Println("myip.wtf says:", publicV4, publicV6)
	fmt.Println("my.ip.fi says:", fiPublicV4, fiPublicV6)
}

func getDNSValues(fqdn string) {
	d, err := doh.NewClient("google", "")
	if err != nil {
		panic(err)
	}
	aVal, err := d.Query(fqdn, dnsmessage.TypeA)
	aaaaVal, err := d.Query(fqdn, dnsmessage.TypeAAAA)
	txtVal, err := d.Query(fqdn, dnsmessage.TypeTXT)
	cnameVal, err := d.Query(fqdn, dnsmessage.TypeCNAME)

	fmt.Println(fqdn, aVal)
	fmt.Println(fqdn, aaaaVal)
	fmt.Println(fqdn, txtVal)
	fmt.Println(fqdn, cnameVal)
}

func setDNSValue(fqdn, ip string) {
	zone := "cellebyte.de"
	dnsProviderClient, err := dyndns.PrepaidHoster.DNSProvider("")
	if err != nil {
		panic(fmt.Errorf("getting dns provider %w", err))
	}
	parsedIP, _ := netip.ParseAddr(ip)
	record := libdns.Address{
		Name: libdns.RelativeName(fqdn, zone),
		TTL:  time.Duration(600 * time.Second),
		IP:   parsedIP,
	}
	dyndns.Update(context.TODO(), zone, record, dnsProviderClient)

}

func main() {
	getDiscoveredIPs()
	getDNSValues("www.selfnet.de")
	getDNSValues("my.ip.fi")
	getDNSValues("myip.wtf")
	getDNSValues("example.com")

	setDNSValue("my.cellebyte.de", "10.0.0.0")
}
