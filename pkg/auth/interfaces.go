package auth

import (
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
		Encode(identity Identifiable) string
	}

	TokenHandler interface {
		TokenEncoder
		HttpVerifier() func(http.Handler) http.Handler
		HttpAuthenticator() func(http.Handler) http.Handler
	}

	Signer interface {
		Sign(userID uint64, pp ...interface{}) string
		Verify(signature string, userID uint64, pp ...interface{}) bool
	}
)
