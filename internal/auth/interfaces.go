package auth

import (
	"net/http"
)

type (
	Identifiable interface {
		Identity() uint64
		Roles() []uint64
		Valid() bool
	}

	TokenEncoder interface {
		Encode(identity Identifiable) string
	}

	TokenHandler interface {
		Encode(identity Identifiable) string
		Verifier() func(http.Handler) http.Handler
		Authenticator() func(http.Handler) http.Handler
	}

	Signer interface {
		Sign(userID uint64, pp ...interface{}) string
		Verify(signature string, userID uint64, pp ...interface{}) bool
	}
)
