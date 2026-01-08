package doh

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type Client struct {
	endpoint   *url.URL
	httpClient *http.Client
}

func NewClient(provider, dohEndpoint string) (Client, error) {
	var url *url.URL
	var client = Client{}
	dohProvider, err := ParseProvider(provider)
	if err != nil {
		return client, fmt.Errorf("finding provider %q: %w", provider, err)
	}
	if dohProvider != custom {
		dohEndpoint = dohProvider.Endpoint()
	}
	endpoint, err := url.Parse(dohEndpoint)
	if err != nil {
		return client, fmt.Errorf("parse doh endpoint url %q: %w", dohEndpoint, err)
	}
	client.endpoint = endpoint
	client.httpClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 3 * time.Second,
			DisableKeepAlives:   false,
			MaxIdleConns:        256,
			MaxIdleConnsPerHost: 256,
		},
	}
	return client, nil
}

func (d Client) parse(dnsRawMessage []byte, queryType dnsmessage.Type) ([]string, error) {
	var parser dnsmessage.Parser
	if _, err := parser.Start(dnsRawMessage); err != nil {
		return nil, fmt.Errorf("starting parser on %v: %w", dnsRawMessage, err)
	}
	if err := parser.SkipAllQuestions(); err != nil {
		return nil, fmt.Errorf("skipping question section: %w", err)
	}
	var values []string
	for {
		h, err := parser.AnswerHeader()
		if err == dnsmessage.ErrSectionDone {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("getting AnswerHeader: %w", err)
		}
		if h.Type != queryType || h.Class != dnsmessage.ClassINET {
			// skip everything we are not querying for
			fmt.Println("skipping", h)
			err := parser.SkipAnswer()
			if err != nil {
				return nil, fmt.Errorf("skipping answer %v: %w", h, err)
			}
			continue
		}
		switch h.Type {
		case dnsmessage.TypeA:
			ar, err := parser.AResource()
			if err != nil {
				return nil, fmt.Errorf("parsing A record: %w", err)
			}
			ip := netip.AddrFrom4(ar.A)
			values = append(values, ip.String())
		case dnsmessage.TypeAAAA:
			ar, err := parser.AAAAResource()
			if err != nil {
				return nil, fmt.Errorf("parsing AAAA record: %w", err)
			}
			ip := netip.AddrFrom16(ar.AAAA)
			values = append(values, ip.String())
		case dnsmessage.TypeCNAME:
			cn, err := parser.CNAMEResource()
			if err != nil {
				return nil, fmt.Errorf("parsing CNAME record: %w", err)
			}
			val := cn.CNAME.Data
			values = append(values, string(val[:]))
		case dnsmessage.TypeTXT:
			txt, err := parser.TXTResource()
			if err != nil {
				return nil, fmt.Errorf("parsing TXT record: %w", err)
			}
			val := txt.TXT
			values = append(values, strings.Join(val, " "))
		}
	}
	if len(values) == 0 {
		return nil, errors.New("not implemented")

	}
	return values, nil
}

func (d Client) get(dnsMessage string) ([]byte, error) {
	// Using RFC 8484
	// ref: https://datatracker.ietf.org/doc/html/rfc8484
	q := d.endpoint.Query()
	q.Set("dns", dnsMessage)
	d.endpoint.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", d.endpoint.String(), nil)
	// Using RFC 8484
	// ref: https://datatracker.ietf.org/doc/html/rfc8484
	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("doing request %q: %w", resp.Request.URL.String(), err)
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("reading body to buffer: %w", err)
	}
	return buf, nil

}

func (d Client) Query(name string, queryType dnsmessage.Type) ([]string, error) {
	// Ensure we always have an absolute name
	if !strings.HasSuffix(name, ".") {
		name = fmt.Sprintf("%s.", name)
	}
	validName, err := dnsmessage.NewName(name)
	if err != nil {
		return nil, fmt.Errorf("validating name %q: %w", name, err)
	}
	var dnsBuf [1024]byte
	b := dnsmessage.NewBuilder(dnsBuf[:0], dnsmessage.Header{RecursionDesired: true})
	b.EnableCompression()
	b.StartQuestions()
	b.Question(dnsmessage.Question{Name: validName, Type: queryType, Class: dnsmessage.ClassINET})
	binaryMessage, err := b.Finish()
	if err != nil {
		return nil, fmt.Errorf("constructing binary blob %v %s: %w", queryType, name, err)
	}
	b64Message := base64.RawURLEncoding.EncodeToString(binaryMessage)
	dnsRawMessage, err := d.get(b64Message)
	values, err := d.parse(dnsRawMessage, queryType)

	return values, nil
}
