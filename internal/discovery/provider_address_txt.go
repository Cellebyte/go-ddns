package discovery

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

type AddressTxtClient struct {
	endpoint     *url.URL
	v4HttpClient *http.Client
	v6HttpClient *http.Client
}

func NewAddressTxtClient(endpoint string) (c AddressTxtClient, err error) {
	if endpoint == "" {
		return c, errors.New("endpoint missing")
	}
	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return c, fmt.Errorf("parse doh endpoint url %q: %w", endpoint, err)
	}
	c.endpoint = parsedEndpoint
	c.v4HttpClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: 60 * time.Second,
				}).DialContext(ctx, "tcp4", addr)
			},
			TLSHandshakeTimeout: 3 * time.Second,
			DisableKeepAlives:   false,
			MaxIdleConns:        256,
			MaxIdleConnsPerHost: 256,
		},
	}
	c.v6HttpClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: 60 * time.Second,
				}).DialContext(ctx, "tcp6", addr)
			},
			TLSHandshakeTimeout: 3 * time.Second,
			DisableKeepAlives:   false,
			MaxIdleConns:        256,
			MaxIdleConnsPerHost: 256,
		},
	}
	return c, err
}

func (c AddressTxtClient) parseIP(rawIP string) (netip.Addr, error) {
	rawIP = strings.TrimRight(rawIP, "\n\r")
	ip, err := netip.ParseAddr(rawIP)
	if err != nil {
		return ip, fmt.Errorf("not an IP %q: %w", rawIP, err)
	}
	return ip, nil
}

func (c AddressTxtClient) GetIPv4() (ip netip.Addr, err error) {
	r, err := c.v4HttpClient.Get(c.endpoint.String())
	if err != nil {
		return ip, fmt.Errorf("getting endpoint %s: %w", c.endpoint.String(), err)
	}
	defer r.Body.Close()
	byteString, err := io.ReadAll(r.Body)
	if err != nil {
		return ip, fmt.Errorf("reading content length %d: %w", r.ContentLength, err)
	}
	ip, err = c.parseIP(string(byteString))
	if err != nil {
		return ip, fmt.Errorf("parsing ip: %w", err)
	}
	return ip, nil
}

func (c AddressTxtClient) GetIPv6() (ip netip.Addr, err error) {
	r, err := c.v6HttpClient.Get(c.endpoint.String())
	if err != nil {
		return ip, fmt.Errorf("getting endpoint %s: %w", c.endpoint.String(), err)
	}
	defer r.Body.Close()
	byteString, err := io.ReadAll(r.Body)
	if err != nil {
		return ip, fmt.Errorf("reading content length %d: %w", r.ContentLength, err)
	}
	ip, err = c.parseIP(string(byteString))
	if err != nil {
		return ip, fmt.Errorf("parsing ip: %w", err)
	}
	return ip, nil
}
