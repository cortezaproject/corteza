package dal

import (
	"context"
	"fmt"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/pkg/gvalfnc"
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
	globalGvalConverter converterGval
)

var (
	refToGvalExp = map[string]*exprHandlerGval{
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
		"xor": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("(%s != %s)", args[0], args[1])
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

		// - strings
		"concat": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("concat(%s)", strings.Join(args, ", "))
			},
		},

		// @todo implement; the commented versions are not good enough
		// "like": {
		// 	Handler: func(args ...string) string {
		// 		// @todo better regex construction
		// 		nn := strings.Replace(strings.Trim(args[1], "\""), "%", ".*", -1)
		// 		nn = strings.Replace(nn, "_", ".[1]", -1)

		// 		return fmt.Sprintf("%s =~ ^%s$", args[0], nn)
		// 	},
		// },
		// "nlike": {
		// 	Handler: func(args ...string) string {
		// 		// @todo better regex construction
		// 		nn := strings.Replace(args[1], "%", ".*", -1)
		// 		nn = strings.Replace(nn, "_", ".[1]", -1)

		// 		return fmt.Sprintf("!(%s =~ ^%s$)", args[0], nn)
		// 	},
		// },

		// "is": {
		// 	Handler: func(args ...string) string {
		// 		return fmt.Sprintf("%s == %s", args[0], args[1])
		// 	}
		// }

		// - filtering
		"now": {
			Handler: func(args ...string) string {
				return "now()"
			},
		},
		"quarter": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("quarter(%s)", args[0])
			},
		},
		"year": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("year(%s)", args[0])
			},
		},
		"month": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("month(%s)", args[0])
			},
		},
		"date": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("date(%s)", args[0])
			},
		},
		"date_format": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("strftime(%s, %s)", args[0], args[1])
			},
		},

		// generic stuff
		"null": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("isNil(%s)", args[0])
			},
		},
		"nnull": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("!isNil(%s)", args[0])
			},
		},
		"exists": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("!isNil(%s)", args[0])
			},
		},

		// - typecast
		"float": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("float(%s)", args[0])
			},
		},
		"int": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("int(%s)", args[0])
			},
		},
		"string": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("string(%s)", args[0])
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
	if globalGvalConverter.parser == nil {
		globalGvalConverter = converterGval{
			parser: newQlParser(),
		}
	}

	return globalGvalConverter
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

	out.eval, err = newGval(expr)
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

	out.eval, err = newGval(expr)
	return
}

// newGval initializes a new gval evaluatable for the provided expression
//
// The eval. includes a small subset of supported expr. functions which may be
// used in the pipeline.
//
// @note the subset is limited to simplify the (eventual) offloading to the DB.
//       At some point, more functions will be supported, and the ones which can't
//       be offloaded will be performed in some exec. step.
func newGval(expr string) (gval.Evaluable, error) {
	return gval.Full(
		// Extra functions we'll need
		// @note don't bring in all of the expr. pkg functions as we'll need to
		//       support these on the DB as well
		gval.Function("now", gvalfnc.Now),
		gval.Function("quarter", gvalfnc.Quarter),
		gval.Function("year", gvalfnc.Year),
		gval.Function("month", gvalfnc.Month),
		gval.Function("strftime", gvalfnc.StrfTime),
		gval.Function("date", gvalfnc.Date),
		gval.Function("isNil", gvalfnc.IsNil),
		gval.Function("float", gvalfnc.CastFloat),
		gval.Function("int", gvalfnc.CastInt),
		gval.Function("string", gvalfnc.CastString),
		gval.Function("concat", gvalfnc.ConcatStrings),
	).NewEvaluable(expr)
}

func newQlParser() *ql.Parser {
	pp := ql.NewParser()
	pp.OnIdent = func(ident ql.Ident) (ql.Ident, error) {
		ident.Value = NormalizeAttrNames(ident.Value)
		return ident, nil
	}

	return pp
}

func (e *runnerGval) Test(ctx context.Context, rows any) (bool, error) {
	o, err := e.eval(ctx, rows)
	if err != nil {
		return false, err
	}

	v, ok := o.(bool)
	return ok && v, nil
}

func (e *runnerGval) Eval(ctx context.Context, rows any) (any, error) {
	o, err := e.eval(ctx, rows)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// Parse parses the QL expression into QL ASTNodes
func (c converterGval) Parse(expr string) (*ql.ASTNode, error) {
	return c.parser.Parse(expr)
}

// Convert converts the given nodes into a GVal expression
// @todo add more validation so we can potentially omit exec. error checks
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
	r := strings.ToLower(n.Ref)
	if refToGvalExp[r] == nil {
		return "", fmt.Errorf("unknown ref %q", n.Ref)
	}
	return refToGvalExp[r].Handler(args...), nil
}
