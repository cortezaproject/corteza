package {package}

{load warning.tpl}

import (
	"net/http"

	"github.com/go-chi/chi"
)

func ({self}h *{name|expose}Handlers)MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func (r chi.Router) {
			r.Use(middlewares...)
		r.Route("{api.path}", func(r chi.Router) {
{foreach $api.apis as $call}
			r.{eval echo capitalize(strtolower($call.method))}("{call.path}", {self}h.{call.name|capitalize})
{/foreach}
		})
	})
}
