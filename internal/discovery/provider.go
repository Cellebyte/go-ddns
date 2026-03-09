package discovery

import (
	"fmt"
	"net/netip"
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
	Type() Provider
	IPv4() (netip.Addr, error)
	IPv6() (netip.Addr, error)
}
