package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	header struct {
		types.FilterMeta
		params struct {
			Expr string `json:"expr"`
		}
	}

	queryParam struct {
		types.FilterMeta
		params struct {
			Expr string `json:"expr"`
		}
	}

	origin struct {
		types.FilterMeta
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
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h header) Type() types.FilterKind {
	return h.Kind
}

func (h header) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h header) Weight() int {
	return h.Wgt
}

func (v *header) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&v.params)
	return v, err
}

func (h header) Exec(ctx context.Context, scope *types.Scp) error {
	vv := map[string]interface{}{}
	headers := scope.Request().Header

	for k, v := range headers {
		// sanitize header keys?
		vv[k] = v[0]
	}

	// get the request data and put it into vars
	out, err := expr.NewVars(vv)

	if err != nil {
		return err
	}

	pp := expr.NewParser()
	tt, err := pp.Parse(h.params.Expr)

	if err != nil {
		return fmt.Errorf("could not parse matching expression: %s", err)
	}

	b, err := tt.Test(ctx, out)

	if err != nil {
		return fmt.Errorf("could not validate headers: %s", err)
	}

	if !b {
		return fmt.Errorf("could not validate headers")
	}

	return nil
}

func NewOrigin() (v *origin) {
	v = &origin{}

	v.Name = "origin"
	v.Label = "Origin"
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

func (h origin) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h origin) Type() types.FilterKind {
	return h.Kind
}

func (h origin) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h origin) Weight() int {
	return h.Wgt
}

func (v *origin) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&v.params)
	return v, err
}

func (h origin) Exec(ctx context.Context, scope *types.Scp) error {
	vv := map[string]interface{}{
		"origin": scope.Request().Header.Get("Origin"),
	}

	// get the request data and put it into vars
	out, err := expr.NewVars(vv)

	if err != nil {
		return err
	}
	// spew.Dump("OUT", out)
	pp := expr.NewParser()
	tt, err := pp.Parse(h.params.Expr)

	if err != nil {
		return fmt.Errorf("could not parse matching expression: %s", err)
	}

	b, err := tt.Test(ctx, out)

	if err != nil {
		return fmt.Errorf("could not validate origin: %s", err)
	}

	if !b {
		return fmt.Errorf("could not validate origin")
	}

	return nil
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
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h queryParam) Type() types.FilterKind {
	return h.Kind
}

func (h queryParam) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h queryParam) Weight() int {
	return h.Wgt
}

func (v *queryParam) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&v.params)
	return v, err
}

func (h queryParam) Exec(ctx context.Context, scope *types.Scp) error {
	vv := map[string]interface{}{}
	vals := scope.Request().URL.Query()

	for k, v := range vals {
		vv[k] = v[0]
	}

	// get the request data and put it into vars
	out, err := expr.NewVars(vv)

	if err != nil {
		return err
	}

	pp := expr.NewParser()
	tt, err := pp.Parse(h.params.Expr)

	if err != nil {
		return fmt.Errorf("could not parse matching expression: %s", err)
	}

	b, err := tt.Test(ctx, out)

	if err != nil {
		return fmt.Errorf("could not validate query params: %s", err)
	}

	if !b {
		return fmt.Errorf("could not validate query params")
	}

	return nil
}
