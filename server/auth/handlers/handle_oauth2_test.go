package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/auth/settings"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/options"
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

func Test_oauth2Userinfo(t *testing.T) {
	// Initialize TokenIssuer for tests with a mock lookup that always succeeds
	var err error
	auth.TokenIssuer, err = auth.NewTokenIssuer(
		auth.WithSecretSigner("test-secret"),
		auth.WithLookup(func(ctx context.Context, tokenID string) error {
			// Always return nil to indicate token is valid
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("failed to initialize token issuer: %v", err)
	}

	var (
		now       = time.Now()
		updatedAt = now.Add(-time.Hour)
	)

	testCases := []struct {
		name           string
		user           *types.User
		scope          string
		expectedStatus int
		expectedClaims map[string]interface{}
		userNotFound   bool
		noToken        bool
	}{
		{
			name: "full profile and email claims",
			user: &types.User{
				ID:             12345,
				Name:           "John Doe",
				Handle:         "johndoe",
				Email:          "john@example.com",
				EmailConfirmed: true,
				CreatedAt:      now,
				UpdatedAt:      &updatedAt,
				Meta: &types.UserMeta{
					PreferredLanguage: "en",
				},
			},
			scope:          "openid profile email",
			expectedStatus: http.StatusOK,
			expectedClaims: map[string]interface{}{
				"sub":                "12345",
				"name":               "John Doe",
				"given_name":         "John",
				"family_name":        "Doe",
				"preferred_username": "johndoe",
				"email":              "john@example.com",
				"email_verified":     true,
				"locale":             "en",
			},
		},
		{
			name: "profile scope only",
			user: &types.User{
				ID:        12345,
				Name:      "Jane Smith",
				Handle:    "janesmith",
				Email:     "jane@example.com",
				CreatedAt: now,
				Meta:      &types.UserMeta{},
			},
			scope:          "openid profile",
			expectedStatus: http.StatusOK,
			expectedClaims: map[string]interface{}{
				"sub":                "12345",
				"name":               "Jane Smith",
				"preferred_username": "janesmith",
			},
		},
		{
			name: "email scope only",
			user: &types.User{
				ID:             12345,
				Email:          "test@example.com",
				EmailConfirmed: false,
				CreatedAt:      now,
				Meta:           &types.UserMeta{},
			},
			scope:          "openid email",
			expectedStatus: http.StatusOK,
			expectedClaims: map[string]interface{}{
				"sub":            "12345",
				"email":          "test@example.com",
				"email_verified": false,
			},
		},
		{
			name: "minimal claims - openid only",
			user: &types.User{
				ID:        12345,
				CreatedAt: now,
				Meta:      &types.UserMeta{},
			},
			scope:          "openid",
			expectedStatus: http.StatusOK,
			expectedClaims: map[string]interface{}{
				"sub": "12345",
			},
		},
		{
			name: "single name without family name",
			user: &types.User{
				ID:        12345,
				Name:      "Madonna",
				CreatedAt: now,
				Meta:      &types.UserMeta{},
			},
			scope:          "openid profile",
			expectedStatus: http.StatusOK,
			expectedClaims: map[string]interface{}{
				"sub":        "12345",
				"name":       "Madonna",
				"given_name": "Madonna",
			},
		},
		{
			name:           "no token",
			noToken:        true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "user not found",
			user: &types.User{
				ID:        99999,
				CreatedAt: now,
				Meta:      &types.UserMeta{},
			},
			scope:          "openid profile",
			userNotFound:   true,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/oauth2/userinfo", nil)
			rec := httptest.NewRecorder()

			// Set up context with JWT token
			ctx := req.Context()
			if !tc.noToken && tc.user != nil {
				// Create a properly signed token using the TokenIssuer
				signed, err := auth.TokenIssuer.Sign(
					auth.WithIdentity(auth.Authenticated(tc.user.ID)),
					auth.WithScope(tc.scope),
				)
				rq.NoError(err)

				// Parse the signed token to get jwt.Token
				token, err := jwt.Parse(signed)
				rq.NoError(err)

				ctx = jwtauth.NewContext(ctx, token, nil)
				req = req.WithContext(ctx)
			}

			// Set up mock user service
			userService := userServiceMocked{
				findByAny: func(ctx context.Context, any interface{}) (*types.User, error) {
					if tc.userNotFound {
						return nil, fmt.Errorf("user not found")
					}
					return tc.user, nil
				},
			}

			// Create handler
			h := AuthHandlers{
				Log:         zap.NewNop(),
				UserService: userService,
				Opt:         options.AuthOpt{BaseURL: "https://example.com"},
			}

			// Call handler
			h.oauth2Userinfo(rec, req)

			// Check status
			rq.Equal(tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				// Parse response
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				rq.NoError(err)

				// Check expected claims are present
				for key, expected := range tc.expectedClaims {
					rq.Contains(response, key, "missing claim: %s", key)
					rq.Equal(expected, response[key], "claim %s mismatch", key)
				}

				// Verify sub is always present
				rq.Contains(response, "sub")
			} else {
				// Check error response
				rq.Contains(rec.Header().Get("WWW-Authenticate"), "Bearer")
				rq.Equal("application/json", rec.Header().Get("Content-Type"))
			}
		})
	}
}

func Test_userinfoError(t *testing.T) {
	rq := require.New(t)

	rec := httptest.NewRecorder()
	h := AuthHandlers{Log: zap.NewNop()}

	h.userinfoError(rec, fmt.Errorf("test error"))

	rq.Equal(http.StatusUnauthorized, rec.Code)
	rq.Equal("application/json", rec.Header().Get("Content-Type"))
	rq.Contains(rec.Header().Get("WWW-Authenticate"), `Bearer error="invalid_token"`)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	rq.NoError(err)
	rq.Equal("invalid_token", response["error"])
	rq.Equal("test error", response["error_description"])
}
