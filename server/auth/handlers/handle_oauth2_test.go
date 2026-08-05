package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/auth/settings"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/go-chi/jwtauth"
	oauth2models "github.com/go-oauth2/oauth2/v4/models"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_oauth2AuthorizeSuccess(t *testing.T) {
	var (
		user = makeMockUser()

		req = &http.Request{
			Form:     url.Values{},
			PostForm: url.Values{},
		}

		oauthService oauth2Service
		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	tcc := []testingExpect{
		{
			name:     "authorization success",
			payload:  -1,
			err:      "",
			template: "",
			fn: func(_ *settings.Settings) {
				oauthService = &oauth2ServiceMocked{
					handleAuthorizeRequest: func(w http.ResponseWriter, r *http.Request) error {
						return nil
					},
				}
			},
		},
		{
			name:     "authorization failure",
			payload:  http.StatusInternalServerError,
			err:      "not authorized",
			template: TmplInternalError,
			fn: func(_ *settings.Settings) {
				oauthService = &oauth2ServiceMocked{
					handleAuthorizeRequest: func(w http.ResponseWriter, r *http.Request) error {
						return fmt.Errorf("not authorized")
					},
				}
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			authSettings := &settings.Settings{}

			tc.fn(authSettings)

			authHandlers = &AuthHandlers{
				Log:         zap.NewNop(),
				AuthService: authService,
				Settings:    authSettings,
				OAuth2:      oauthService,
			}
			authReq = prepareClientAuthReq(authHandlers, req, user)

			err := authHandlers.oauth2Authorize(authReq)

			rq.NoError(err)
			rq.Equal(tc.template, authReq.Template)
			rq.Equal(tc.payload, authReq.Status)

			if tc.err != "" {
				rq.EqualError(errors.New(tc.err), authReq.Data["error"].(error).Error())
			}
		})
	}
}

func Test_oauth2AuthorizeSuccessSetParams(t *testing.T) {
	var (
		user = makeMockUser()

		req = &http.Request{
			Form:     url.Values{},
			PostForm: url.Values{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}

		rq = require.New(t)
	)

	oauthService := &oauth2ServiceMocked{
		handleAuthorizeRequest: func(w http.ResponseWriter, r *http.Request) error {
			return nil
		},
	}

	authHandlers = &AuthHandlers{
		Log:         zap.NewNop(),
		AuthService: authService,
		Settings:    authSettings,
		OAuth2:      oauthService,
	}
	authReq = prepareClientAuthReq(authHandlers, req, user)
	authReq.Session.Values["oauth2AuthParams"] = url.Values{"foo": []string{"bar"}}

	err := authHandlers.oauth2Authorize(authReq)

	rq.NoError(err)
	rq.Equal("", authReq.Template)
	rq.Equal(-1, authReq.Status)
	rq.Equal(nil, authReq.Data["error"])
}

func Test_generateIdToken(t *testing.T) {
	var (
		req = require.New(t)
	)

	signed, err := generateIdToken(
		&types.User{},
		&types.AuthClient{
			Secret: "correct horse battery stable",
		},
		&oauth2models.Token{},
		"http://cortezaproject.org",
	)

	req.NoError(err)
	req.NotEmpty(signed)
}

func Test_generateIdTokenClaims(t *testing.T) {
	var (
		updatedAt = time.Now()

		user = &types.User{
			ID:             42,
			Name:           "John Doe",
			Handle:         "jdoe",
			Email:          "john@example.tld",
			EmailConfirmed: true,
			UpdatedAt:      &updatedAt,
			Meta:           &types.UserMeta{PreferredLanguage: "sl"},
		}

		client = &types.AuthClient{ID: 1, Secret: "correct horse battery stable"}
	)

	tcc := []struct {
		name    string
		scope   string
		present []string
		absent  []string
	}{
		{
			name:    "without profile scope",
			scope:   "openid",
			present: []string{"sub", "user_id", "email", "email_verified", "iat", "aud"},
			absent:  []string{"name", "preferred_username", "locale", "updated_at"},
		},
		{
			name:    "with profile scope",
			scope:   "openid profile",
			present: []string{"sub", "email", "name", "preferred_username", "locale", "updated_at"},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
				ti  = &oauth2models.Token{}
			)

			ti.SetScope(tc.scope)

			signed, err := generateIdToken(user, client, ti, "https://example.tld/auth")
			req.NoError(err)

			token, err := jwt.Parse(signed)
			req.NoError(err)

			for _, claim := range tc.present {
				_, has := token.Get(claim)
				req.True(has, "expecting claim %q", claim)
			}

			for _, claim := range tc.absent {
				_, has := token.Get(claim)
				req.False(has, "not expecting claim %q", claim)
			}

			req.Equal("https://example.tld/auth", token.Issuer())
			req.Equal("42", token.Subject())
		})
	}
}

