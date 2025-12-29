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
	"fmt"

	"github.com/cellebyte/go-ddns/internal/doh"
	"github.com/cellebyte/go-ddns/internal/dyndns"

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

	DohURL string `json:"doh_endpoint,omitempty"`
}

func main() {
	d, err := doh.NewClient("google", "")
	if err != nil {
		panic(err)
	}
	urlString := "www.selfnet.de"
	aVal, err := d.Query(urlString, dnsmessage.TypeA)
	aaaaVal, err := d.Query(urlString, dnsmessage.TypeAAAA)
	txtVal, err := d.Query(urlString, dnsmessage.TypeTXT)
	cnameVal, err := d.Query(urlString, dnsmessage.TypeCNAME)

	fmt.Println(urlString, aVal)
	fmt.Println(urlString, aaaaVal)
	fmt.Println(urlString, txtVal)
	fmt.Println(urlString, cnameVal)
	return

}
