package {package}

{load warning.tpl}

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/{project}/rest/server"
)

func MountRoutes(r chi.Router) {
{foreach $apis as $api}
	{api.interface|strtolower} := &server.{api.interface|capitalize}Handlers{{api.interface|capitalize}{ldelim}{rdelim}.New()}
{/foreach}
{foreach $apis as $api}
	r.Group(func (r chi.Router) {
			r.Use({api.interface|strtolower}.{api.interface}.Authenticator())
		r.Route("{api.path}", func(r chi.Router) {
{foreach $api.apis as $call}
			r.{eval echo capitalize(strtolower($call.method))}("{call.path}", {api.interface|strtolower}.{call.name|capitalize})
{/foreach}
		})
	})
{/foreach}
}