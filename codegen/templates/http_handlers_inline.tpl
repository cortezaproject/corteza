package {package}

{load warning.tpl}

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/{project}/rest/request"
	"github.com/crusttech/crust/internal/logger"
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

func New{name|expose}(h {name|expose}API) *{name|expose} {
	return &{name|expose}{ldelim}{newline}
{foreach $calls as $call}
		{call.name|capitalize}: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.New{name|capitalize}{call.name|capitalize}()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("{name|expose}.{call.name|capitalize}", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := h.{call.name|capitalize}(r.Context(), params); err != nil {
				logger.LogControllerError("{name|expose}.{call.name|capitalize}", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("{name|expose}.{call.name|capitalize}", r, params.Auditable())
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

func (h {name|expose})MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func (r chi.Router) {
		r.Use(middlewares...)
{foreach $api.apis as $call}
		r.{eval echo capitalize(strtolower($call.method))}("{api.path}{call.path}", h.{call.name|capitalize})
{/foreach}
	})
}
