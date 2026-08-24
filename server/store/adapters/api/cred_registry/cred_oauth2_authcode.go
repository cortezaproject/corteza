package cred_registry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// OAuth2AuthCodeCredential handles the OAuth2 authorization-code grant: a
// delegated, per-user token obtained once at consent and kept alive by the
// refresh token. It is provider-agnostic — every provider difference (token
// endpoint, scopes, client id/secret) is passed in from the connector blueprint
// and the instance OAuth app, never hardcoded here.
//
// The refresh/rotation mechanics are delegated to golang.org/x/oauth2, which
// autodetects the token-endpoint auth style (Basic vs body) and preserves the
// existing refresh token when a provider does not return a new one.
type OAuth2AuthCodeCredential struct {
	connID uint64
	config oauth2.Config

	// Mutable runtime state
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time

	// onRotate, when set, persists a rotated refresh token durably. Providers that
	// rotate (Zoom, Atlassian, …) invalidate the old refresh token on use, so the
	// new value must be saved immediately or the connection is locked out after a
	// restart.
	onRotate func(ctx context.Context, refreshToken string) error
}

// SetOnRotate registers a callback that durably persists a rotated refresh token.
func (c *OAuth2AuthCodeCredential) SetOnRotate(fn func(ctx context.Context, refreshToken string) error) {
	c.onRotate = fn
}

func NewOAuth2AuthCodeCredential(connID uint64, clientID, clientSecret, tokenURL string, scopes []string, accessToken, refreshToken string, expiresAt time.Time) *OAuth2AuthCodeCredential {
	return &OAuth2AuthCodeCredential{
		connID: connID,
		config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			Endpoint:     oauth2.Endpoint{TokenURL: tokenURL},
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
}

func (c *OAuth2AuthCodeCredential) ConnectionID() uint64   { return c.connID }
func (c *OAuth2AuthCodeCredential) AuthType() string       { return "oauth2_authorization_code" }
func (c *OAuth2AuthCodeCredential) GetAccessToken() string { return c.AccessToken }

func (c *OAuth2AuthCodeCredential) NeedsRefresh() bool {
	// Nothing to refresh with (e.g. providers that issue non-expiring tokens and
	// no refresh token) — use the access token as-is until it is rejected.
	if c.RefreshToken == "" {
		return false
	}
	// Loaded from durable storage with only the refresh token (e.g. on boot): mint
	// an access token before first use.
	if c.AccessToken == "" {
		return true
	}
	// A zero expiry means the provider did not declare one; treat as non-expiring.
	if c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().Add(OAuth2TokenRefreshBuffer).After(c.ExpiresAt)
}

func (c *OAuth2AuthCodeCredential) Refresh(ctx context.Context, client *http.Client) error {
	if c.RefreshToken == "" {
		return fmt.Errorf("no refresh token available; reconnect required")
	}

	if client != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
	}

	// Force x/oauth2 to run the refresh now: its TokenSource only refreshes once
	// the token is invalid, and our proactive buffer fires before x/oauth2's own
	// (10s) delta would. Handing it an already-expired copy triggers the refresh
	// while keeping the real refresh token.
	src := c.config.TokenSource(ctx, &oauth2.Token{
		AccessToken:  c.AccessToken,
		RefreshToken: c.RefreshToken,
		Expiry:       time.Now().Add(-time.Hour),
	})

	tok, err := src.Token()
	if err != nil {
		return fmt.Errorf("failed to refresh oauth2 token: %w", err)
	}

	c.AccessToken = tok.AccessToken
	c.ExpiresAt = tok.Expiry

	// Capture a rotated refresh token; providers that rotate (Zoom, Atlassian, …)
	// invalidate the old one, so the new value must replace it and be persisted
	// durably. x/oauth2 already carries the previous token forward when the
	// response omits a new one, so a changed value here means a real rotation.
	if tok.RefreshToken != "" && tok.RefreshToken != c.RefreshToken {
		c.RefreshToken = tok.RefreshToken
		if c.onRotate != nil {
			if perr := c.onRotate(ctx, c.RefreshToken); perr != nil {
				return fmt.Errorf("failed to persist rotated refresh token: %w", perr)
			}
		}
	}

	return nil
}

func (c *OAuth2AuthCodeCredential) MarshalState() map[string]any {
	return map[string]any{
		"accessToken":    c.AccessToken,
		"refreshToken":   c.RefreshToken,
		"tokenExpiresAt": c.ExpiresAt.Format(time.RFC3339),
	}
}
