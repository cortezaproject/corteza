package {package}

{load warning.tpl}

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/{project}/rest/request"
)

// Internal API interface
type {name|expose}API interface {
{foreach $calls as $call}
	{call.name|capitalize}(context.Context, *request.{name|expose}{call.name|capitalize}) (interface{}, error)
{/foreach}
}

// HTTP API interface
type {name|expose} struct {
{foreach $calls as $call}
	{call.name|capitalize} func(http.ResponseWriter, *http.Request)
{/foreach}
}

func New{name|expose}({self}h {name|expose}API) *{name|expose} {
	return &{name|expose}{ldelim}{newline}
{foreach $calls as $call}
		{call.name|capitalize}: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.New{name|capitalize}{call.name|capitalize}()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := {self}h.{call.name|capitalize}(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
					case func(http.ResponseWriter, *http.Request):
						fn(w, r)
						return
				}
				resputil.JSON(w, value)
				return
			}
		},
{/foreach}
	{rdelim}
}

func ({self}h *{name|expose})MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func (r chi.Router) {
		r.Use(middlewares...)
{foreach $api.apis as $call}
		r.{eval echo capitalize(strtolower($call.method))}("{api.path}{call.path}", {self}h.{call.name|capitalize})
{/foreach}
	})
}
