package auth

import (
	"github.com/go-chi/jwtauth"
	"github.com/titpetric/factory/resputil"
	"net/http"
	"strconv"
	"time"
)

type jwt struct {
	tokenAuth *jwtauth.JWTAuth
}

func JWT() (*jwt, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	jwt := &jwt{tokenAuth: jwtauth.New("HS256", []byte(config.jwtSecret), nil)}

	return jwt, nil
}

// Verifies JWT and stores it into context
func (t *jwt) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(t.tokenAuth)
}

func (t *jwt) Encode(identity Identifiable) string {
	// @todo Set expiry
	claims := jwtauth.Claims{}
	claims.Set("sub", strconv.FormatUint(identity.Identity(), 10))
	claims.SetExpiryIn(time.Duration(config.jwtExpiry) * time.Minute)

	_, jwt, _ := t.tokenAuth.Encode(claims)
	return jwt
}

// Extracts and authenticates JWT from context, validates claims
func (t *jwt) Authenticator() func(http.Handler) http.Handler {
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
