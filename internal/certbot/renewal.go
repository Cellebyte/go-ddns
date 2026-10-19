package certbot

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/cellebyte/go-ddns/internal/config"
	"github.com/libdns/libdns"
	"golang.org/x/crypto/acme/autocert"
)

// AccountKeyNeedsRotation returns true when the account key is missing, cannot be parsed,
// or is smaller than 2048 bits.
func AccountKeyNeedsRotation(ctx context.Context, cache autocert.Cache) (bool, error) {
	accountKeyBytes, err := cache.Get(ctx, defaultAccountKeyName)
	if err != nil {
		// autocert.ErrCacheMiss means there's no key stored yet -> rotate/create
		if errors.Is(err, autocert.ErrCacheMiss) {
			return true, nil
		}
		return false, fmt.Errorf("error reading account key from cache: %w", err)
	}

	// Try parse as PKCS1 (the form used elsewhere in this repo)
	accountKey, err := x509.ParsePKCS1PrivateKey(accountKeyBytes)
	if err != nil {
		// Can't parse -> rotate
		return true, nil
	}

	// If key is smaller than 2048 bits, request rotation
	if accountKey.N.BitLen() < 2048 {
		return true, nil
	}

	return false, nil
}

// CertNeedsRenewal checks if a cached certificate for domain needs renewal.
// It returns (needsRenewal, parsedCert, error).
// If there is no cert in cache, needsRenewal==true and parsedCert==nil.
func CertNeedsRenewal(ctx context.Context, cache autocert.Cache, domain string, renewBefore time.Duration) (bool, *x509.Certificate, time.Duration, error) {
	data, err := cache.Get(ctx, getFileName(domain, "crt.pem"))
	if err != nil {
		if errors.Is(err, autocert.ErrCacheMiss) {
			// No cert -> needs issuance
			return true, nil, 0, nil
		}
		return false, nil, 0, fmt.Errorf("error reading cert from cache: %w", err)
	}

	// Try to parse certificate. Accept PEM or DER.
	var cert *x509.Certificate
	if block, _ := pem.Decode(data); block != nil {
		parsed, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			// Parsing failed -> treat as needing renewal
			return true, nil, 0, nil
		}
		cert = parsed
	} else {
		parsed, err := x509.ParseCertificate(data)
		if err != nil {
			// Parsing failed -> treat as needing renewal
			return true, nil, 0, nil
		}
		cert = parsed
	}
	timeLeft := time.Until(cert.NotAfter)

	// If cert is already expired -> renew
	if timeLeft <= 0 {
		return true, cert, timeLeft, nil
	}

	// If cert is within the renewBefore window -> renew
	if timeLeft <= renewBefore {
		return true, cert, timeLeft, nil
	}

	// Otherwise, no renewal needed
	return false, cert, timeLeft, nil
}

// ManageCertificate is the top-level helper you call to ensure the certificate (and account key)
// are present and fresh. It will run CertChallenge only when required.
//
// This function is intended to be invoked by an external scheduler (cron) or orchestration system.
func ManageCertificate(dyndns config.DynDNS) error {
	cfg, err := ParseACMEConfig()
	if err != nil {
		return fmt.Errorf("failed building config: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	cache := autocert.DirCache(cfg.CacheDir)

	// 1) Check account key
	keyRot, err := AccountKeyNeedsRotation(ctx, cache)
	if err != nil {
		return fmt.Errorf("account key check failed: %w", err)
	}

	// 2) Check cert
	domain := libdns.AbsoluteName(dyndns.RecordName, dyndns.Zone)
	certRot, _, validFor, err := CertNeedsRenewal(ctx, cache, domain, cfg.RenewBefore)
	if err != nil {
		return fmt.Errorf("certificate check failed: %w", err)
	}

	// If nothing needs rotation, return early
	if !keyRot && !certRot {
		// Nothing to do
		fmt.Println(domain, "is valid for", validFor.Round(time.Second), ": not due for renewal")
		return nil
	}

	fmt.Println(domain, "is valid for", validFor.Round(time.Second), ": starting renewal")
	// Otherwise, run the issuance/renewal flow
	if err := CertChallenge(dyndns); err != nil {
		return fmt.Errorf("cert issuance/renewal failed: %w", err)
	}

	return nil
}
