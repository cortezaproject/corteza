package auth

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	token struct {
		// Expiration time in minutes
		expiry    time.Duration
		tokenAuth *jwtauth.JWTAuth
	}
)

var (
	DefaultJwtHandler TokenHandler
)

func SetupDefault(secret string, expiry time.Duration) {
	// Use JWT secret for hmac signer for now
	DefaultSigner = HmacSigner(secret)
	DefaultJwtHandler, _ = JWT(secret, expiry)

}

func JWT(secret string, expiry time.Duration) (tkn *token, err error) {
	if len(secret) == 0 {
		return nil, errors.New("JWT secret missing")
	}

	tkn = &token{
		expiry:    expiry,
		tokenAuth: jwtauth.New(jwt.SigningMethodHS512.Alg(), []byte(secret), nil),
	}

	return tkn, nil
}

// Verifies JWT and stores it into context
func (t *token) HttpVerifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(t.tokenAuth)
}

func (t *token) Encode(i Identifiable) string {

	roles := ""
	for _, r := range i.Roles() {
		roles += fmt.Sprintf(" %d", r)
	}

	_, tkn, _ := t.tokenAuth.Encode(jwt.MapClaims{
		"sub": i.String(),
		"exp": time.Now().Add(t.expiry).Unix(),

		// when possible, extend this with the client
		"aud": "0", // we should be more strict here and provide a client!

		// at some point, we'll need to do this a bit more carefully;
		// for now, we just add all necessary scopes so that users of this (cli, corredor)
		// can create a useful JWT
		"scope": "profile api",

		"roles": strings.TrimSpace(roles),
	})

	return tkn
}

// HttpAuthenticator converts JWT claims into Identity and stores it into context
func (t *token) HttpAuthenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tkn, claims, err := jwtauth.FromContext(ctx)

			// When token is present, expect no errors and valid claims!
			if tkn != nil {
				if err != nil {
					// But if token is present, the shouldn't be an error
					api.Send(w, r, err)
					return
				}

				ctx = SetIdentityToContext(ctx, claimsToIdentity(claims))
				ctx = context.WithValue(ctx, scopeCtxKey{}, claims["scope"])

				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// decodes sub & roles claims into identity
func claimsToIdentity(c jwt.MapClaims) (i *Identity) {
	var (
		aux       interface{}
		ok        bool
		id, roles string
	)

	if aux, ok = c["sub"]; !ok {
		return
	}

	i = &Identity{}
	if id, ok = aux.(string); ok {
		i.id, _ = strconv.ParseUint(id, 10, 64)
	}

	if i.id == 0 {
		// pointless to decode roles if id is 0
		return nil
	}

	if aux, ok = c["roles"]; !ok {
		return
	}

	if roles, ok = aux.(string); !ok {
		return
	}

	for _, role := range strings.Split(roles, " ") {
		if roleID, _ := strconv.ParseUint(role, 10, 64); roleID != 0 {
			i.memberOf = append(i.memberOf, roleID)
		}
	}

	return
}
