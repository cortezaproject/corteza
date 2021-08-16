package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
)

type (
	redirection struct {
		types.FilterMeta
		params struct {
			HTTPStatus int    `json:"status,string"`
			Location   string `json:"location"`
		}
	}

	// support for arbitrary response
	// obfuscation
	customResponse struct {
		types.FilterMeta
		params struct {
			Source string `json:"source"`
		}
	}

	defaultJsonResponse struct {
		types.FilterMeta
	}

	errorHandler struct {
		name string
		args []string
	}
)

func NewRedirection() (e *redirection) {
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

func (h redirection) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h redirection) Type() types.FilterKind {
	return h.Kind
}

func (h redirection) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h redirection) Weight() int {
	return h.Wgt
}

func (f *redirection) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)
	return f, err
}

func (h redirection) Exec(ctx context.Context, scope *types.Scp) error {
	loc, err := url.ParseRequestURI(h.params.Location)

	if err != nil {
		return fmt.Errorf("could not redirect: %s", err)
	}

	status := h.params.HTTPStatus

	if !checkStatus("redirect", status) {
		return fmt.Errorf("could not redirect: wrong status %d", status)
	}

	http.Redirect(scope.Writer(), scope.Request(), loc.String(), status)

	return nil
}

func NewErrorHandler(name string, args []string) (e *errorHandler) {
	e = &errorHandler{
		name: name,
		args: args,
	}

	return
}

func (pp errorHandler) Exec(ctx context.Context, scope *types.Scp, err error) {
	type (
		responseHelper struct {
			ErrResponse struct {
				Msg string `json:"msg"`
			} `json:"error"`
		}
	)

	resp := responseHelper{
		ErrResponse: struct {
			Msg string "json:\"msg\""
		}{
			Msg: err.Error(),
		},
	}

	// set http status code
	scope.Writer().WriteHeader(http.StatusInternalServerError)

	// set body
	json.NewEncoder(scope.Writer()).Encode(resp)

}

func NewDefaultJsonResponse() (e *defaultJsonResponse) {
	e = &defaultJsonResponse{}

	e.Name = "defaultJsonResponse"
	e.Label = "Default JSON response"
	e.Kind = types.PostFilter

	return
}

func (h defaultJsonResponse) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h defaultJsonResponse) Type() types.FilterKind {
	return h.Kind
}

func (h defaultJsonResponse) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h defaultJsonResponse) Weight() int {
	return h.Wgt
}

func (f *defaultJsonResponse) Merge(params []byte) (h types.Handler, err error) {
	return f, err
}

func (h defaultJsonResponse) Exec(ctx context.Context, scope *types.Scp) (err error) {
	scope.Writer().Header().Set("Content-Type", "application/json")
	scope.Writer().WriteHeader(http.StatusAccepted)

	_, err = scope.Writer().Write([]byte(`{}`))

	return
}

func checkStatus(typ string, status int) bool {
	switch typ {
	case "redirect":
		return status >= 300 && status <= 399
	default:
		return true
	}
}
