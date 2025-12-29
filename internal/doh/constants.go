package doh

import (
	"errors"
	"fmt"
)

type Provider int

const (
	custom Provider = iota
	CloudFlare
	Google
	Quad9
	WikiMedia
	Joindns4EU
)

var (
	endpoints = map[Provider]string{
		CloudFlare: "https://cloudflare-dns.com/dns-query",
		Google:     "https://dns.google/dns-query",
		Quad9:      "https://dns.quad9.net/dns-query",
		WikiMedia:  "https://wikimedia-dns.org/dns-query",
		Joindns4EU: "https://unfiltered.joindns4.eu/dns-query",
	}
)

func (p Provider) String() string {
	switch p {
	case custom:
		return "custom"
	case CloudFlare:
		return "cloudflare"
	case Google:
		return "google"
	case Quad9:
		return "quad9"
	case WikiMedia:
		return "wikimedia"
	case Joindns4EU:
		return "joindns4eu"
	}
	return fmt.Sprintf("Provider(%q)", int(p))
}

func ParseProvider(in string) (Provider, error) {
	switch in {
	case custom.String():
		return custom, nil
	case CloudFlare.String():
		return CloudFlare, nil
	case Google.String():
		return Google, nil
	case Quad9.String():
		return Quad9, nil
	case WikiMedia.String():
		return WikiMedia, nil
	case Joindns4EU.String():
		return Joindns4EU, nil
	}
	return custom, fmt.Errorf("%q is not a valid provider: %w", in, errors.New("invalid provider"))
}

func (p Provider) Endpoint() string {
	endpoint, ok := endpoints[p]
	if !ok {
		return ""
	}
	return endpoint
}
