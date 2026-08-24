package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/api/cred_registry"
	"github.com/crusttech/human/server/system/types"
)

// This file owns the server side of the generic OAuth2 authorization-code "Connect"
// flow. The auth-handlers mux drives the browser round-trip and single-use state;
// these methods build the consent URL and complete the exchange. Everything
// provider-specific is read off the connector's ConnectionAuth block + the instance
// OAuth app registry — no per-provider code.

// oauth2ConfigFor builds an oauth2.Config from the instance OAuth app record
// (client id/secret + endpoints + shared redirect — all provider-level) and the
// connector blueprint (per-connector scopes).
func oauth2ConfigFor(app *types.ConnectionOAuthApp, auth types.ConnectionAuth, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       auth.Scopes,
		Endpoint:     oauth2.Endpoint{AuthURL: app.AuthURL, TokenURL: app.TokenURL},
	}
}

// parseAuthParams decodes the provider record's authParams JSON (extra consent-URL
// params, e.g. access_type=offline). Malformed or empty yields no params.
func parseAuthParams(s string) map[string]string {
	if s == "" {
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil
	}
	return m
}

// findOAuthOption returns the connector's oauth2_authorization_code auth block,
// whether it is the connector's single auth method or one of its authOptions.
// "Connect" means "use the oauth method", so this does not depend on a prior
// authMethod selection on the configured connection.
func findOAuthOption(conn *types.Connection) (types.ConnectionAuth, bool) {
	if conn.Service.Auth.Method == "oauth2_authorization_code" {
		return conn.Service.Auth, true
	}
	for _, opt := range conn.Service.AuthOptions {
		if opt.Method == "oauth2_authorization_code" {
			return opt, true
		}
	}
	return types.ConnectionAuth{}, false
}

// resolveOAuth2 loads the configured connection (with RBAC + live connection) and
// resolves the oauth2_authorization_code auth block + instance OAuth app.
func (svc *configuredConnection) resolveOAuth2(ctx context.Context, ccID uint64) (cc *types.ConfiguredConnection, auth types.ConnectionAuth, app *types.ConnectionOAuthApp, redirect string, err error) {
	if cc, err = svc.FindByID(ctx, ccID); err != nil {
		return
	}

	var ok bool
	if auth, ok = findOAuthOption(&cc.Connection); !ok {
		err = errors.InvalidData("connection does not offer oauth2_authorization_code")
		return
	}

	app, redirect, ok = CurrentSettings.ConnectionOAuthApp(auth.OAuthApp)
	if !ok || redirect == "" {
		err = errors.InvalidData("no OAuth app configured for provider: %s", auth.OAuthApp)
	}
	return
}

// OAuth2AuthorizeURL builds the provider consent URL for a configured connection:
// scopes + any blueprint authParams + a PKCE challenge (when declared) + the given
// single-use state. verifier is the PKCE verifier the caller must keep for callback.
func (svc *configuredConnection) OAuth2AuthorizeURL(ctx context.Context, ccID uint64, state, verifier string) (string, error) {
	_, auth, app, redirect, err := svc.resolveOAuth2(ctx, ccID)
	if err != nil {
		return "", err
	}

	cfg := oauth2ConfigFor(app, auth, redirect)

	authParams := parseAuthParams(app.AuthParams)
	opts := make([]oauth2.AuthCodeOption, 0, len(authParams)+1)
	for k, v := range authParams {
		opts = append(opts, oauth2.SetAuthURLParam(k, v))
	}
	if app.PKCE {
		opts = append(opts, oauth2.S256ChallengeOption(verifier))
	}

	return cfg.AuthCodeURL(state, opts...), nil
}

