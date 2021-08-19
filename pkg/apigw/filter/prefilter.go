package filter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	pe "github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	header struct {
		types.FilterMeta
		eval   expr.Evaluable
		params struct {
			Expr string `json:"expr"`
		}
	}

	queryParam struct {
		types.FilterMeta
		eval   expr.Evaluable
		params struct {
			Expr string `json:"expr"`
		}
	}

	origin struct {
		types.FilterMeta
		eval   expr.Evaluable
		params struct {
			Expr string `json:"expr"`
		}
	}
)

func NewHeader() (v *header) {
	v = &header{}

	v.Name = "header"
	v.Label = "Header"
	v.Kind = types.PreFilter

	v.Args = []*types.FilterMetaArg{
		{
			Type:    "expr",
			Label:   "expr",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h header) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h header) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (v *header) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&v.params)

	if err != nil {
		return nil, err
	}

	parser := expr.NewParser()
	v.eval, err = parser.Parse(v.params.Expr)

	if err != nil {
		return nil, fmt.Errorf("could not validate origin parameters: %s", err)
	}

	return v, err
}

func (h header) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		var (
			ctx = r.Context()
		)

		vv := map[string]interface{}{}
		headers := r.Header

		for k, v := range headers {
			vv[k] = v[0]
		}

		// get the request data and put it into vars
		out, err := expr.NewVars(vv)

		if err != nil {
			return pe.Internal("could not validate headers: (%v) (%s)", err, h.params.Expr)
		}

		b, err := h.eval.Test(ctx, out)

		if err != nil {
			return pe.InvalidData("could not validate headers: (%v) (%s)", err, h.params.Expr)
		}

		if !b {
			return pe.InvalidData("could not validate headers: (%v) (%s)", errors.New("validation failed"), h.params.Expr)
		}

		return nil
	}
}

func NewQueryParam() (v *queryParam) {
	v = &queryParam{}

	v.Name = "queryParam"
	v.Label = "Query parameters"
	v.Kind = types.PreFilter

	v.Args = []*types.FilterMetaArg{
		{
			Type:    "expr",
			Label:   "expr",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h queryParam) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h queryParam) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (v *queryParam) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&v.params)

	if err != nil {
		return nil, err
	}

	parser := expr.NewParser()
	v.eval, err = parser.Parse(v.params.Expr)

	if err != nil {
		return nil, fmt.Errorf("could not validate query parameters: %s", err)
	}

	return v, err
}

func (h *queryParam) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		var (
			ctx = r.Context()
		)

		vv := map[string]interface{}{}
		vals := r.URL.Query()

		for k, v := range vals {
			vv[k] = v[0]
		}

		// get the request data and put it into vars
		out, err := expr.NewVars(vv)

		if err != nil {
			return pe.Internal("could not validate query parameters: (%v) (%s)", err, h.params.Expr)
		}

		b, err := h.eval.Test(ctx, out)

		if err != nil {
			return pe.InvalidData("could not validate query parameters: (%v) (%s)", err, h.params.Expr)
		}

		if !b {
			return pe.InvalidData("could not validate query parameters: (%v) (%s)", errors.New("validation failed"), h.params.Expr)
		}

		return nil
	}
}
