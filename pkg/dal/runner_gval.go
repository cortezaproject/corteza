package dal

import (
	"context"
	"fmt"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/pkg/ql"
)

type (
	runnerGval struct {
		eval gval.Evaluable
	}
	converterGval struct {
		parser *ql.Parser
	}
	exprHandlerGval struct {
		Handler func(...string) string
	}
)

var (
	refToGvalExp = map[string]*exprHandlerGval{
		// keywords
		"null": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s == null", args[0])
			},
		},
		"nnull": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s != null", args[0])
			},
		},

		// operators
		"not": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("(! %s)", args[0])
			},
		},

		// - bool
		"and": {
			Handler: func(args ...string) string {
				return strings.Join(args, " && ")
			},
		},
		"or": {
			Handler: func(args ...string) string {
				return strings.Join(args, " || ")
			},
		},

		// - comp.
		"eq": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s == %s", args[0], args[1])
			},
		},
		"ne": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s != %s", args[0], args[1])
			},
		},
		"lt": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s < %s", args[0], args[1])
			},
		},
		"le": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s <= %s", args[0], args[1])
			},
		},
		"gt": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s > %s", args[0], args[1])
			},
		},
		"ge": {
			//Handler: makeGenericCompHandler(">="),
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s >= %s", args[0], args[1])
			},
		},

		// - math
		"add": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s + %s", args[0], args[1])
			},
		},
		"sub": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s - %s", args[0], args[1])
			},
		},
		"mult": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s * %s", args[0], args[1])
			},
		},
		"div": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s / %s", args[0], args[1])
			},
		},

		"group": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("(%s)", args[0])
			},
		},
	}
)

// newConverterGval initializes a new gval exp. converter
func newConverterGval() converterGval {
	return converterGval{
		parser: ql.NewParser(),
	}
}

// newRunnerGval initializes a new gval exp. runner from the provided expression
func newRunnerGval(expr string) (out *runnerGval, err error) {
	out = &runnerGval{}

	c := newConverterGval()
	n, err := c.Parse(expr)
	if err != nil {
		return nil, err
	}

	expr, err = c.Convert(n)
	if err != nil {
		return nil, err
	}

	// @todo add all those extra functions the expr package uses?
	//       Potentially not ok since we allow a subset of operations (for compatability)
	out.eval, err = gval.Full().NewEvaluable(expr)
	return
}

// newRunnerGvalParsed initializes a new gval exp. runner from the pre-parsed expression
func newRunnerGvalParsed(n *ql.ASTNode) (out *runnerGval, err error) {
	out = &runnerGval{}
	c := newConverterGval()

	expr, err := c.Convert(n)
	if err != nil {
		return
	}

	// @todo add all those extra functions the expr package uses?
	//       Potentially not ok since we allow a subset of operations (for compatability)
	out.eval, err = gval.Full().NewEvaluable(expr)
	return
}

func (c converterGval) Parse(expr string) (*ql.ASTNode, error) {
	return c.parser.Parse(expr)
}

func (c converterGval) Convert(n *ql.ASTNode) (expr string, err error) {
	return c.convert(n)
}

func (c converterGval) convert(n *ql.ASTNode) (_ string, err error) {
	switch {
	case n.Symbol != "":
		return n.Symbol, nil
	case n.Value != nil:
		// @todo I don't think this is quite ok, but it works for now so it'll do for now
		if n.Value.V.Type() == "String" {
			return fmt.Sprintf("\"%v\"", n.Value.V.Get()), nil
		}
		return fmt.Sprintf("%v", n.Value.V.Get()), nil
	}

	args := make([]string, len(n.Args))
	for i, a := range n.Args {
		args[i], err = c.convert(a)
		if err != nil {
			return
		}
	}

	return c.refHandler(n, args...)
}

func (c converterGval) refHandler(n *ql.ASTNode, args ...string) (out string, err error) {
	if refToGvalExp[n.Ref] == nil {
		return "", fmt.Errorf("unknown ref %q", n.Ref)
	}
	return refToGvalExp[n.Ref].Handler(args...), nil
}

func (e *runnerGval) Test(ctx context.Context, rows any) bool {
	o, err := e.eval(ctx, rows)
	if err != nil {
		return false
	}

	v, ok := o.(bool)
	return ok && v
}

func (e *runnerGval) Eval(ctx context.Context, rows any) any {
	o, err := e.eval(ctx, rows)
	if err != nil {
		return nil
	}

	return o
}
