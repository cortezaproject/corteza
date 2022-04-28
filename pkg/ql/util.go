package ql

import (
	"regexp"
	"strings"
)

var (
	truthy = regexp.MustCompile(`^(t(rue)?|y(es)?|1)$`)
)

// Check what boolean value the given string conforms to
func evalBool(v string) bool {
	return truthy.MatchString(strings.ToLower(v))
}

// Check if the given string is a float
func isFloaty(v string) bool {
	return strings.Contains(v, ".")
}
