package {package}

{load warning.tpl}

import (
	"github.com/go-chi/chi"
	"net/http"
)

func ({self}h *{name}Handlers)MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func (r chi.Router) {
			r.Use(middlewares...)
		r.Route("{api.path}", func(r chi.Router) {
{foreach $api.apis as $call}
			r.{eval echo capitalize(strtolower($call.method))}("{call.path}", {self}h.{call.name|capitalize})
{/foreach}
		})
	})
}
