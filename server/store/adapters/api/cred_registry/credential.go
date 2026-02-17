package cred_registry

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	// OAuth2TokenRefreshBuffer is the time before expiry when token refresh should occur
	OAuth2TokenRefreshBuffer = 5 * time.Minute
)

type (
	// Credential defines the interface for all credential types.
	// Each auth method implements its own refresh and state management.
	Credential interface {
		ConnectionID() uint64
		AuthType() string
		GetAccessToken() string
		NeedsRefresh() bool
		Refresh(ctx context.Context, client *http.Client) error
		MarshalState() map[string]any
	}

	// CredentialConfig holds the parameters needed to construct a Credential.
	CredentialConfig struct {
		ConnectionID uint64
		AuthType     string
		Token        string
		APIKey       string
		ClientID     string
		ClientSecret string
		TokenURL     string

		// JWT bearer fields
		Issuer        string
		Subject       string
		Audience      string
		Scopes        []string
		PrivateKey    string
		TokenLifetime time.Duration
	}
)

// NewCredential constructs the appropriate Credential implementation
// based on the auth type specified in the config.
func NewCredential(cfg CredentialConfig) (Credential, error) {
	switch cfg.AuthType {
	case "bearer":
		return NewBearerCredential(cfg.ConnectionID, cfg.Token), nil
	case "apikey":
		return NewAPIKeyCredential(cfg.ConnectionID, cfg.APIKey, "", ""), nil
	case "oauth2_client_credentials":
		return NewOAuth2ClientCredsCredential(cfg.ConnectionID, cfg.ClientID, cfg.ClientSecret, cfg.TokenURL), nil
	case "jwt_bearer":
		return NewJWTBearerCredential(cfg.ConnectionID, cfg.Issuer, cfg.Subject, cfg.Audience, cfg.TokenURL, cfg.PrivateKey, cfg.Scopes, cfg.TokenLifetime), nil
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", cfg.AuthType)
	}
}
