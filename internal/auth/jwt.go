package auth

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/titpetric/factory/resputil"
)

type token struct {
	expiry       int64
	cookieDomain string
	tokenAuth    *jwtauth.JWTAuth
}

func JWT() (*token, error) {
	if err := flags.Validate(); err != nil {
		return nil, err
	}

	jwt := &token{
		expiry:       flags.jwt.Expiry,
		cookieDomain: flags.jwt.CookieDomain,
		tokenAuth:    jwtauth.New("HS256", []byte(flags.jwt.Secret), nil),
	}

	if flags.jwt.DebugToken {
		log.Println("DEBUG JWT TOKEN:", jwt.Encode(NewIdentity(1)))
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

// Extracts and authenticates JWT from context, validates claims
func (t *token) SetCookie(w http.ResponseWriter, r *http.Request, identity Identifiable) {
	cookie := &http.Cookie{
		Name:    "jwt",
		Expires: time.Now().Add(time.Duration(t.expiry) * time.Minute),
		Secure:  r.URL.Scheme == "https",
		Domain:  t.cookieDomain,
		Path:    "/",
	}

	if identity == nil {
		cookie.Expires = time.Unix(0, 0)
	} else {
		cookie.Value = t.Encode(identity)
	}

	http.SetCookie(w, cookie)
}
