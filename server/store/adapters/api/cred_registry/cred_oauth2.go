package cred_registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// OAuth2ClientCredsCredential handles OAuth2 client credentials grant.
type OAuth2ClientCredsCredential struct {
	connID       uint64
	ClientID     string
	ClientSecret string
	TokenURL     string

	// Mutable runtime state
	AccessToken string
	ExpiresAt   time.Time
}

func NewOAuth2ClientCredsCredential(connID uint64, clientID, clientSecret, tokenURL string) *OAuth2ClientCredsCredential {
	return &OAuth2ClientCredsCredential{
		connID:       connID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		// Force immediate refresh on first use
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
}

func (c *OAuth2ClientCredsCredential) ConnectionID() uint64  { return c.connID }
func (c *OAuth2ClientCredsCredential) AuthType() string       { return "oauth2_client_credentials" }
func (c *OAuth2ClientCredsCredential) GetAccessToken() string { return c.AccessToken }

func (c *OAuth2ClientCredsCredential) NeedsRefresh() bool {
	return time.Now().Add(OAuth2TokenRefreshBuffer).After(c.ExpiresAt)
}

func (c *OAuth2ClientCredsCredential) MarshalState() map[string]any {
	return map[string]any{
		"accessToken":    c.AccessToken,
		"tokenExpiresAt": c.ExpiresAt.Format(time.RFC3339),
	}
}

type oauth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *OAuth2ClientCredsCredential) Refresh(ctx context.Context, client *http.Client) error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", c.TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp oauth2TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	c.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}
