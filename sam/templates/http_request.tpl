package {package}

import (
	"net/http"

	"github.com/pkg/errors"
)

{foreach $calls as $call}
// {name} {call.name} request parameters
type {name}{call.name|ucfirst}Request struct {
{foreach $call.parameters as $params}
{foreach $params as $method => $param}
	{param.name} {param.type}{newline}
{/foreach}
{/foreach}
}

func ({name}{call.name|ucfirst}Request) new() *{name}{call.name|ucfirst}Request {
	return &{name}{call.name|ucfirst}Request{}
}

func ({self} *{name}{call.name|ucfirst}Request) Fill(r *http.Request) error {
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
{if substr($param.type, 0, 2) !== '[]'}
	{self}.{param.name} = {if $param.type !== "string"}{$parsers[$param.type]}({method|strtolower}["{param.name}"]){else}{method|strtolower}["{param.name}"]{/if}{newline}
{/if}
{/foreach}
{/foreach}
	return errors.New("Not implemented: {name}{call.name|ucfirst}Request.Fill")
}

var _ RequestFiller = {name}{call.name|ucfirst}Request{}.new()
{/foreach}
