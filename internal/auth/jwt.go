package auth

import (
	"net/http"
	"strconv"
	"strings"
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
)

var (
	DefaultJwtHandler TokenHandler
)

func SetupDefault(secret string, expiry int) {
	// Use JWT secret for hmac signer for now
	DefaultSigner = HmacSigner(secret)
	DefaultJwtHandler, _ = JWT(secret, int64(expiry))

}

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
	claims := jwt.MapClaims{
		"userID": strconv.FormatUint(identity.Identity(), 10),
		"exp":    time.Now().Add(time.Duration(t.expiry) * time.Minute).Unix(),
	}

	if rr := identity.Roles(); len(rr) > 0 {
		var memberOf string
		for _, r := range identity.Roles() {
			memberOf = memberOf + " " + strconv.FormatUint(r, 10)
		}

		claims["memberOf"] = memberOf[1:] // trim leading space
	}

	_, jwt, _ := t.tokenAuth.Encode(claims)
	return jwt
}

// Authenticator converts JWT claims into Identity and stores it into context
func (t *token) Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwt, claims, err := jwtauth.FromContext(r.Context())

			// When token is present, expect no errors and valid claims!
			if jwt != nil {
				if err != nil {
					// But if token is present, the shouldn't be an error
					resputil.JSON(w, err)
					return
				}

				identity := &Identity{}
				if userID, ok := claims["userID"].(string); ok && len(userID) >= 0 {
					identity.id, _ = strconv.ParseUint(userID, 10, 64)
				}

				if memberOf, ok := claims["memberOf"].(string); ok && len(memberOf) >= 0 {
					ss := strings.Split(memberOf, " ")
					identity.memberOf = make([]uint64, len(ss))
					for i, s := range ss {
						identity.memberOf[i], _ = strconv.ParseUint(s, 10, 64)
					}
				}

				r = r.WithContext(SetIdentityToContext(r.Context(), identity))
			}

			next.ServeHTTP(w, r)
		})
	}
}
