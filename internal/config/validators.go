package config

import (
	"regexp"
	"strings"
	"sync"
)

const (
	// Regexes
	dnsRegexStringRFC1035Label = "^[a-z]([-a-z0-9]*[a-z0-9])?$"
)

func lazyRegexCompile(str string) func() *regexp.Regexp {
	var regex *regexp.Regexp
	var once sync.Once
	return func() *regexp.Regexp {
		once.Do(func() {
			regex = regexp.MustCompile(str)
		})
		return regex
	}
}

var (
	dnsRegexRFC1035Label = lazyRegexCompile(dnsRegexStringRFC1035Label)
)

func ValidRFC1035Label(dnsLabel string) bool {
	size := len(dnsLabel)
	if size > 63 {
		return false
	}

	return dnsRegexRFC1035Label().MatchString(dnsLabel)
}

func ValidRFC1035Name(dnsName string) bool {
	size := len(dnsName)
	if size > 253 {
		return false
	}
	for _, dnsLabel := range strings.Split(dnsName, ".") {
		if ok := ValidRFC1035Label(dnsLabel); !ok {
			return false
		}
	}

	return true
}

func ValidRecordType(recordType string) bool {
	switch recordType {
	case "A":
		fallthrough
	case "AAAA":
		return true
	}
	return false
}

func ValidRecordTypes(recordTypes []string) bool {
	size := len(recordTypes)
	if size > 2 {
		return false
	}
	for _, recordType := range recordTypes {
		if ok := ValidRecordType(recordType); !ok {
			return false
		}
	}
	return true
}
