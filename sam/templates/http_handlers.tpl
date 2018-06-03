package {package}

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

{foreach $calls as $call}
func ({self}h *{name}Handlers) {call.name|ucfirst}(w http.ResponseWriter, r *http.Request) {
	params := {name}{call.name|ucfirst}Request{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return {self}h.{name}.{call.name|ucfirst}(params) })
}
{/foreach}
