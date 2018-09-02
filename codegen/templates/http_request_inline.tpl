package {package}

{load warning.tpl}

import (
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi"
)

var _ = chi.URLParam

{foreach $calls as $call}
// {name} {call.name} request parameters
type {name|expose}{call.name|capitalize} struct {
{foreach $call.parameters as $params}
{foreach $params as $method => $param}
	{param.name|expose} {param.type}{newline}
{/foreach}
{/foreach}
}

func New{name|expose}{call.name|capitalize}() *{name|expose}{call.name|capitalize} {
	return &{name|expose}{call.name|capitalize}{}
}

func ({self} *{name|expose}{call.name|capitalize}) Fill(r *http.Request) error {
	json.NewDecoder(r.Body).Decode({self})

	r.ParseForm()
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}
{foreach $call.parameters as $method => $params}
{foreach $params as $param}
{if strtolower($method) === "path"}
	{self}.{param.name|expose} = {if ($param.type !== "string")}{$parsers[$param.type]}({/if}chi.URLParam(r, "{param.name}"){if ($param.type !== "string")}){/if}
{elseif substr($param.type, 0, 2) !== '[]'}
        if val, ok := {method|strtolower}["{param.name}"]; ok {
		{self}.{param.name|expose} = {if ($param.type !== "string")}{$parsers[$param.type]}(val){else}val{/if}
	}{/if}
{/foreach}
{/foreach}{newline}
	return nil
}

var _ RequestFiller = New{name|expose}{call.name|capitalize}()
{/foreach}