func Test_oauth2Userinfo(t *testing.T) {
	var (
		updatedAt = time.Now()

		user = &types.User{
			ID:             42,
			Name:           "John Doe",
			Handle:         "jdoe",
			Email:          "john@example.tld",
			EmailConfirmed: true,
			UpdatedAt:      &updatedAt,
			Meta:           &types.UserMeta{PreferredLanguage: "sl"},
		}
	)

	// swap the global token issuer for one we can sign test tokens with
	original := auth.TokenIssuer
	t.Cleanup(func() { auth.TokenIssuer = original })

	issuer, err := auth.NewTokenIssuer(
		auth.WithSecretSigner("correct horse battery stable"),
		auth.WithLookup(func(ctx context.Context, tokenID string) error { return nil }),
	)
	require.NoError(t, err)
	auth.TokenIssuer = issuer

	tcc := []struct {
		name     string
		scope    string
		noToken  bool
		notFound bool
		status   int
		claims   map[string]interface{}
		absent   []string
	}{
		{
			name:   "openid only",
			scope:  "openid",
			status: http.StatusOK,
			claims: map[string]interface{}{"sub": "42"},
			absent: []string{"name", "preferred_username", "locale", "email", "email_verified"},
		},
		{
			name:   "profile scope",
			scope:  "openid profile",
			status: http.StatusOK,
			claims: map[string]interface{}{
				"sub":                "42",
				"name":               "John Doe",
				"preferred_username": "jdoe",
				"locale":             "sl",
			},
			absent: []string{"email", "email_verified"},
		},
		{
			name:   "email scope",
			scope:  "openid email",
			status: http.StatusOK,
			claims: map[string]interface{}{
				"sub":            "42",
				"email":          "john@example.tld",
				"email_verified": true,
			},
			absent: []string{"name", "preferred_username", "locale"},
		},
		{
			name:    "no token",
			noToken: true,
			status:  http.StatusUnauthorized,
		},
		{
			name:     "unknown user",
			scope:    "openid",
			notFound: true,
			status:   http.StatusUnauthorized,
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				rq  = require.New(t)
				r   = httptest.NewRequest(http.MethodGet, "/oauth2/userinfo", nil)
				rec = httptest.NewRecorder()
			)

			if !tc.noToken {
				signed, err := auth.TokenIssuer.Sign(
					auth.WithIdentity(auth.Authenticated(user.ID)),
					auth.WithScope(strings.Split(tc.scope, " ")...),
				)
				rq.NoError(err)

				token, err := jwt.Parse(signed)
				rq.NoError(err)

				r = r.WithContext(jwtauth.NewContext(r.Context(), token, nil))
			}

			h := &AuthHandlers{
				Log: zap.NewNop(),
				UserService: userServiceMocked{
					findByAny: func(context.Context, interface{}) (*types.User, error) {
						if tc.notFound {
							return nil, fmt.Errorf("not found")
						}

						return user, nil
					},
				},
			}

			h.oauth2Userinfo(rec, r)
			rq.Equal(tc.status, rec.Code)

			if tc.status != http.StatusOK {
				rq.Contains(rec.Header().Get("WWW-Authenticate"), "Bearer")
				return
			}

			var response map[string]interface{}
			rq.NoError(json.Unmarshal(rec.Body.Bytes(), &response))

			for claim, value := range tc.claims {
				rq.Equal(value, response[claim], "claim %q", claim)
			}

			for _, claim := range tc.absent {
				rq.NotContains(response, claim)
			}
		})
	}
}

func Test_SubSplitRoles(t *testing.T) {
	type (
		exp struct {
			id string
			i  string
			ii []string
		}
	)
	var (
		req = require.New(t)
		ti  = &oauth2models.Token{}
		d   = make(map[string]interface{})

		tii = []exp{
			{id: "1", i: "1", ii: []string{}},
			{id: "1 2", i: "1", ii: []string{"2"}},
			{id: "1 2 33 444", i: "1", ii: []string{"2", "33", "444"}},
		}
	)

	for _, v := range tii {
		ti.SetUserID(v.id)
		SubSplit(ti, d)

		req.Equal(v.i, d["sub"])

		if _, is := d["roles"]; is {
			req.Equal(v.ii, d["roles"])
		}
	}
}
