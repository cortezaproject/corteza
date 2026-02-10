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

	"go.uber.org/zap"
)

type oauth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (r *registry) refreshToken(ctx context.Context, cred *Credential) error {
	if cred.AuthType != "oauth2_client_credentials" {
		return fmt.Errorf("credential type is not OAuth2")
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", cred.ClientID)
	data.Set("client_secret", cred.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", cred.TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.httpClient.Do(req)
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

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	cred.UpdateToken(tokenResp.AccessToken, expiresAt)

	// Persist the updated token to database
	if err := r.persistTokenUpdate(ctx, cred); err != nil {
		r.logger.Warn("failed to persist token update",
			zap.Uint64("connectionID", cred.ConnectionID),
			zap.Error(err),
		)
	}

	return nil
}
