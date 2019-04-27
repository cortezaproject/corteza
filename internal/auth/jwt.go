package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"
)

type (
	token struct {
		expiry    int64
		tokenAuth *jwtauth.JWTAuth
	}

	jwtSettingsGetter interface {
		GetGlobalString(name string) (out string, err error)
	}
)

var (
	DefaultJwtHandler TokenHandler
)

func JWT(secret string, expiry int64) (jwt *token, err error) {
	if len(secret) == 0 {
		return nil, errors.New("JWT secret missing")
	}

	jwt = &token{
		expiry:    expiry,
		tokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}

	return jwt, nil
}

// Verifies JWT and stores it into context
func (t *token) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(t.tokenAuth)
}

func (t *token) Encode(identity Identifiable) string {
	claims := jwt.StandardClaims{}
	claims.Subject = strconv.FormatUint(identity.Identity(), 10)
	claims.ExpiresAt = time.Now().Add(time.Duration(t.expiry) * time.Minute).Unix()

	_, jwt, _ := t.tokenAuth.Encode(claims)
	return jwt
}

// Extracts and authenticates JWT from context, validates claims
func (t *token) Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			if identityId, err := getIdentityClaimFromContext(ctx); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				// Request validated, identity confirmed
				r = r.WithContext(SetIdentityToContext(ctx, NewIdentity(identityId)))
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}
