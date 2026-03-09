module github.com/cellebyte/go-ddns

go 1.25.5

require (
	github.com/libdns/libdns v1.1.1
	golang.org/x/net v0.48.0
)

require (
	github.com/cellebyte/go-pph v0.0.2 // indirect
	github.com/libdns/cloudflare v0.2.2 // indirect
	github.com/libdns/pph v0.0.0-00010101000000-000000000000 // indirect
)

replace github.com/libdns/pph => github.com/cellebyte/pph v0.0.1
