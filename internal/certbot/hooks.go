package certbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cellebyte/go-ddns/internal/config"
	"github.com/cellebyte/go-ddns/internal/doh"
	"github.com/libdns/libdns"
	"github.com/sethvargo/go-retry"
	"golang.org/x/net/dns/dnsmessage"
)

func GetChallengeRecordName(identifier, zone string) string {
	return libdns.RelativeName(fmt.Sprintf("%s.%s", ACMETXTPrefix, identifier), zone)
}

func Hook(config config.DynDNS, identifier, challenge string, cleanup bool) error {
	ctx := context.Background()
	dohClient, err := doh.NewClient(config.DOHProvider.String(), config.DOHProvider.Endpoint())
	if err != nil {
		return fmt.Errorf("creating doh client for %v+: %w", config.DOHProvider, err)
	}
	dnsProviderClient, err := config.Provider.New(config.APIToken)
	if err != nil {
		return fmt.Errorf("instantiating client: %w", err)
	}
	if cleanup {
		records, err := dnsProviderClient.GetRecords(ctx, config.Zone)
		if err != nil {
			return fmt.Errorf("getting Records: %w", err)
		}
		name := GetChallengeRecordName(identifier, config.Zone)
		var toDelete []libdns.Record
		for _, record := range records {
			if record.RR().Name == name {
				toDelete = append(toDelete, record)
			}
		}
		deleted, err := dnsProviderClient.DeleteRecords(ctx, config.Zone, toDelete)
		if err != nil {
			return fmt.Errorf("deleting [%v] only deleted [%v]: %w", toDelete, deleted, err)
		}
		return nil
	}
	// ACME DNS-01 challenge
	txtRecord := libdns.TXT{
		Name: GetChallengeRecordName(identifier, config.Zone),
		TTL:  time.Duration(600 * time.Second),
		Text: challenge,
	}
	updated, err := dnsProviderClient.SetRecords(ctx, config.Zone, []libdns.Record{txtRecord})
	if err != nil {
		return fmt.Errorf("creating %v only created [%v]: %w", txtRecord, updated, err)
	}
	// Check DNS for existence
	fullName := libdns.AbsoluteName(txtRecord.Name, config.Zone)

	if err := retry.Do(ctx, retry.WithMaxRetries(
		4, retry.WithJitter(
			1*time.Second, retry.NewFibonacci(
				2*time.Second),
		),
	), func(ctx context.Context) error {
		responses, err := dohClient.Query(fullName, dnsmessage.TypeTXT)
		if err != nil {
			return fmt.Errorf("resolving Challenge for %s: %w", fullName, err)
		}
		for _, response := range responses {
			if txtRecord.Text == strings.Trim(response, " ") {
				return nil
			}
		}
		err = fmt.Errorf("finding token for %s {%s != %v} on endpoint %s", fullName, txtRecord.Text, responses, config.DOHProvider.Endpoint())
		fmt.Println(ctx, err)
		return retry.RetryableError(err)
	}); err != nil {
		return fmt.Errorf("retry for existence: %w", err)
	}
	return nil
}

func Auth(config config.DynDNS, params CertBotParameters) {
	auth(config, params.Identifier, params.Validation)
}

func Cleanup(config config.DynDNS, params CertBotParameters) {
	cleanup(config, params.Identifier, params.Validation)
}

func auth(config config.DynDNS, identifier, challenge string) {
	err := Hook(config, identifier, challenge, false)
	if err != nil {
		panic(fmt.Errorf("calling auth hook: %w", err))
	}
}
func cleanup(config config.DynDNS, identifier, challenge string) {
	err := Hook(config, identifier, challenge, true)
	if err != nil {
		panic(fmt.Errorf("calling cleanup hook: %w", err))
	}
}
