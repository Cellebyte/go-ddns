package dyndns

import (
	"strings"
	"testing"
)

func TestParseProvider_Valid(t *testing.T) {
	cases := map[string]Provider{
		"cloudflare":    CloudFlare,
		"prepaidhoster": PrepaidHoster,
	}

	for input, want := range cases {
		t.Run(input, func(t *testing.T) {
			got, err := ParseProvider(input)
			if err != nil {
				t.Fatalf("unexpected error parsing %q: %v", input, err)
			}
			if got != want {
				t.Fatalf("unexpected provider: got %q want %q", got, want)
			}
			if got.String() != input {
				t.Fatalf("String() mismatch: got %q want %q", got.String(), input)
			}
		})
	}
}

func TestParseProvider_Invalid(t *testing.T) {
	input := "not-a-provider"
	got, err := ParseProvider(input)
	if err == nil {
		t.Fatalf("expected error for invalid provider %q, got nil", input)
	}
	if got.String() != "" {
		t.Fatalf("expected zero-value provider on error, got %q", got)
	}
	if !strings.Contains(err.Error(), "cannot parse") || !strings.Contains(err.Error(), input) {
		t.Fatalf("unexpected error message: %v", err)
	}
}
