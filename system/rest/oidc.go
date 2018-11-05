package rest

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/system/repository"
	"github.com/crusttech/crust/system/service"
	"github.com/crusttech/crust/system/types"
	"github.com/crusttech/go-oidc"
	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"
	"golang.org/x/oauth2"
)

const DB_SETTINGS_KEY_OIDC_CLIENT = "oidc-client"

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

	oidcProfile struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Sub   string `json:"sub"`
	}

	jwtEncodeCookieSetter interface {
		auth.TokenEncoder
		SetCookie(w http.ResponseWriter, r *http.Request, identity auth.Identifiable)
	}
)

const openIdConnectStateCookie = "oidc-state"

// Sets-up OIDC connection (issuer discovery, client registration)
//
// Client registration is done when no cfg.ClientID is provided.
func OpenIdConnect(ctx context.Context, cfg *config.OIDC, usvc service.UserService, jwt jwtEncodeCookieSetter, settings repository.Settings) (c *openIdConnect, err error) {
	c = &openIdConnect{
		appURL:            cfg.AppURL,
		stateCookieExpiry: cfg.StateCookieExpiry,
		userService:       usvc,
		jwt:               jwt,
	}

	// Allow 5 seconds for issuer discovery process
	c.provider, err = oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, err
	}

	if len(cfg.ClientID) > 0 {
		// System is configured with fixed OIDC client ID (probably through AUTH_OIDC_CLIENT_ID)
		// Construct oauth2 config from provided configuration
		c.config = oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,

			// Discovery returns the OAuth2 endpoints.
			Endpoint: c.provider.Endpoint(),
		}
	} else {
		client := &oidc.Client{}

		if found, err := settings.Get(DB_SETTINGS_KEY_OIDC_CLIENT, client); err != nil {
			return nil, errors.Wrap(err, "could not load oidc client settings from the database")
		} else if !found {
			// Perform dynamic client registration
			client, err = c.provider.RegisterClient(ctx, &oidc.ClientRegistration{
				Name:          "Crust",
				RedirectURIs:  []string{cfg.RedirectURL},
				ResponseTypes: []string{"token id_token", "code"},
			})

			if err := settings.Set(DB_SETTINGS_KEY_OIDC_CLIENT, client); err != nil {
				return nil, errors.Wrap(err, "could not store oidc client settings from the database")
			}
		}

		c.config = c.provider.OAuth2Config(client)
	}

	c.config.Scopes = []string{
		oidc.ScopeOpenID,
		"email",
		"profile",
		"address",
		"phone_number",
	}

	c.verifier = c.provider.Verifier(&oidc.Config{ClientID: c.config.ClientID})

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
	_, err = c.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		resputil.JSON(w, err)
		return
	}

	u, _ := c.provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	p := &oidcProfile{}
	u.Claims(p)

	var user = &types.User{
		SatosaID: p.Sub,
		Email:    p.Email,
		Name:     p.Name,
	}

	if user, err = c.userService.With(ctx).FindOrCreate(user); err != nil {
		resputil.JSON(w, err)
		return
	} else {
		c.jwt.SetCookie(w, r, user)
	}

	http.Redirect(w, r, c.appURL, http.StatusSeeOther)
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
		Name:     openIdConnectStateCookie,
		Value:    value,
		Expires:  time.Now().Add(time.Duration(c.stateCookieExpiry) * time.Minute),
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		Path:     "/oidc",
		Domain:   ".rustbucket.io", // @todo make this configurable (like stateCookieExpiry)
	})
}
