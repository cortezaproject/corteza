package sam

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/oauth2"
	"net/http"
)

type (
	openIdConnect struct {
		provider *oidc.Provider
		verifier *oidc.IDTokenVerifier
		config   oauth2.Config
	}
)

func OpenIdConnect(ctx context.Context, issuer string, cfg oauth2.Config) (c *openIdConnect, err error) {
	c = &openIdConnect{}

	c.provider, err = oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	// Configure an OpenID Connect aware OAuth2 client.
	c.config = oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: c.provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	c.verifier = c.provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return
}

func (c *openIdConnect) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	state := "@todo"
	http.Redirect(w, r, c.config.AuthCodeURL(state), http.StatusFound)
}

func (c *openIdConnect) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	// @todo check state

	// Verify state and errors.

	oauth2Token, err := c.config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		// handle error
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		// handle missing token
	}

	// Parse and verify ID Token payload.
	idToken, err := c.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		// handle error
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}

	if err := idToken.Claims(&claims); err != nil {
		// handle error
	}

	spew.Dump()
}
