package auth

import (
	"net/http"
)

type (
	Identifiable interface {
		Identity() uint64
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
)
