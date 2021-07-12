package apigw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	verifierQueryParam struct {
		functionMeta
		params struct {
			Expr string `json:"expr"`
		}
	}

	verifierOrigin struct {
		functionMeta
		params struct {
			Expr string `json:"expr"`
		}
	}
)

func NewVerifierOrigin() (v *verifierOrigin) {
	v = &verifierOrigin{}

	v.Step = 0
	v.Name = "verifierOrigin"
	v.Label = "Origin verifier"
	v.Kind = FunctionKindVerifier

	v.Args = []*functionMetaArg{
		{
			Type:    "expr",
			Label:   "expr",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h verifierOrigin) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h verifierOrigin) Meta() functionMeta {
	return h.functionMeta
}

func (v *verifierOrigin) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&v.params)
	return v, err
}

func (h verifierOrigin) Exec(ctx context.Context, scope *scp) error {
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

func NewVerifierQueryParam() (v *verifierQueryParam) {
	v = &verifierQueryParam{}

	v.Step = 0
	v.Name = "verifierQueryParam"
	v.Label = "Query parameters verifier"
	v.Kind = "verifier"

	v.Args = []*functionMetaArg{
		{
			Type:    "expr",
			Label:   "expr",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h verifierQueryParam) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h verifierQueryParam) Meta() functionMeta {
	return h.functionMeta
}

func (v *verifierQueryParam) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&v.params)
	return v, err
}

func (h verifierQueryParam) Exec(ctx context.Context, scope *scp) error {
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
	// }

	// testing
	scope.Request().Header.Add(fmt.Sprintf("step_%d", h.Step), h.Name)

	return nil
}