// OAuth2Complete exchanges the authorization code and links the resulting delegated
// token to the configured connection. It applies Model B: one credential per
// (oauthApp, user), reused across that provider's connectors — so a user connects a
// provider once and its connectors share the token (scopes accumulate via the
// blueprint's include_granted_scopes). Returns the consenting account email, if any.
func (svc *configuredConnection) OAuth2Complete(ctx context.Context, ccID uint64, code, verifier string, userID uint64) (email string, err error) {
	cc, auth, app, redirect, err := svc.resolveOAuth2(ctx, ccID)
	if err != nil {
		return "", err
	}

	cfg := oauth2ConfigFor(app, auth, redirect)

	var opts []oauth2.AuthCodeOption
	if app.PKCE {
		opts = append(opts, oauth2.VerifierOption(verifier))
	}

	tok, err := cfg.Exchange(ctx, code, opts...)
	if err != nil {
		return "", fmt.Errorf("oauth2 token exchange failed: %w", err)
	}

	email = fetchIdentityEmail(ctx, app.IdentityURL, app.IdentityEmailField, tok.AccessToken)

	// Model B: find-or-create one credential per (oauthApp, user).
	rec, err := svc.upsertOAuth2Credential(ctx, userID, auth.OAuthApp, email, tok.RefreshToken)
	if err != nil {
		return "", err
	}

	// Link the shared credential + method to this configured connection.
	cc.Config.CredentialID = rec.ID
	cc.Config.AuthMethod = "oauth2_authorization_code"
	cc.UpdatedAt = now()
	if err = store.UpdateConfiguredConnection(ctx, svc.store, cc); err != nil {
		return "", err
	}

	// Load into the registry so the connection is usable immediately and the
	// refresher keeps it alive; route rotated refresh tokens to the shared record.
	cred := cred_registry.NewOAuth2AuthCodeCredential(
		cc.ID, app.ClientID, app.ClientSecret, app.TokenURL, auth.Scopes,
		tok.AccessToken, tok.RefreshToken, tok.Expiry,
	)
	credID := rec.ID
	s := svc.store
	cred.SetOnRotate(func(ctx context.Context, refreshToken string) error {
		c, err := store.LookupCredentialByID(ctx, s, credID)
		if err != nil {
			return err
		}
		c.Credentials = refreshToken
		return store.UpdateCredential(ctx, s, c)
	})
	_ = cred_registry.Default().Store(cred)

	return email, nil
}

// upsertOAuth2Credential returns the single credential for (oauthApp, user),
// creating it or updating its refresh token/label. The refresh token is only
// overwritten when the provider returned a new one (Google omits it on re-consent
// without prompt=consent).
func (svc *configuredConnection) upsertOAuth2Credential(ctx context.Context, userID uint64, oauthApp, email, refreshToken string) (*types.Credential, error) {
	kind := "oauth2:" + oauthApp
	label := oauthApp
	if email != "" {
		label = oauthApp + " " + email
	}

	set, _, err := store.SearchCredentials(ctx, svc.store, types.CredentialFilter{OwnerID: userID, Kind: kind})
	if err != nil {
		return nil, err
	}
	for _, c := range set {
		if c.DeletedAt != nil {
			continue
		}
		if refreshToken != "" {
			c.Credentials = refreshToken
		}
		if email != "" {
			c.Label = label
		}
		c.UpdatedAt = now()
		if err = store.UpdateCredential(ctx, svc.store, c); err != nil {
			return nil, err
		}
		return c, nil
	}

	rec := &types.Credential{
		ID:          nextID(),
		OwnerID:     userID,
		Kind:        kind,
		Label:       label,
		Credentials: refreshToken,
		CreatedAt:   *now(),
	}
	if err = store.CreateCredential(ctx, svc.store, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// fetchIdentityEmail calls the blueprint's identity endpoint (when declared) with
// the fresh access token and returns the email from the declared field. Best-effort:
// any failure yields an empty label rather than an error.
func fetchIdentityEmail(ctx context.Context, identityURL, emailField, accessToken string) string {
	if identityURL == "" || emailField == "" {
		return ""
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, identityURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var m map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return ""
	}
	if v, ok := m[emailField].(string); ok {
		return v
	}
	return ""
}
