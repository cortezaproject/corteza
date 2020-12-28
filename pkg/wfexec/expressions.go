package wfexec

import (
	"context"
	"fmt"
	"github.com/PaesslerAG/gval"
	"strings"
)

type (
	Expression struct {
		// where to assign the evaluated Expression
		name string

		// Expression
		expr string

		// Expression, ready to be executed
		eval gval.Evaluable
	}

	Expressions struct {
		lang gval.Language
		set  []*Expression
	}
)

func NewExpression(lang gval.Language, dst, expr string) (e *Expression, err error) {
	e = &Expression{name: dst, expr: expr}

	if e.eval, err = lang.NewEvaluable(expr); err != nil {
		return nil, fmt.Errorf("can not parse Expression %s: %w", expr, err)
	}

	return e, nil
}

func NewExpressions(lang gval.Language, ee ...*Expression) *Expressions {
	return &Expressions{
		lang: lang,
		set:  ee,
	}
}

func (ee *Expressions) Set(dst, expr string) error {
	var (
		e, err = NewExpression(ee.lang, dst, expr)
	)

	if err != nil {
		return err
	}

	for i := range ee.set {
		if ee.set[i].name == dst {
			ee.set[i] = e
			return nil
		}
	}

	ee.set = append(ee.set, e)
	return nil
}

func (ee *Expressions) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	if result, err := ee.Eval(ctx, r.Scope); err != nil {
		return nil, err
	} else {
		return r.Scope.Merge(result), nil
	}
}

func (ee *Expressions) Eval(ctx context.Context, in Variables) (Variables, error) {
	var (
		err error
		// Copy/create scope
		scope = Variables.Merge(in)
		out   = Variables{}
	)

	for _, e := range ee.set {
		if strings.Contains(e.name, ".") {
			// handle property setting
			return nil, fmt.Errorf("dot/prop setting not supported at the moment")
		}

		if scope[e.name], err = e.eval(ctx, scope); err != nil {
			return nil, fmt.Errorf("could not evaluate %q for %q: %w", e.expr, e.name, err)
		}

		out[e.name] = scope[e.name]
	}

	return out, nil
}
