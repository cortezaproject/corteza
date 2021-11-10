package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type (
	Identifiable interface {
		Identity() uint64
		Roles() []uint64
		Valid() bool
		String() string
	}

	TokenGenerator interface {
		Generate(ctx context.Context, identity Identifiable) (string, error)
	}

	TokenHandler interface {
		TokenGenerator
		Authenticate(token string) (jwt.MapClaims, error)
		HttpVerifier() func(http.Handler) http.Handler
		HttpAuthenticator() func(http.Handler) http.Handler
	}

	Signer interface {
		Sign(userID uint64, pp ...interface{}) string
		Verify(signature string, userID uint64, pp ...interface{}) bool
	}
)
