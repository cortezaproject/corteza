package auth

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type (
	Identifiable interface {
		Identity() uint64
		Roles() []uint64
		Valid() bool
		String() string
	}

	TokenEncoder interface {
		Encode(identity Identifiable, scope ...string) string
	}

	TokenHandler interface {
		TokenEncoder
		Authenticate(token string) (jwt.MapClaims, error)
		HttpVerifier() func(http.Handler) http.Handler
		HttpAuthenticator() func(http.Handler) http.Handler
	}

	Signer interface {
		Sign(userID uint64, pp ...interface{}) string
		Verify(signature string, userID uint64, pp ...interface{}) bool
	}
)
