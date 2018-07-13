package sam

import (
	"net/http"
)

var pass = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (*OrganisationHandlers) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*TeamHandlers) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*ChannelHandlers) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*MessageHandlers) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*UserHandlers) Authenticator() func(http.Handler) http.Handler {
	return pass
}

func (*WebsocketHandlers) Authenticator() func(http.Handler) http.Handler {
	return pass
}
