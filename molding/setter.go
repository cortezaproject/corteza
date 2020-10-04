package molding

import (
	"context"
	"fmt"
	"github.com/PaesslerAG/gval"
	"strings"
)

type (
	expr struct {
		// target for the result of the evaluated expression
		target string

		// expression
		source string

		// expression, ready to be executed
		eval gval.Evaluable
	}

	setActivity struct {
		next Node

		// List of expressions that will be evaluated
		// with given scope
		expressions []*expr
	}
)

var (
	_ Executor = &setActivity{}
)

func Expr(dst, source string) *expr {
	return &expr{target: dst, source: source}
}

func NewSetActivity(next Node, ee ...*expr) (*setActivity, error) {
	var (
		err error
		set = &setActivity{
			next:        next,
			expressions: ee,
		}

		lang = gval.Full()
	)

	for _, e := range set.expressions {
		if e.eval, err = lang.NewEvaluable(e.source); err != nil {
			return nil, fmt.Errorf("can no parse %s for %s: %w", e.source, e.target, err)
		}
	}

	return set, nil
}

func (s setActivity) NodeRef() string { return "setter" }
func (s setActivity) Next() Node      { return s.next }
func (s *setActivity) SetNext(n Node) { s.next = n }

func (s setActivity) Exec(ctx context.Context, params Variables) (Variables, error) {
	var (
		err error

		// Create scope from params
		//
		// We'll use it for evaluation of preconfigured expressions on the setter
		// and as a container for each value that comes out of that evaluation
		scope = params.Merge()
	)

	for _, e := range s.expressions {
		if strings.Contains(e.target, ".") {
			// handle property setting
			return nil, fmt.Errorf("no support for prop setting ATM")
		}

		if scope[e.target], err = e.eval(ctx, scope); err != nil {
			return nil, fmt.Errorf("could not evaluate %s for %s: %w", e.source, e.target, err)
		}
	}

	return scope, nil
}
