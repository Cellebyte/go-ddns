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
	"flag"
	"fmt"
	"os"

	"github.com/cellebyte/go-ddns/internal/certbot"
	"github.com/cellebyte/go-ddns/internal/config"
	"github.com/cellebyte/go-ddns/internal/ddns"
)

var acmeCmd *flag.FlagSet

var certbotCmd *flag.FlagSet
var certbotAuthHook bool
var certbotCleanupHook bool

var ddnsCmd *flag.FlagSet

func init() {
	certbotCmd = flag.NewFlagSet("certbot", flag.ExitOnError)
	certbotCmd.BoolVar(&certbotAuthHook, "auth-hook", false, "provide it to trigger the auth hook. (ref: https://eff-certbot.readthedocs.io/en/stable/using.html#pre-and-post-validation-hooks)")
	certbotCmd.BoolVar(&certbotCleanupHook, "cleanup-hook", false, "provide it to trigger the cleanup hook. (ref: https://eff-certbot.readthedocs.io/en/stable/using.html#pre-and-post-validation-hooks)")
	ddnsCmd = flag.NewFlagSet("ddns", flag.ExitOnError)
	acmeCmd = flag.NewFlagSet("acme", flag.ExitOnError)
}

func updateDns() {
	config, err := config.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("parsing config: %w", err))
	}
	switch os.Args[1] {
	case acmeCmd.Name():
		err := acmeCmd.Parse(os.Args[2:])
		if err != nil {
			panic(fmt.Errorf("parsing acme command args: %w", err))
		}
		err = certbot.ManageCertificate(config)
		if err != nil {
			panic(fmt.Errorf("manage certificate: %w", err))
		}
	case certbotCmd.Name():
		err = certbotCmd.Parse(os.Args[2:])
		if err != nil {
			panic(fmt.Errorf("parsing certbot command args: %w", err))
		}
		hook_values, err := certbot.ParseParams()
		if err != nil {
			panic(fmt.Errorf("parsing certbot manual hook env: %w", err))
		}
		if certbotAuthHook != certbotCleanupHook {
			if certbotAuthHook {
				certbot.Auth(config, hook_values)
			}
			if certbotCleanupHook {
				certbot.Cleanup(config, hook_values)
			}
		} else {
			panic(fmt.Errorf("provide either auth=%v or cleanup=%v", certbotAuthHook, certbotCleanupHook))
		}
		fmt.Printf("certbot: hook auth=%v, cleanup=%v\n", certbotAuthHook, certbotCleanupHook)
	case ddnsCmd.Name():
		ddns.UpdateDNS(config)
	default:
		fmt.Println("expected 'certbot' or 'ddns' or 'acme' subcommands")
		os.Exit(1)
	}
}

func main() {
	updateDns()
}
