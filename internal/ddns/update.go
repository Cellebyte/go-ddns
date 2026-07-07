package ddns

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/cellebyte/go-ddns/internal/config"
	"github.com/cellebyte/go-ddns/internal/discovery"
	"github.com/cellebyte/go-ddns/internal/doh"
	"github.com/cellebyte/go-ddns/internal/dyndns"
	"github.com/cellebyte/go-ddns/internal/helpers"
	"github.com/libdns/libdns"
	"golang.org/x/net/dns/dnsmessage"
)

func discoverIP(d discovery.DisoveryProvider, recordType string) (netip.Addr, error) {
	switch recordType {
	case config.AType:
		return d.IPv4()
	case config.AAAAType:
		return d.IPv6()
	}
	return netip.Addr{}, helpers.ErrNotImplemeted
}

func getDNSValue(d doh.Client, dnsName, recordType string) (netip.Addr, error) {
	var messageType dnsmessage.Type
	switch recordType {
	case config.AType:
		messageType = dnsmessage.TypeA
	case config.AAAAType:
		messageType = dnsmessage.TypeAAAA
	default:
		return netip.Addr{}, helpers.ErrNotImplemeted
	}
	val, err := d.Query(dnsName, messageType)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("querying for %q: %w", dnsName, err)
	}
	var addr netip.Addr
	if len(val) > 0 {
		addr, err = netip.ParseAddr(val[0])
		if err != nil {
			return addr, fmt.Errorf("parsing %s: %w", val[0], err)
		}
	}
	return addr, nil
}

func update(config config.DynDNS) error {
	dnsName := libdns.AbsoluteName(config.RecordName, config.Zone)
	dohClient, err := doh.NewClient(config.DOHProvider.String(), config.DOHProvider.Endpoint())
	if err != nil {
		return fmt.Errorf("creating doh client for %v+: %w", config.DOHProvider, err)
	}
	discoveryClient, err := config.DiscoveryProvider.New(config.DiscoveryProviderEndpoint)
	if err != nil {
		return fmt.Errorf("creating discovery client for %v+: %w", config.DiscoveryProvider, err)
	}
	dnsProviderClient, err := config.Provider.New(config.APIToken)
	if err != nil {
		return fmt.Errorf("creating dns provider client for %v+: %w", config.Provider, err)
	}
	for _, recordType := range config.RecordTypes {
		oldAddr, err := getDNSValue(dohClient, dnsName, recordType)
		if err != nil {
			return fmt.Errorf("getting dns content for [%s] %s: %w", recordType, dnsName, err)
		}
		newAddr, err := discoverIP(discoveryClient, recordType)
		if err != nil {
			return fmt.Errorf("getting ip for [%s] %s: %w", recordType, dnsName, err)
		}
		if newAddr == oldAddr {
			fmt.Printf("%s: %s already has address in DNS (old %q == new %q)\n", recordType, dnsName, oldAddr.String(), newAddr.String())
		}
		record := libdns.Address{
			Name: libdns.RelativeName(dnsName, config.Zone),
			TTL:  time.Duration(600 * time.Second),
			IP:   newAddr,
		}
		err = dyndns.Update(context.TODO(), config.Zone, record, dnsProviderClient)
		if err != nil {
			return fmt.Errorf("updating ip for [%s] %s to %q: %w", recordType, dnsName, newAddr.String(), err)
		}
		fmt.Printf("%s: %s now has address %q\n", recordType, dnsName, newAddr.String())
	}
	return nil
}

func UpdateDNS(config config.DynDNS) {
	fmt.Println(time.Now().Format("15:04:05"), ":: Found config:", config)
	ticker := time.NewTicker(config.RecordTTL / 3)
	defer ticker.Stop()
	fmt.Println(time.Now().Format("15:04:05"), ":: Update record", libdns.AbsoluteName(config.RecordName, config.Zone), "every", config.RecordTTL/3)
	err := update(config)
	if err != nil {
		panic(fmt.Errorf("running update: %w", err))
	}
	for {
		t := <-ticker.C
		// This code runs every third part of config.RecordTTL
		fmt.Println(t.Format("15:04:05"), ":: Update record", libdns.AbsoluteName(config.RecordName, config.Zone))
		err := update(config)
		if err != nil {
			panic(fmt.Errorf("%s :: running update: %w", t.Format("15:04:05"), err))
		}
	}
}
