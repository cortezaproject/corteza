package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

var (
	HttpTokenVerifier func(http.Handler) http.Handler
)

// tokenClockSkew is how far the verifying clock may trail or lead the issuing
// one when a token's iat, nbf and exp are checked. A clock stepped back by time
// sync, or a second node a little behind, otherwise refuses a token it has just
// issued as "issued in the future".
const tokenClockSkew = 30 * time.Second

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
			token, err := verifyRequest(ja, r, jwtauth.TokenFromHeader, jwtauth.TokenFromQuery, jwtauth.TokenFromCookie)
			ctx := r.Context()

			if token != nil && err == nil {
				if err = TokenIssuer.Validate(ctx, token); err != nil {
					errors.ProperlyServeHTTP(w, r, err, false)
					return
				}
			}

			ctx = jwtauth.NewContext(ctx, token, err)
			ctx = SetIdentityToContext(ctx, IdentityFromToken(token))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// verifyRequest is jwtauth.VerifyRequest with tokenClockSkew allowed on the
// time claims.
func verifyRequest(ja *jwtauth.JWTAuth, r *http.Request, findTokenFns ...func(r *http.Request) string) (jwt.Token, error) {
	var tokenString string
	for _, fn := range findTokenFns {
		if tokenString = fn(r); tokenString != "" {
			break
		}
	}

	if tokenString == "" {
		return nil, jwtauth.ErrNoTokenFound
	}

	token, err := ja.Decode(tokenString)
	if err != nil {
		return token, jwtauth.ErrorReason(err)
	}

	if token == nil {
		return nil, jwtauth.ErrUnauthorized
	}

	if err = jwt.Validate(token, jwt.WithAcceptableSkew(tokenClockSkew)); err != nil {
		return token, jwtauth.ErrorReason(err)
	}

	return token, nil
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

// HttpAuthenticatedOnly rejects requests that carry no identity.
//
// HttpTokenValidator deliberately lets jwtauth.ErrNoTokenFound through, so it
// validates the scope of a token that is present rather than requiring one.
// That is right for the REST surface, where handlers reach RBAC and RBAC denies
// the anonymous role. It is not enough where a handler discloses something
// before any RBAC check runs — MCP tool discovery being the case in point: it
// lists every registered tool with its full description straight from the
// registry.
//
// Use this in addition to HttpTokenValidator, not instead of it: this one
// establishes that there is a caller, that one establishes the token is good
// for this surface.
func HttpAuthenticatedOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !GetIdentityFromContext(r.Context()).Valid() {
				errors.ProperlyServeHTTP(w, r, errUnauthorized(), false)
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
		// ErrNoTokenFound is passed through unchanged: HttpTokenValidator
		// identifies it to let an unauthenticated request continue to RBAC.
		if errors.Is(err, jwtauth.ErrNoTokenFound) {
			return err
		}

		// Anything else is a token that failed to parse or verify. Returned
		// raw it carries no Kind, so ProperlyServeHTTP rendered it as a 500 —
		// a malformed credential reported as a server fault. It is a rejected
		// credential, which is 401.
		return errUnauthorized()
	}

	if token == nil {
		return errUnauthorized()
	}

	if len(scope) > 0 && !CheckJwtScope(token, scope...) {
		return errUnauthorizedScope()
	}

	return
}
