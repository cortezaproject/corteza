package cred_registry

import (
	"time"
)

const (
	// OAuth2TokenRefreshBuffer is the time before expiry when token refresh should occur
	OAuth2TokenRefreshBuffer = 5 * time.Minute
)

type (
	Credential struct {
		ConnectionID uint64

		// AuthType: bearer, apikey, oauth2_client_credentials
		AuthType string

		// Static credentials
		Token  string // For bearer tokens
		APIKey string // For API key auth

		// OAuth2 client credentials
		ClientID     string
		ClientSecret string
		TokenURL     string

		// Dynamic OAuth2 state
		AccessToken string
		ExpiresAt   time.Time
	}
)

func (c *Credential) GetAccessToken() string {
	switch c.AuthType {
	case "bearer":
		return c.Token
	case "apikey":
		return c.APIKey
	case "oauth2_client_credentials":
		return c.AccessToken
	default:
		return ""
	}
}

func (c *Credential) UpdateToken(token string, expiresAt time.Time) {
	c.AccessToken = token
	c.ExpiresAt = expiresAt
}

func (c *Credential) NeedsRefresh() bool {
	if c.AuthType != "oauth2_client_credentials" {
		return false
	}
	return time.Now().Add(OAuth2TokenRefreshBuffer).After(c.ExpiresAt)
}
