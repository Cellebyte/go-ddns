package dyndns

import (
	"fmt"
)

type Provider string

const (
	CloudFlare    Provider = "cloudflare"
	PrepaidHoster Provider = "prepaidhoster"
	FritzBox      Provider = "fritzbox"
)

func (p Provider) String() string {
	return string(p)
}

func ParseProvider(provider string) (p Provider, err error) {
	providers := map[Provider]struct{}{
		CloudFlare:    {},
		PrepaidHoster: {},
	}

	prov := Provider(provider)
	_, ok := providers[prov]
	if !ok {
		return p, fmt.Errorf(`cannot parse:[%s] as provider`, provider)
	}
	return prov, nil
}
