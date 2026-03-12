package discovery

import (
	"net/netip"

	"github.com/cellebyte/go-ddns/internal/helpers"
)

type InterfaceNameTracker struct {
	provider      Provider
	interfaceName string
}

func NewInterfaceNameTracker(interfaceName string) (c InterfaceNameTracker, err error) {
	c.provider = InterfaceName
	c.interfaceName = interfaceName
	return c, helpers.ErrNotImplemeted
}

func (c InterfaceNameTracker) Provider() Provider {
	return c.provider
}

func (c InterfaceNameTracker) IPv4() (a netip.Addr, err error) {
	return a, helpers.ErrNotImplemeted

}

func (c InterfaceNameTracker) IPv6() (a netip.Addr, err error) {
	return a, helpers.ErrNotImplemeted
}
