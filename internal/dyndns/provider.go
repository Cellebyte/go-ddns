package dyndns

import (
	"fmt"

	"github.com/cellebyte/go-ddns/internal/helpers"
	"github.com/libdns/cloudflare"
	pph "github.com/libdns/pph"
)

type Provider int

const (
	unknown Provider = iota
	CloudFlare
	PrepaidHoster
)

func (p Provider) String() string {
	switch p {
	case CloudFlare:
		return "cloudflare"
	case PrepaidHoster:
		return "prepaidhoster"
	}
	return ""
}

func ParseProvider(provider string) (Provider, error) {
	switch provider {
	case CloudFlare.String():
		return CloudFlare, nil
	case PrepaidHoster.String():
		return PrepaidHoster, nil
	}
	return unknown, fmt.Errorf("cannot parse %q", provider)
}

func (p Provider) New(token string) (RecordGetterSetter, error) {
	switch p {
	case PrepaidHoster:
		return &pph.Provider{
			APIToken: token,
		}, nil
	case CloudFlare:
		return &cloudflare.Provider{
			APIToken: token,
		}, nil
	default:
		return nil, helpers.ErrNotImplemeted
	}
}
