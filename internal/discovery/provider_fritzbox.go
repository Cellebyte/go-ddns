package discovery

import (
	"net/netip"

	"github.com/cellebyte/go-ddns/internal/helpers"
)

type FritzBoxClient struct {
	provider Provider
	endpoint string
}

func NewFritzBoxClient(endpoint string) (c FritzBoxClient, err error) {
	c.provider = FritzBox
	c.endpoint = endpoint
	return c, helpers.ErrNotImplemeted
}

func (c FritzBoxClient) Provider() Provider {
	return c.provider
}

func (c FritzBoxClient) IPv4() (a netip.Addr, err error) {
	return a, helpers.ErrNotImplemeted

}

func (c FritzBoxClient) IPv6() (a netip.Addr, err error) {
	return a, helpers.ErrNotImplemeted
}
