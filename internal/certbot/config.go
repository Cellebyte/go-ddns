package certbot

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// ref: https://eff-certbot.readthedocs.io/en/stable/using.html#pre-and-post-validation-hooks
	certBotIdentifier          = "CERTBOT_IDENTIFIER"           // The domain or IP address being authenticated
	certBotValidation          = "CERTBOT_VALIDATION"           // The validation string
	certBotToken               = "CERTBOT_TOKEN"                // Resource name part of the HTTP-01 challenge (HTTP-01 only)
	certBotRemainingChallenges = "CERTBOT_REMAINING_CHALLENGES" // Number of challenges remaining after the current challenge
	certBotAllIdentifiers      = "CERTBOT_ALL_IDENTIFIERS"      // A comma-separated list of all identifiers challenged for the current certificate

	ACMETXTPrefix = "_acme-challenge"
)

type CertBotParameters struct {
	Identifier          string
	Validation          string
	Token               string
	RemainingChallenges int64
	AllIdentifiers      []string
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
