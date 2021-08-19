package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	pe "github.com/cortezaproject/corteza-server/pkg/errors"
)

type (
	redirection struct {
		types.FilterMeta

		location *url.URL
		status   int

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

func NewDefaultJsonResponse() (e *defaultJsonResponse) {
	e = &defaultJsonResponse{}

	e.Name = "defaultJsonResponse"
	e.Label = "Default JSON response"
	e.Kind = types.PostFilter

	return
}

func (h defaultJsonResponse) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h defaultJsonResponse) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (f *defaultJsonResponse) Merge(params []byte) (h types.Handler, err error) {
	return f, err
}

func (h defaultJsonResponse) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusAccepted)

		if _, err := rw.Write([]byte(`{}`)); err != nil {
			return pe.Internal("could not write to body: (%v)", err)
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
