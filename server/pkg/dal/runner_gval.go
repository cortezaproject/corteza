package dal

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/gvalfnc"
	"github.com/cortezaproject/corteza/server/pkg/ql"
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

		// @todo temporarily added to support the current aggregate attribute type calculation
		OutType        Type
		OutTypeUnknown bool
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
			OutType: &TypeBoolean{},
		},

		// - bool
		"and": {
			Handler: func(args ...string) string {
				return strings.Join(args, " && ")
			},
			OutType: &TypeBoolean{},
		},
		"or": {
			Handler: func(args ...string) string {
				return strings.Join(args, " || ")
			},
			OutType: &TypeBoolean{},
		},
		"xor": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("(%s != %s)", args[0], args[1])
			},
			OutType: &TypeBoolean{},
		},

		// - comp.
		"eq": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s == %s", args[0], args[1])
			},
			OutType: &TypeBoolean{},
		},
		"ne": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s != %s", args[0], args[1])
			},
			OutType: &TypeBoolean{},
		},
		"lt": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s < %s", args[0], args[1])
			},
			OutType: &TypeBoolean{},
		},
		"le": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s <= %s", args[0], args[1])
			},
			OutType: &TypeBoolean{},
		},
		"gt": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s > %s", args[0], args[1])
			},
			OutType: &TypeBoolean{},
		},
		"ge": {
			// Handler: makeGenericCompHandler(">="),
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s >= %s", args[0], args[1])
			},
			OutType: &TypeBoolean{},
		},

		// - math
		"add": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s + %s", args[0], args[1])
			},
			OutTypeUnknown: true,
		},
		"sub": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s - %s", args[0], args[1])
			},
			OutTypeUnknown: true,
		},
		"mult": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s * %s", args[0], args[1])
			},
			OutTypeUnknown: true,
		},
		"div": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("%s / %s", args[0], args[1])
			},
			OutTypeUnknown: true,
		},

		"in": {
			Handler: func(args ...string) string {
				// The arguments must be reversed!!
				return fmt.Sprintf("has(%s, %s)", args[1], args[0])
			},
			OutTypeUnknown: true,
		},
		"nin": {
			Handler: func(args ...string) string {
				// The arguments must be reversed!!
				return fmt.Sprintf("!has(%s, %s)", args[1], args[0])
			},
			OutTypeUnknown: true,
		},

		// - strings
		"concat": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("concat(%s)", strings.Join(args, ", "))
			},
			OutType: &TypeText{},
		},
		"instr": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("instr(%s, %s)", args[0], args[1])
			},
			OutType: &TypeNumber{},
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
			OutType: &TypeTimestamp{},
		},
		"quarter": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("quarter(%s)", args[0])
			},
			OutType: &TypeNumber{},
		},
		"year": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("year(%s)", args[0])
			},
			OutType: &TypeNumber{},
		},
		"month": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("month(%s)", args[0])
			},
			OutType: &TypeNumber{},
		},
		"date": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("date(%s)", args[0])
			},
			OutType: &TypeDate{},
		},
		"day": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("day(%s)", args[0])
			},
			OutType: &TypeNumber{},
		},
		"week": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("week(%s)", args[0])
			},
			OutType: &TypeNumber{},
		},
		"date_format": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("strftime(%s, %s)", args[0], args[1])
			},
			OutType: &TypeText{},
		},

		// generic stuff
		"null": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("isNil(%s)", args[0])
			},
			OutType: &TypeBoolean{},
		},
		"nnull": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("!isNil(%s)", args[0])
			},
			OutType: &TypeBoolean{},
		},
		"exists": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("!isNil(%s)", args[0])
			},
			OutType: &TypeBoolean{},
		},

		// - typecast
		"float": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("float(%s)", args[0])
			},
			OutType: &TypeNumber{},
		},
		"int": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("int(%s)", args[0])
			},
			OutType: &TypeNumber{},
		},
		"string": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("string(%s)", args[0])
			},
			OutType: &TypeText{},
		},

		"group": {
			Handler: func(args ...string) string {
				return fmt.Sprintf("(%s)", args[0])
			},
			OutTypeUnknown: true,
		},
	}
)

// newConverterGval initializes a new gval exp. converter
func newConverterGval(ii ...ql.IdentHandler) converterGval {
	if globalGvalConverter.parser == nil {
		globalGvalConverter = converterGval{
			parser: newQlParser(append(
				ii,
				// This should always happen for gval expressions
				func(i ql.Ident) (ql.Ident, error) {
					i.Value = wrapNestedGvalIdent(i.Value)
					return i, nil
				},
			)...),
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
func newGval(e string) (gval.Evaluable, error) {
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
		gval.Function("day", gvalfnc.Day),
		gval.Function("isNil", gvalfnc.IsNil),
		gval.Function("float", gvalfnc.CastFloat),
		gval.Function("int", gvalfnc.CastInt),
		gval.Function("string", gvalfnc.CastString),
		gval.Function("concat", gvalfnc.ConcatStrings),
		gval.Function("has", arrHas),
	).NewEvaluable(e)
}

func newQlParser(onIdent ...ql.IdentHandler) *ql.Parser {
	pp := ql.NewParser()
	pp.OnIdent = func(ident ql.Ident) (_ ql.Ident, err error) {
		ident.Value = NormalizeAttrNames(ident.Value)

		for _, oi := range onIdent {
			ident, err = oi(ident)
			if err != nil {
				return
			}
		}

		return ident, err
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

// arrHas is a helper to assure the expr.Has always gets an array (or map) as
// the first argument
//
// @todo this is needed because how the ValueGetters returns multi-value fields so
//       an edge case where a field would have [a] but here, it would be presented
//       as a.
//       This would become obsolete when we address the actual issue.
func arrHas(arr interface{}, vv ...interface{}) (b bool, err error) {
	arr = expr.UntypedValue(arr)

	if isMap(arr) {
		return expr.Has(arr, vv...)
	}

	var (
		c = reflect.ValueOf(arr)
	)

	switch c.Kind() {
	case reflect.Slice:
		return expr.Has(arr, vv...)

	default:
		return expr.Has([]any{arr}, vv...)
	}
}

func isMap(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}
