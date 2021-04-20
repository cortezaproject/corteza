package websocket

import (
	"github.com/go-chi/chi"
	"net/http"
)

func middlewareAllowedAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if !service.DefaultAccessControl.CanAccess(r.Context()) {
		//	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}

func (ws Websocket) MountRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/websocket", func(r chi.Router) {
			//r.Use(middlewareAllowedAccess)
			r.Get("/", ws.Open)
		})
	})
}
