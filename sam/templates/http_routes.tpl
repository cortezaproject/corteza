package {package}

import (
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
{foreach $apis as $api}
	{api.interface|strtolower} := {api.interface|capitalize}Handlers{}.new()
{/foreach}
{foreach $apis as $api}
	r.Route("{api.path}", func(r chi.Router) {
{foreach $api.apis as $call}
		r.{eval echo capitalize(strtolower($call.method))}("{call.path}", {api.interface|strtolower}.{call.name|capitalize})
{/foreach}
	})
{/foreach}
}