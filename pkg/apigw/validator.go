package apigw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	validatorHeader struct {
		functionMeta
		params struct {
			Expr string `json:"expr"`
		}
	}
)

func NewValidatorHeader() (v *validatorHeader) {
	v = &validatorHeader{}

	v.Step = 3
	v.Name = "validatorHeader"
	v.Label = "Header validator"
	v.Kind = FunctionKindValidator

	v.Args = []*functionMetaArg{
		{
			Type:    "expr",
			Label:   "expr",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h validatorHeader) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h validatorHeader) Meta() functionMeta {
	return h.functionMeta
}

func (v *validatorHeader) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&v.params)
	return v, err
}

func (h validatorHeader) Exec(ctx context.Context, scope *scp) error {
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
