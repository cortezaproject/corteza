package {package}

{load warning.tpl}

import (
	"io"
	"net/http"
	"encoding/json"
	"mime/multipart"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

{if !empty($imports)}{foreach ($imports as $import)}
	{import}{EOL}{/foreach}{/if}
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

{foreach $calls as $call}
// {name} {call.name} request parameters
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

func ({self} *{name|expose}{call.name|capitalize}) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode({self})

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

{eval $parseForm = "ParseForm()"}
{foreach $call.parameters as $method => $params}{foreach $params as $param}{if $param.type === "*multipart.FileHeader"}{eval $parseForm = "ParseMultipartForm(32 << 20)"}{/if}{/foreach}{/foreach}
	if err = r.{$parseForm}; err != nil {
		return err
	}

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
{elseif (substr($param.type, 0, 2) === '[]' || substr($param.type, -3) === "Set") && isset($parsers[$param.type])}
	{self}.{param.name|expose} = {$parsers[$param.type]}({if $method === "post"}r.Form["{param.name}"]{else}urlQuery["{param.name}"]{/if})
{elseif $param.type === "*multipart.FileHeader"}
	if _, {self}.{param.name|expose}, err = r.FormFile("{$param.name}"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}
{elseif substr($param.type, 0, 2) !== '[]' && substr($param.type, -3) !== 'Set'}
	if val, ok := {method|strtolower}["{param.name}"]; ok {
{if substr($parsers[$param.type], -7) === 'WithErr'}
		if {self}.{param.name|expose}, err = {$parsers[$param.type]}(val); err != nil {
			return err
		}
{else}
		{self}.{param.name|expose} = {if ($param.type !== "string")}{$parsers[$param.type]}(val){else}val{/if}{EOL}
{/if}
	}{/if}
{/foreach}
{/foreach}{newline}
	return err
}

var _ RequestFiller = New{name|expose}{call.name|capitalize}()
{/foreach}
