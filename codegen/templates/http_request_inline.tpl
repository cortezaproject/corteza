package {package}

{load warning.tpl}

import (
	"io"
	"strings"

	"net/http"
	"encoding/json"
	"mime/multipart"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

{if !empty($imports)}{foreach ($imports as $import)}
	{import}{EOL}{/foreach}{/if}
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

{foreach $calls as $call}
// {name|expose} {call.name} request parameters
type {name|expose}{call.name|capitalize} struct {
{foreach $call.parameters as $params}
{foreach $params as $method => $param}
	{param.name|expose} {param.type}{if $param.type === "uint64" || $param.type === "[]uint64"} `json:",string"`{/if}{newline}
{/foreach}
{/foreach}
}

func New{name|expose}{call.name|capitalize}() *{name|expose}{call.name|capitalize} {
	return &{name|expose}{call.name|capitalize}{}
}

func (r {name|expose}{call.name|capitalize}) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}
{foreach $call.parameters as $method => $params}{foreach $params as $param}{if $param.sensitive}
	out["{param.name}"] = "*masked*sensitive*data*"
{elseif $param.type === "*multipart.FileHeader"}
	out["{param.name}.size"] = r.{param.name|expose}.Size
	out["{param.name}.filename"] = r.{param.name|expose}.Filename
{else}
	out["{param.name}"] = r.{param.name|expose}
{/if}{/foreach}{/foreach}{newline}
	return out
}

func (r *{name|expose}{call.name|capitalize}) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

{eval $parseForm = "ParseForm()"}
{foreach $call.parameters as $method => $params}{foreach $params as $param}{if $param.type === "*multipart.FileHeader"}{eval $parseForm = "ParseMultipartForm(32 << 20)"}{/if}{/foreach}{/foreach}
	if err = req.{$parseForm}; err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}
{foreach $call.parameters as $method => $params}
{foreach $params as $param}
{if strtolower($method) === "path"}
	r.{param.name|expose} = {if ($param.type !== "string")}{$parsers[$param.type]}({/if}chi.URLParam(req, "{param.name}"){if ($param.type !== "string")}){/if}
{elseif (substr($param.type, 0, 2) === '[]' || substr($param.type, -3) === "Set") && isset($parsers[$param.type])}
	{if strtolower($method) === "post"}
	if val, ok := req.Form["{param.name}"]; ok {
		r.{param.name|expose} = {$parsers[$param.type]}(val)
	}
	{elseif strtolower($method) === "get"}
	if val, ok := urlQuery["{param.name}[]"]; ok {
		r.{param.name|expose} = {$parsers[$param.type]}(val)
	} else if val, ok = urlQuery["{param.name}"]; ok {
		r.{param.name|expose} = {$parsers[$param.type]}(val)
	}
	{/if}

{elseif $param.type === "*multipart.FileHeader"}
	if _, r.{param.name|expose}, err = req.FormFile("{$param.name}"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}
{elseif substr($param.type, 0, 2) !== '[]' && substr($param.type, -3) !== 'Set'}
	if val, ok := {method|strtolower}["{param.name}"]; ok {
{if substr($parsers[$param.type], -7) === 'WithErr'}
		if r.{param.name|expose}, err = {$parsers[$param.type]}(val); err != nil {
			return err
		}
{else}		r.{param.name|expose} = {if ($param.type !== "string")}{if isset($parsers[$param.type])}{$parsers[$param.type]}{else}{$param.type}{/if}(val){else}val{/if}{EOL}
{/if}
	}{/if}
{/foreach}
{/foreach}{newline}
	return err
}

var _ RequestFiller = New{name|expose}{call.name|capitalize}()
{/foreach}
