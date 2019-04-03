package {package}

{load warning.tpl}

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

{foreach $calls as $call}
func ({self}h *{name|expose}Handlers) {call.name|capitalize}(w http.ResponseWriter, r *http.Request) {
	params := {name|capitalize}{call.name|capitalize}Request{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return {self}h.{name|expose}.{call.name|capitalize}(r.Context(), params) })
}
{/foreach}
