package dyndns

import (
	"fmt"

	"github.com/cellebyte/go-ddns/internal/helpers"
	"github.com/libdns/cloudflare"
	hetzner "github.com/libdns/hetzner/v2"
	"github.com/libdns/pph"
)

type Provider int

const (
	unknown Provider = iota
	CloudFlare
	PrepaidHoster
	Hetzner
)

func (p Provider) String() string {
	switch p {
	case CloudFlare:
		return "cloudflare"
	case PrepaidHoster:
		return "pph"
	case Hetzner:
		return "hetzner"
	}
	return ""
}

func ParseProvider(provider string) (Provider, error) {
	switch provider {
	case CloudFlare.String():
		return CloudFlare, nil
	case PrepaidHoster.String():
		return PrepaidHoster, nil
	case Hetzner.String():
		return Hetzner, nil
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
	case Hetzner:
		return &hetzner.Provider{
			APIToken: token,
		}, nil
	default:
		return nil, helpers.ErrNotImplemeted
	}
}
