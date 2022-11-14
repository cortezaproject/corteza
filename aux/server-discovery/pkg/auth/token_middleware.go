package auth

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

var (
	HttpTokenVerifier func(http.Handler) http.Handler
)

// verifier returns a jwt verification middleware
//
// Tasks
// 1. picks token from header, query or cookie
// 2. extracts identity (if any) and adds it into request context
//
// In there is no token
func verifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := jwtauth.VerifyRequest(ja, r, jwtauth.TokenFromHeader, jwtauth.TokenFromQuery, jwtauth.TokenFromCookie)
			ctx := r.Context()

			ctx = jwtauth.NewContext(ctx, token, err)
			ctx = SetIdentityToContext(ctx, IdentityFromToken(token))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TokenVerifierMiddlewareWithSecretSigner returns HTTP handler with simple jwa.HS512 + secret verifier
//
// This should be 1:1 with token issuer!
func TokenVerifierMiddlewareWithSecretSigner(secret string) (_ func(http.Handler) http.Handler, err error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("JWK missing")
	}

	var key jwk.Key
	if key, err = jwk.New([]byte(secret)); err != nil {
		return nil, fmt.Errorf("could not parse JWK: %w", err)
	}

	return verifier(jwtauth.New(jwa.HS512.String(), key, nil)), nil
}

// HttpTokenValidator checks if there is a token with identity and matching scope claim
//
// Empty scope defaults to "api"!
func HttpTokenValidator(scope ...string) func(http.Handler) http.Handler {
	if len(scope) == 0 {
		// ensure that scope is not empty
		scope = []string{"api"}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := verifyToken(r.Context(), scope...)
			if err != nil && !errors.Is(err, jwtauth.ErrNoTokenFound) {
				errors.ProperlyServeHTTP(w, r, err, false)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// pulls token from context and validates scope & access-token
func verifyToken(ctx context.Context, scope ...string) (err error) {
	var token jwt.Token
	if token, _, err = jwtauth.FromContext(ctx); err != nil {
		return
	}

	if token == nil {
		return fmt.Errorf("unauthorized")
	}

	// if len(scope) > 0 && !CheckJwtScope(token, scope...) {
	// 	return errUnauthorizedScope()
	// }

	return
}
