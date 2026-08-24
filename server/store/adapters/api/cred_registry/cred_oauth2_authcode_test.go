package cred_registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// stubTokenEndpoint returns a token endpoint that echoes a configurable response
// and records the grant_type / refresh_token it received.
func stubTokenEndpoint(t *testing.T, resp string, gotRefresh *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if gotRefresh != nil {
			*gotRefresh = r.Form.Get("refresh_token")
		}
		if r.Form.Get("grant_type") != "refresh_token" {
			t.Errorf("expected grant_type=refresh_token, got %q", r.Form.Get("grant_type"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(resp))
	}))
}

func TestOAuth2AuthCode_Refresh(t *testing.T) {
	var sentRefresh string
	srv := stubTokenEndpoint(t,
		`{"access_token":"new-access","refresh_token":"rotated-refresh","expires_in":3600,"token_type":"Bearer"}`,
		&sentRefresh)
	defer srv.Close()

	c := NewOAuth2AuthCodeCredential(1, "cid", "secret", srv.URL, []string{"scope"},
		"old-access", "old-refresh", time.Now().Add(-time.Minute))

	if !c.NeedsRefresh() {
		t.Fatal("expected NeedsRefresh true for an expired token with a refresh token")
	}

	if err := c.Refresh(context.Background(), srv.Client()); err != nil {
		t.Fatalf("Refresh: %v", err)
	}

	if sentRefresh != "old-refresh" {
		t.Errorf("endpoint received refresh_token %q, want old-refresh", sentRefresh)
	}
	if c.GetAccessToken() != "new-access" {
		t.Errorf("access token = %q, want new-access", c.GetAccessToken())
	}
	// Rotation: the new refresh token must replace the old one.
	if c.RefreshToken != "rotated-refresh" {
		t.Errorf("refresh token = %q, want rotated-refresh", c.RefreshToken)
	}
	if c.NeedsRefresh() {
		t.Error("expected NeedsRefresh false after a successful refresh")
	}
}

func TestOAuth2AuthCode_Refresh_PersistsRotatedToken(t *testing.T) {
	srv := stubTokenEndpoint(t,
		`{"access_token":"new-access","refresh_token":"rotated","expires_in":3600,"token_type":"Bearer"}`, nil)
	defer srv.Close()

	c := NewOAuth2AuthCodeCredential(1, "cid", "secret", srv.URL, nil,
		"old-access", "old-refresh", time.Now().Add(-time.Minute))

	var persisted string
	c.SetOnRotate(func(_ context.Context, refreshToken string) error {
		persisted = refreshToken
		return nil
	})

	if err := c.Refresh(context.Background(), srv.Client()); err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if persisted != "rotated" {
		t.Errorf("onRotate persisted %q, want rotated", persisted)
	}
}

func TestOAuth2AuthCode_BootCaseMintsAccessToken(t *testing.T) {
	srv := stubTokenEndpoint(t,
		`{"access_token":"minted","expires_in":3600,"token_type":"Bearer"}`, nil)
	defer srv.Close()

	// Loaded from storage: refresh token only, no access token.
	c := NewOAuth2AuthCodeCredential(1, "cid", "secret", srv.URL, nil, "", "stored-refresh", time.Time{})

	if !c.NeedsRefresh() {
		t.Fatal("expected NeedsRefresh true when only a refresh token is loaded")
	}
	if err := c.Refresh(context.Background(), srv.Client()); err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if c.GetAccessToken() != "minted" {
		t.Errorf("access token = %q, want minted", c.GetAccessToken())
	}
}

func TestOAuth2AuthCode_Refresh_KeepsOldRefreshWhenAbsent(t *testing.T) {
	srv := stubTokenEndpoint(t,
		`{"access_token":"new-access","expires_in":3600,"token_type":"Bearer"}`, nil)
	defer srv.Close()

	c := NewOAuth2AuthCodeCredential(1, "cid", "secret", srv.URL, nil,
		"old-access", "keep-me", time.Now().Add(-time.Minute))

	if err := c.Refresh(context.Background(), srv.Client()); err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	// No refresh_token in the response → the existing one is preserved.
	if c.RefreshToken != "keep-me" {
		t.Errorf("refresh token = %q, want keep-me (preserved)", c.RefreshToken)
	}
}

func TestOAuth2AuthCode_NoRefreshToken(t *testing.T) {
	c := NewOAuth2AuthCodeCredential(1, "cid", "secret", "https://example.invalid/token", nil,
		"access-only", "", time.Time{})

	// Non-expiring / no-refresh providers must not be flagged for refresh.
	if c.NeedsRefresh() {
		t.Error("expected NeedsRefresh false when there is no refresh token")
	}
	if err := c.Refresh(context.Background(), http.DefaultClient); err == nil {
		t.Error("expected Refresh to error when there is no refresh token")
	}
}

func TestOAuth2AuthCode_MarshalState(t *testing.T) {
	exp := time.Now().Add(time.Hour).Truncate(time.Second)
	c := NewOAuth2AuthCodeCredential(1, "cid", "secret", "https://x/token", nil, "a", "r", exp)

	st := c.MarshalState()
	if st["accessToken"] != "a" || st["refreshToken"] != "r" {
		t.Errorf("MarshalState missing tokens: %#v", st)
	}
	if _, err := time.Parse(time.RFC3339, st["tokenExpiresAt"].(string)); err != nil {
		t.Errorf("tokenExpiresAt not RFC3339: %v", err)
	}
}
