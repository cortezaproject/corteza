package service

import (
	"context"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	// argsProc is a helper to process arguments
	//
	// when initialised it evaluates expressions using given scope
	// it provides methods to extract values from the result
	argsProc struct {
		expr   types.ExprSet
		result *expr.Vars
	}
)

func processArguments(ctx context.Context, expr []*types.Expr, scope *expr.Vars) (p *argsProc, err error) {
	p = &argsProc{expr: expr}

	if p.result, err = p.expr.Eval(ctx, scope); err != nil {
		return nil, err
	}

	return
}

func (p *argsProc) bool(name string, val *bool) bool     { return getStepArg(p, name, val) }
func (p *argsProc) string(name string, val *string) bool { return getStepArg(p, name, val) }
func (p *argsProc) uint64(name string, val *uint64) bool { return getStepArg(p, name, val) }

// vars handles assigning of expr.Vars
//
// not as straightforward as plain types
func (p *argsProc) vars(name string, val *expr.Vars) bool {
	// make an auxiliary variable to hold the results q
	aux := make(map[string]expr.TypedValue)

	// extract it
	if !getStepArg(p, name, &aux) {
		return false
	}

	if val == nil {
		panic("initialize Vars before calling the function")
	}

	vars, _ := expr.NewVars(aux)
	*val = *vars
	return true
}

// getStepArg is a helper to extract step argument from result
func getStepArg[T any](p *argsProc, arg string, val *T) bool {
	var (
		aux any
	)

	if p.result.Has(arg) {
		aux = expr.Must(p.result.Select(arg)).Get()
	} else if exp := p.expr.GetByTarget(arg); exp != nil {
		aux = p.expr.GetByTarget(arg).Value.(T)
	} else {
		return false
	}

	if conv, is := aux.(T); is {
		*val = conv
		return true
	}

	return false
}
