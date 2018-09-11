package auth

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/crusttech/crust/auth/types"
	"github.com/go-chi/jwtauth"
	"github.com/titpetric/factory/resputil"
)

type jwt struct {
	tokenAuth *jwtauth.JWTAuth
}

func JWT() (*jwt, error) {
	if err := flags.Validate(); err != nil {
		return nil, err
	}

	jwt := &jwt{tokenAuth: jwtauth.New("HS256", []byte(flags.jwt.Secret), nil)}

	if flags.jwt.DebugToken {
		log.Println("DEBUG JWT TOKEN:", jwt.Encode(NewIdentity(1)))
	}

	return jwt, nil
}

// Verifies JWT and stores it into context
func (t *jwt) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(t.tokenAuth)
}

func (t *jwt) Encode(identity types.Identifiable) string {
	// @todo Set expiry
	claims := jwtauth.Claims{}
	claims.Set("sub", strconv.FormatUint(identity.Identity(), 10))
	claims.SetExpiryIn(time.Duration(flags.jwt.Expiry) * time.Minute)

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
