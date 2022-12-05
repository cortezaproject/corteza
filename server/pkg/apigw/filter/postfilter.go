package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	agctx "github.com/cortezaproject/corteza/server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	errors "github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/options"
)

type (
	typesRegistry interface {
		Type(ref string) expr.Type
	}
	redirection struct {
		types.FilterMeta

		location *url.URL
		status   int

		params struct {
			HTTPStatus int    `json:"status,string"`
			Location   string `json:"location"`
		}
	}

	response struct {
		types.FilterMeta

		reg typesRegistry

		params struct {
			Header http.Header `json:"header"`

			Exp *atypes.Expr `json:"input"`

			Evaluable expr.Evaluable
		}
	}

	defaultJsonResponse struct {
		types.FilterMeta
	}
)

func NewRedirection(opts options.ApigwOpt) (e *redirection) {
	e = &redirection{}

	e.Name = "redirection"
	e.Label = "Redirection"
	e.Kind = types.PostFilter

	e.Args = []*types.FilterMetaArg{
		{
			Type:    "status",
			Label:   "status",
			Options: map[string]interface{}{},
		},
		{
			Type:    "text",
			Label:   "location",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h redirection) New(opts options.ApigwOpt) types.Handler {
	return NewRedirection(opts)
}

func (h redirection) Enabled() bool {
	return true
}

func (h redirection) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h redirection) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h redirection) Weight() int {
	return h.Wgt
}

func (h *redirection) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h.params)

	if err != nil {
		return h, err
	}

	loc, err := url.ParseRequestURI(h.params.Location)

	if err != nil {
		return nil, fmt.Errorf("could not validate parameters, invalid URL: %s", err)
	}

	if !checkStatus("redirect", h.params.HTTPStatus) {
		return nil, fmt.Errorf("could not validate parameters, wrong status %d", h.params.HTTPStatus)
	}

	h.location = loc
	h.status = h.params.HTTPStatus

	return h, err
}

func (h redirection) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		http.Redirect(rw, r, h.location.String(), h.status)
		return nil
	}
}

func NewDefaultJsonResponse(opts options.ApigwOpt) (e *defaultJsonResponse) {
	e = &defaultJsonResponse{}

	e.Name = "defaultJsonResponse"
	e.Label = "Default JSON response"
	e.Kind = types.PostFilter

	return
}

func (j defaultJsonResponse) New(opts options.ApigwOpt) types.Handler {
	return NewDefaultJsonResponse(opts)
}

func (j defaultJsonResponse) Enabled() bool {
	return true
}

func (j defaultJsonResponse) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", j.Name, j.Label)
}

func (j defaultJsonResponse) Meta() types.FilterMeta {
	return j.FilterMeta
}

func (j *defaultJsonResponse) Merge(params []byte) (h types.Handler, err error) {
	return j, err
}

func (j defaultJsonResponse) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusAccepted)

		if _, err := rw.Write([]byte(`{}`)); err != nil {
			return errors.Internal("could not write to body: %v", err)
		}

		return nil
	}
}

func checkStatus(typ string, status int) bool {
	switch typ {
	case "redirect":
		return status >= 300 && status <= 399
	default:
		return true
	}
}

func NewResponse(opts options.ApigwOpt, reg typesRegistry) (e *response) {
	e = &response{}

	e.Name = "response"
	e.Label = "JSON response"
	e.Kind = types.PostFilter

	e.Args = []*types.FilterMetaArg{
		{
			Type:    "input",
			Label:   "input",
			Options: map[string]interface{}{},
		},
		{
			Type:    "header",
			Label:   "header",
			Options: map[string]interface{}{},
		},
	}

	e.reg = reg

	return
}

func (j response) New(opts options.ApigwOpt) types.Handler {
	return NewResponse(opts, j.reg)
}

func (j response) Enabled() bool {
	return true
}

func (j response) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", j.Name, j.Label)
}

func (j response) Meta() types.FilterMeta {
	return j.FilterMeta
}

func (j *response) Merge(params []byte) (h types.Handler, err error) {
	var (
		parser = expr.NewParser()
	)

	err = json.NewDecoder(bytes.NewBuffer(params)).Decode(&j.params)

	if err != nil {
		return j, err
	}

	j.params.Evaluable, err = parser.Parse(j.params.Exp.Expr)

	if err != nil {
		return j, fmt.Errorf("could not evaluate expression: %s", err)
	}

	j.params.Exp.SetEval(j.params.Evaluable)

	return j, err
}

func (j response) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) (err error) {
		var (
			ctx           = r.Context()
			scope         = agctx.ScopeFromContext(ctx)
			hasJsonHeader = false

			evald interface{}
		)

		in, err := expr.NewVars(scope.Dict())

		if err != nil {
			return errors.Internal("could not validate request data: %v", err)
		}

		// set type to the registered expression from
		// any of the already registered types
		j.params.Exp.SetType(func(name string) (e expr.Type, err error) {
			if name == "" {
				name = "Any"
			}

			if typ := j.reg.Type(name); typ != nil {
				return typ, nil
			} else {
				return nil, errors.Internal("unknown or unregistered type %s", name)
			}
		})

		evald, err = j.params.Exp.Eval(ctx, in)

		if err != nil {
			return
		}

		for h, v := range j.params.Header {
			for _, vv := range v {
				rw.Header().Add(h, vv)

				if h == "Content-Type" || h == "content-type" {
					hasJsonHeader = vv == "application/json"
				}
			}
		}

		if reflect.ValueOf(evald).Kind() != reflect.String && hasJsonHeader {
			err = (json.NewEncoder(rw)).Encode(expr.UntypedValue(evald))
			return
		}

		fmt.Fprintf(rw, "%v", evald)

		return
	}
}
