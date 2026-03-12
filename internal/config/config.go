package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cellebyte/go-ddns/internal/discovery"
	"github.com/cellebyte/go-ddns/internal/doh"
	"github.com/cellebyte/go-ddns/internal/dyndns"
	"github.com/libdns/libdns"
)

const (
	// DYNDNS Provider mix config
	tokenEnv                     = "DYNDNS_API_TOKEN"
	providerEnv                  = "DYNDNS_PROVIDER"
	dohProviderEnv               = "DYNDNS_DOH_PROVIDER"
	dohProviderEndpointEnv       = "DYNDNS_DOH_PROVIDER_ENDPOINT"
	discoveryProviderEnv         = "DYNDNS_DISCOVERY_PROVIDER"
	discoveryProviderEndpointEnv = "DYNDNS_DISCOVERY_PROVIDER_ENDPOINT"

	// DNS config
	zoneEnv                 = "DYNDNS_DNS_ZONE"
	subDomainEnv            = "DYNDNS_DNS_SUBDOMAIN"
	subDomainTTLEnv         = "DYNDNS_DNS_SUBDOMAIN_TTL"
	subDomainRecordTypesEnv = "DYNDNS_DNS_SUBDOMAIN_RECORDTYPES"

	// Defaults
	defaultDohProvider               = "ffmuc.net"
	defaultDiscoveryProvider         = "addressTxt"
	defaultDiscoveryProviderEndpoint = "https://myip.wtf/text"
	defaulRecordTTL                  = 300 * time.Second

	// RecordTypes
	AType    = "A"
	AAAAType = "AAAA"
)

var (
	defaultRecordTypes = []string{AType, AAAAType}
)

type DynDNS struct {
	APIToken                  string
	Provider                  dyndns.Provider
	DOHProvider               doh.Provider
	DOHProviderEndpoint       string
	DiscoveryProvider         discovery.Provider
	DiscoveryProviderEndpoint string

	Zone        string
	RecordName  string
	RecordTypes []string
	RecordTTL   time.Duration
}

func ParseConfig() (config DynDNS, err error) {
	// Token
	config.APIToken = os.Getenv(tokenEnv)
	if config.APIToken == "" {
		return config, fmt.Errorf("missing %q", tokenEnv)
	}

	// dyndns.Provider
	provider := os.Getenv(providerEnv)
	if provider == "" {
		return config, fmt.Errorf("missing %q", providerEnv)
	}
	config.Provider, err = dyndns.ParseProvider(provider)
	if err != nil {
		return config, fmt.Errorf("parsing provider %q: %w", provider, err)
	}

	// doh.Provider
	dohProvider := os.Getenv(dohProviderEnv)
	if dohProvider == "" {
		dohProvider = defaultDohProvider
	}
	config.DOHProvider, err = doh.ParseProvider(dohProvider)
	if err != nil {
		return config, fmt.Errorf("parsing doh provider %q: %w", dohProvider, err)
	}
	// doh.Provider endpoint
	// this value is ignored if a provider is defined
	config.DOHProviderEndpoint = os.Getenv(dohProviderEndpointEnv)

	// discovery.Provider
	discoveryProvider := os.Getenv(discoveryProviderEnv)
	if discoveryProvider == "" {
		discoveryProvider = defaultDiscoveryProvider
	}
	config.DiscoveryProvider, err = discovery.ParseProvider(discoveryProvider)
	if err != nil {
		return config, fmt.Errorf("parsing discovery provider %q: %w", discoveryProvider, err)
	}
	// discover.Provider endpoint
	config.DiscoveryProviderEndpoint = os.Getenv(discoveryProviderEndpointEnv)
	if config.DiscoveryProviderEndpoint == "" {
		config.DiscoveryProviderEndpoint = defaultDiscoveryProviderEndpoint
	}
	// zone
	config.Zone = os.Getenv(zoneEnv)
	if config.Zone == "" {
		return config, fmt.Errorf("missing %q", zoneEnv)
	}
	if !ValidRFC1035Name(config.Zone) {
		return config, fmt.Errorf("invalid value %s=%q", zoneEnv, config.Zone)
	}

	// subdomain
	config.RecordName = os.Getenv(subDomainEnv)
	if config.RecordName == "" {
		return config, fmt.Errorf("missing %q", subDomainEnv)
	}
	if !ValidRFC1035Name(config.RecordName) {
		return config, fmt.Errorf("invalid value %s=%q", subDomainEnv, config.RecordName)
	}
	config.RecordName = libdns.RelativeName(config.RecordName, config.Zone)

	// ttl
	config.RecordTTL = defaulRecordTTL
	ttlString := os.Getenv(subDomainTTLEnv)
	if ttlString != "" {
		ttl, err := strconv.ParseInt(ttlString, 10, 0)
		if err != nil {
			return config, fmt.Errorf("parsing %q, %w", ttlString, err)
		}
		config.RecordTTL = time.Duration(ttl) * time.Second
	}

	// recordTypes
	config.RecordTypes = defaultRecordTypes
	recordTypesString := strings.ToUpper(os.Getenv(subDomainRecordTypesEnv))
	if recordTypesString != "" {
		recordTypes := UniqueSliceElements(strings.Split(recordTypesString, ","))
		if !ValidRecordTypes(recordTypes) {
			return config, fmt.Errorf("invalid value %s=%q", subDomainRecordTypesEnv, recordTypesString)
		}
		config.RecordTypes = recordTypes
	}

	return config, nil
}
