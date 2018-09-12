package rest

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/crusttech/crust/auth/service"
	"github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/config"
	"github.com/titpetric/factory/resputil"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type (
	openIdConnect struct {
		provider *oidc.Provider
		verifier *oidc.IDTokenVerifier
		config   oauth2.Config

		appURL            string
		stateCookieExpiry int64

		userService service.UserService

		jwt jwtEncodeCookieSetter
	}

	jwtEncodeCookieSetter interface {
		types.TokenEncoder
		SetToCookie(w http.ResponseWriter, r *http.Request, identity types.Identifiable)
	}
)

const openIdConnectStateCookie = "oidc-state"

func OpenIdConnect(cfg *config.OIDC, usvc service.UserService, jwt jwtEncodeCookieSetter) (c *openIdConnect, err error) {
	c = &openIdConnect{
		appURL:            cfg.AppURL,
		stateCookieExpiry: cfg.StateCookieExpiry,
		userService:       usvc,
		jwt:               jwt,
	}

	// Allow 5 seconds for issuer discovery process
	c.provider, err = oidc.NewProvider(context.Background(), cfg.Issuer)
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

// Redirects user to the issuer's login screen
func (c *openIdConnect) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	// @todo sure we can improve this...
	rand.Seed(4321)
	var state = strconv.FormatInt(rand.Int63(), 10)

	// Store state to cookie as well
	c.setStateCookie(w, r, state)

	http.Redirect(w, r, c.config.AuthCodeURL(state), http.StatusFound)
}

// Handles callback from issuer
//
// If everything goes well (scope & token verification) it reads issued claims,
// creates Crust JWT and stores it in a cookie.
//
// @todo All failed responses must redirect to appURL as well + some error code that will be displayed on the client
func (c *openIdConnect) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	if !c.stateCheck(r) {
		resputil.JSON(w, "State check failed")
		return
	}

	c.setStateCookie(w, r, "") // remove state cookie

	oauth2Token, err := c.config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		resputil.JSON(w, err)
		return
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		resputil.JSON(w, err)
		return
	}

	// Parse and verify ID Token payload.
	idToken, err := c.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		resputil.JSON(w, err)
		return
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}

	if err := idToken.Claims(&claims); err != nil {
		resputil.JSON(w, err)
		return
	}

	var user *types.User

	if user, err = c.userService.FindOrCreate(claims.Email); err != nil {
		resputil.JSON(w, err)
		return
	} else {
		c.jwt.SetToCookie(w, r, user)
	}

	http.Redirect(w, r, c.appURL+"?jwt="+c.jwt.Encode(user), http.StatusSeeOther)
}

func (c *openIdConnect) stateCheck(r *http.Request) bool {
	if cState, err := r.Cookie(openIdConnectStateCookie); err == nil {
		rState := r.URL.Query().Get("state")
		return len(rState) > 0 && cState.Value == rState
	}

	return false
}

// Sets state cookie
func (c *openIdConnect) setStateCookie(w http.ResponseWriter, r *http.Request, value string) {
	var maxAge int

	if len(value) == 0 {
		// When empty string for a value is received,
		// set maxAge to -1. That will effectively delete the cookie
		maxAge = -1
	}

	// Store state to cookie as well
	http.SetCookie(w, &http.Cookie{
		Name:  openIdConnectStateCookie,
		Value: value,

		Expires:  time.Now().Add(time.Duration(c.stateCookieExpiry) * time.Minute),
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		Path:     "/oidc",
	})
}
