package auth

import (
	"strings"
)

type (
	scopeCtxKey struct{}
)

const (
	scopeDelimiter = " "
)

// Checks if required scope is in claim
// We're using interface{} and casting it if needed to simplify usage of the fn by directly
// using it with map[string]interface{} claims type
func CheckScope(claim interface{}, req string) bool {
	claimStr, ok := claim.(string)
	if !ok {
		return false
	}

	return strings.Contains(
		scopeDelimiter+claimStr+scopeDelimiter,
		scopeDelimiter+strings.TrimSpace(req)+scopeDelimiter,
	)
}
