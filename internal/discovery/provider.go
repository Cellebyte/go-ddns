package discovery

import (
	"fmt"
	"net/netip"

	"github.com/cellebyte/go-ddns/internal/helpers"
)

type Provider int

const (
	unknown Provider = iota
	InterfaceName
	AddressTxt
	FritzBox
)

func (p Provider) String() string {
	switch p {
	case InterfaceName:
		return "interfaceName"
	case AddressTxt:
		return "addressTxt"
	case FritzBox:
		return "fritz.box"
	}
	return ""
}

func (p Provider) New(val string) (DisoveryProvider, error) {
	switch p {
	case InterfaceName:
		client, err := NewInterfaceNameTracker(val)
		return &client, err
	case AddressTxt:
		client, err := NewAddressTxtClient(val)
		return &client, err
	case FritzBox:
		client, err := NewFritzBoxClient(val)
		return &client, err
	}
	return nil, helpers.ErrNotImplemeted
}

func ParseProvider(provider string) (Provider, error) {
	switch provider {
	case InterfaceName.String():
		return InterfaceName, nil
	case AddressTxt.String():
		return AddressTxt, nil
	case FritzBox.String():
		return FritzBox, nil
	}
	return unknown, fmt.Errorf("cannot parse %q", provider)
}

type DisoveryProvider interface {
	Provider() Provider
	IPv4() (netip.Addr, error)
	IPv6() (netip.Addr, error)
}
