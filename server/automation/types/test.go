package types

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	// Used for input evaluation
	Test struct {
		// Expression to evaluate over the input variables; results will be set to scope under variable Name
		Expr string `json:"expr,omitempty"`

		eval expr.Evaluable

		// Error to be   if test fails
		Error string `json:"error"`
	}

	TestSet []*Test
)

func NewTest(expr, error string) (t *Test, err error) {
	return &Test{Expr: expr, Error: error}, nil
}

func (t Test) GetExpr() string              { return t.Expr }
func (t *Test) SetEval(eval expr.Evaluable) { t.eval = eval }
func (t Test) Eval(ctx context.Context, scope *expr.Vars) (interface{}, error) {
	return t.eval.Eval(ctx, scope)
}
func (t Test) Test(ctx context.Context, scope *expr.Vars) (bool, error) {
	return t.eval.Test(ctx, scope)
}

func (set TestSet) Validate(ctx context.Context, scope *expr.Vars) (TestSet, error) {
	vres := make(TestSet, 0, len(set))

	for _, t := range set {
		r, err := t.eval.Test(ctx, scope)
		if err != nil {
			return nil, err
		}

		if !r {
			vres = append(vres, &Test{Error: t.Error})
		}
	}

	return vres, nil
}

func (set TestSet) Test(ctx context.Context, scope *expr.Vars) (bool, error) {
	return set.TestAll(ctx, scope)
}

func (set TestSet) TestAll(ctx context.Context, scope *expr.Vars) (bool, error) {
	for _, t := range set {
		r, err := t.Test(ctx, scope)
		if err != nil || !r {
			return false, err
		}
	}

	return true, nil
}

// Returns true on first true
func (set TestSet) TestAny(ctx context.Context, scope *expr.Vars) (bool, error) {
	for _, t := range set {
		r, err := t.Test(ctx, scope)
		if err != nil {
			return false, err
		}

		if r {
			return true, nil
		}
	}

	return false, nil
}
