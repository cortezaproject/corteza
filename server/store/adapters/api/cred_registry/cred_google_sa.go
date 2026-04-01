package cred_registry

import (
	"context"
	"net/http"
	"time"
)

const (
	googleTokenURL = "https://oauth2.googleapis.com/token"
)

// GoogleServiceAccountCredential wraps JWTBearerCredential with
// Google-specific defaults (token URL, audience, default scopes).
type GoogleServiceAccountCredential struct {
	*JWTBearerCredential
}

func NewGoogleServiceAccountCredential(connID uint64, serviceAccountEmail, privateKey, subject string, scopes []string, tokenLifetime time.Duration) *GoogleServiceAccountCredential {

	return &GoogleServiceAccountCredential{
		JWTBearerCredential: NewJWTBearerCredential(
			connID,
			serviceAccountEmail, // issuer
			subject,             // subject (empty = no DWD; set to user email for impersonation)
			googleTokenURL,      // audience
			googleTokenURL,      // tokenURL
			privateKey,
			scopes,
			tokenLifetime,
		),
	}
}

func (c *GoogleServiceAccountCredential) AuthType() string { return "google_service_account" }

func (c *GoogleServiceAccountCredential) ConnectionID() uint64   { return c.JWTBearerCredential.ConnectionID() }
func (c *GoogleServiceAccountCredential) GetAccessToken() string { return c.JWTBearerCredential.GetAccessToken() }
func (c *GoogleServiceAccountCredential) NeedsRefresh() bool     { return c.JWTBearerCredential.NeedsRefresh() }
func (c *GoogleServiceAccountCredential) MarshalState() map[string]any {
	return c.JWTBearerCredential.MarshalState()
}

func (c *GoogleServiceAccountCredential) Refresh(ctx context.Context, client *http.Client) error {
	return c.JWTBearerCredential.Refresh(ctx, client)
}
