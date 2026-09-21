package oauth2

import (
	"net/url"
	"path"
	"strings"

	"github.com/go-oauth2/oauth2/v4/errors"
)

const (
	// RedirectUriAny allows redirection to any URI when put on the list of allowed redirect URIs
	RedirectUriAny = "*"
)

// AllowedRedirectURIs returns list of redirect URIs allowed for the client
//
// Clients without redirect URIs fall back to the given defaults
func AllowedRedirectURIs(clientRedirectURI string, def []string) []string {
	if allowed := strings.Fields(clientRedirectURI); len(allowed) > 0 {
		return allowed
	}

	return def
}

// ValidateRedirectURI verifies redirect URI against the list of allowed URIs
func ValidateRedirectURI(allowed []string, redirectURI string) error {
	for _, a := range allowed {
		if a == RedirectUriAny {
			return nil
		}
	}

	redirect, err := url.Parse(redirectURI)
	if err != nil || !redirect.IsAbs() || redirect.Host == "" || redirect.User != nil || strings.Contains(redirectURI, `\`) {
		return errors.ErrInvalidRedirectURI
	}

	for _, a := range allowed {
		if base, err := url.Parse(a); err == nil && matchRedirectURI(base, redirect) {
			return nil
		}
	}

	return errors.ErrInvalidRedirectURI
}

// matchRedirectURI compares scheme and host and makes sure redirect path is under the base path
func matchRedirectURI(base, redirect *url.URL) bool {
	if !base.IsAbs() || base.Host == "" {
		return false
	}

	if !strings.EqualFold(base.Scheme, redirect.Scheme) || !strings.EqualFold(base.Host, redirect.Host) {
		return false
	}

	basePath := strings.TrimSuffix(base.Path, "/")
	if basePath == "" {
		return true
	}

	// resolve dot segments before comparing
	redirectPath := path.Clean("/" + redirect.Path)

	return redirectPath == basePath || strings.HasPrefix(redirectPath, basePath+"/")
}
