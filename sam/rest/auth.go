package rest

import (
	"net/http"
)

var pass = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (*Organisation) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*Team) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*Message) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*Channel) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*User) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*Websocket) Authenticator() func(http.Handler) http.Handler {
	return pass
}
