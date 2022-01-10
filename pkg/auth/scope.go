package auth

import (
	"strings"

	"github.com/lestrrat-go/jwx/jwt"
)

const (
	scopeDelimiter = " "
)

// CheckJwtScope verifies if required scope is in claim
// We're using interface{} and casting it if needed to simplify usage of the fn by directly
// using it with map[string]interface{} claims type
func CheckJwtScope(token jwt.Token, required ...string) bool {
	scopeClaimRaw, has := token.Get("scope")
	if !has {
		return false
	}

	scopeClaim, ok := scopeClaimRaw.(string)
	if !ok {
		return false
	}

	return CheckScope(scopeClaim, required...)
}

func CheckScope(scope string, required ...string) bool {
	scope = scopeDelimiter + strings.TrimSpace(scope) + scopeDelimiter

	for _, req := range required {
		req = scopeDelimiter + strings.TrimSpace(req) + scopeDelimiter
		if strings.Contains(scope, req) {
			return true
		}
	}

	return false
}
