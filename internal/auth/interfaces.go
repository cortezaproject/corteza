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
		SetCookie(w http.ResponseWriter, r *http.Request, identity Identifiable)
	}
)
