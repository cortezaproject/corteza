package cred_registry

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// DefaultTokenLifetime is the default JWT assertion lifetime
	DefaultTokenLifetime = time.Hour
)

// JWTBearerCredential handles the JWT bearer assertion grant
// (grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer).
//
// Works with any provider that accepts a signed JWT assertion
// (e.g. Google service accounts, Azure AD, custom IdPs).
type JWTBearerCredential struct {
	connID uint64

	// Static config
	Issuer        string        // JWT "iss" claim (e.g. service account email)
	Subject       string        // JWT "sub" claim (optional, for domain-wide delegation)
	Audience      string        // JWT "aud" claim (token endpoint URL)
	Scopes        []string      // space-joined into JWT "scope" claim
	PrivateKey    string        // PEM-encoded RSA private key
	TokenURL      string        // token endpoint
	TokenLifetime time.Duration // JWT assertion lifetime (default: 1h)

	// Mutable runtime state
	AccessToken string
	ExpiresAt   time.Time

	parsedKey *rsa.PrivateKey
}

func NewJWTBearerCredential(connID uint64, issuer, subject, audience, tokenURL, privateKey string, scopes []string, tokenLifetime time.Duration) *JWTBearerCredential {
	if tokenLifetime == 0 {
		tokenLifetime = DefaultTokenLifetime
	}
	return &JWTBearerCredential{
		connID:        connID,
		Issuer:        issuer,
		Subject:       subject,
		Audience:      audience,
		TokenURL:      tokenURL,
		PrivateKey:    privateKey,
		Scopes:        scopes,
		TokenLifetime: tokenLifetime,
		// Force immediate refresh on first use
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
}

func (c *JWTBearerCredential) ConnectionID() uint64   { return c.connID }
func (c *JWTBearerCredential) AuthType() string       { return "jwt_bearer" }
func (c *JWTBearerCredential) GetAccessToken() string { return c.AccessToken }

func (c *JWTBearerCredential) NeedsRefresh() bool {
	return time.Now().Add(OAuth2TokenRefreshBuffer).After(c.ExpiresAt)
}

func (c *JWTBearerCredential) MarshalState() map[string]any {
	return map[string]any{
		"accessToken":    c.AccessToken,
		"tokenExpiresAt": c.ExpiresAt.Format(time.RFC3339),
	}
}

func (c *JWTBearerCredential) Refresh(ctx context.Context, client *http.Client) error {
	if err := c.ensureKey(); err != nil {
		return err
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": c.Issuer,
		"aud": c.Audience,
		"iat": now.Unix(),
		"exp": now.Add(c.TokenLifetime).Unix(),
	}
	if c.Subject != "" {
		claims["sub"] = c.Subject
	}
	if len(c.Scopes) > 0 {
		scope := ""
		for i, s := range c.Scopes {
			if i > 0 {
				scope += " "
			}
			scope += s
		}
		claims["scope"] = scope
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(c.parsedKey)
	if err != nil {
		return fmt.Errorf("failed to sign JWT assertion: %w", err)
	}

	// Exchange assertion for access token
	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	data.Set("assertion", signed)

	req, err := http.NewRequestWithContext(ctx, "POST", c.TokenURL, strings.NewReader(data.Encode()))
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

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	if tokenResp.ExpiresIn > 0 {
		c.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	} else {
		// Default to 1 hour if not specified
		c.ExpiresAt = time.Now().Add(time.Hour)
	}

	return nil
}

func (c *JWTBearerCredential) ensureKey() error {
	if c.parsedKey != nil {
		return nil
	}

	block, _ := pem.Decode([]byte(c.PrivateKey))
	if block == nil {
		return fmt.Errorf("failed to decode PEM block from private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1 as fallback
		rsaKey, err2 := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err2 != nil {
			return fmt.Errorf("failed to parse private key (tried PKCS8 and PKCS1): %w", err)
		}
		c.parsedKey = rsaKey
		return nil
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("private key is not RSA")
	}
	c.parsedKey = rsaKey
	return nil
}
