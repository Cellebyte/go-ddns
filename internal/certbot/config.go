package certbot

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	// ref: https://eff-certbot.readthedocs.io/en/stable/using.html#pre-and-post-validation-hooks
	certBotIdentifier          = "CERTBOT_IDENTIFIER"           // The domain or IP address being authenticated
	certBotValidation          = "CERTBOT_VALIDATION"           // The validation string
	certBotToken               = "CERTBOT_TOKEN"                // Resource name part of the HTTP-01 challenge (HTTP-01 only)
	certBotRemainingChallenges = "CERTBOT_REMAINING_CHALLENGES" // Number of challenges remaining after the current challenge
	certBotAllIdentifiers      = "CERTBOT_ALL_IDENTIFIERS"      // A comma-separated list of all identifiers challenged for the current certificate

	acmeMailAddresses    = "ACME_MAIL_ADDRESSES"    // A comms-separated list of all mail addresses used on the account
	acmeChallengeTimeout = "ACME_CHALLENGE_TIMEOUT" // Parsed by the time library :default: 1m
	acmeRenewBefore      = "ACME_RENEW_BEFORE"
	// ACME configuration
	ACMETXTPrefix         = "_acme-challenge"
	defaultAccountKeyName = "account.key"
	defaultACMECacheDir   = ".cache/go-ddns"

	//ACMEDomain      = "acme.cellebyte.de"
	ACMEStagingURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

type CertBotParameters struct {
	Identifier          string
	Validation          string
	Token               string
	RemainingChallenges int64
	AllIdentifiers      []string
}
type ACMEConfig struct {
	CacheDir      string
	MailAddresses []string
	ACMEURL       string
	Timeout       time.Duration
	RenewBefore   time.Duration
}

func ParseACMEConfig() (config ACMEConfig, err error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return config, fmt.Errorf("getting default homedir: %w", err)
	}
	config.CacheDir = path.Join(homedir, defaultACMECacheDir)

	mailAddresses := strings.Split(os.Getenv(acmeMailAddresses), ",")
	if len(mailAddresses) == 0 {
		return config, fmt.Errorf("%s is required", acmeMailAddresses)
	}
	// TODO: swap this for production
	config.ACMEURL = ACMEStagingURL

	// Timeout
	config.Timeout = 1 * time.Minute
	// Check if we have user configured timeout
	timeoutString := os.Getenv(acmeChallengeTimeout)
	if timeoutString != "" {
		config.Timeout, err = time.ParseDuration(timeoutString)
		if err != nil {
			return config, fmt.Errorf("%s=%s is invalid: %w", acmeChallengeTimeout, timeoutString, err)
		}
	}
	config.RenewBefore = 30 * 24 * time.Hour
	renewDurationString := os.Getenv(acmeRenewBefore)
	if renewDurationString != "" {
		config.RenewBefore, err = time.ParseDuration(renewDurationString)
		if err != nil {
			return config, fmt.Errorf("%s=%s is invalid: %w", acmeRenewBefore, renewDurationString, err)
		}
	}
	return config, err
}

func ParseParams() (params CertBotParameters, err error) {
	params.Identifier = os.Getenv(certBotIdentifier)
	if params.Identifier == "" {
		return params, fmt.Errorf("%s is required", certBotIdentifier)
	}
	params.Validation = os.Getenv(certBotValidation)
	if params.Validation == "" {
		return params, fmt.Errorf("%s is required", certBotValidation)
	}
	params.Token = os.Getenv(certBotToken)
	remainChallenges := os.Getenv(certBotRemainingChallenges)
	if remainChallenges != "" {
		params.RemainingChallenges, err = strconv.ParseInt(remainChallenges, 10, 0)
		if err != nil {
			return params, fmt.Errorf("%s not a number: %w", certBotRemainingChallenges, err)
		}
	}
	allIdentifiers := os.Getenv(certBotRemainingChallenges)
	if allIdentifiers != "" {
		params.AllIdentifiers = strings.Split(allIdentifiers, ",")
	}
	return params, err
}
