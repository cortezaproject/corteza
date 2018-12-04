package {package}

{load warning.tpl}

import (
	"io"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx/types"
	"github.com/crusttech/crust/internal/rbac"
	"mime/multipart"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}
var _ = multipart.FileHeader{}
var _ = rbac.Operation

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
{elseif substr($param.type, 0, 2) === '[]' && isset($parsers[$param.type])}
	{self}.{param.name|expose} = {$parsers[$param.type]}({if $method === "post"}r.Form["{param.name}"]{else}urlQuery["{param.name}"]{/if})
{elseif $param.type === "*multipart.FileHeader"}
        if _, {self}.{param.name|expose}, err = r.FormFile("{$param.name}"); err != nil {
        	return errors.Wrap(err, "error procesing uploaded file")
        }
{elseif substr($param.type, 0, 2) !== '[]'}
        if val, ok := {method|strtolower}["{param.name}"]; ok {
{if $param.type === "types.JSONText"}
 		if {self}.{param.name|expose}, err = {$parsers[$param.type]}(val); err != nil {
 			return err
 		}
{else}
        {self}.{param.name|expose} = {if ($param.type !== "string")}{$parsers[$param.type]}(val){else}val{/if}
{/if}
	}{/if}
{/foreach}
{/foreach}{newline}
	return err
}

var _ RequestFiller = New{name|expose}{call.name|capitalize}()
{/foreach}
