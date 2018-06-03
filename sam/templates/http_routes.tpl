package {package}

import (
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
{foreach $apis as $api}
	{api.entrypoint|strtolower} := {api.entrypoint|ucfirst}Handlers{}.new()
{/foreach}
{foreach $apis as $api}
	r.Route("{api.path}", func(r chi.Router) {
{foreach $api.apis as $call}
		r.{eval echo ucfirst(strtolower($call.method))}("{call.path}", {api.entrypoint|strtolower}.{call.name|ucfirst})
{/foreach}
	})
{/foreach}
}