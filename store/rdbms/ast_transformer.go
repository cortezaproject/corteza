package rdbms

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/qlng"
)

type (
	astTransformer struct {
		placeholders bool
		root         *qlng.ASTNode
		custom       []ASTFormatterFn
	}

	// FormattedASTArgs is a temporary storage for already handled nested arguments
	FormattedASTArgs struct {
		S    string
		Args []interface{}

		ResultType string
	}
	ASTFormatterFn func(n *qlng.ASTNode, aa ...FormattedASTArgs) (bool, string, []interface{}, error)
	HandlerSig     func(aa ...FormattedASTArgs) (string, []interface{}, error)

	exprHandler struct {
		Args   argSet
		RArgs  bool
		Result *arg

		Handler HandlerSig
	}
	argSet []*arg
	arg    struct {
		Required bool
		Type     string
	}
)

var (
	bracketHandler = makeGenericBracketHandler("(", ")")

	// @todo IS and IS NOT; should this be calculated with eq operators?
	sqlExprRegistry = map[string]exprHandler{
		// operators
		// - bool
		"and": {
			Args:    collectParams(true, "Boolean"),
			RArgs:   true,
			Result:  wrapRes("Boolean"),
			Handler: makeGenericBoolHandler("AND"),
		},
		"or": {
			Args:    collectParams(true, "Boolean"),
			RArgs:   true,
			Result:  wrapRes("Boolean"),
			Handler: makeGenericBoolHandler("OR"),
		},
		"xor": {
			Args:    collectParams(true, "Boolean"),
			RArgs:   true,
			Result:  wrapRes("Boolean"),
			Handler: makeGenericBoolHandler("XOR"),
		},

		// - comp.
		"eq": {
			Args:    collectParams(true, "Any", "Any"),
			Result:  wrapRes("Boolean"),
			Handler: makeGenericCompHandler("="),
		},
		"ne": {
			Args:    collectParams(true, "Any", "Any"),
			Result:  wrapRes("Boolean"),
			Handler: makeGenericCompHandler("!="),
		},
		"lt": {
			Args:    collectParams(true, "Any", "Any"),
			Result:  wrapRes("Boolean"),
			Handler: makeGenericCompHandler("<"),
		},
		"le": {
			Args:    collectParams(true, "Any", "Any"),
			Result:  wrapRes("Boolean"),
			Handler: makeGenericCompHandler("<="),
		},
		"gt": {
			Args:    collectParams(true, "Any", "Any"),
			Result:  wrapRes("Boolean"),
			Handler: makeGenericCompHandler(">"),
		},
		"ge": {
			Args:    collectParams(true, "Any", "Any"),
			Result:  wrapRes("Boolean"),
			Handler: makeGenericCompHandler(">="),
		},

		// - math
		"add": {
			Args:    collectParams(true, "Number", "Number"),
			RArgs:   true,
			Result:  wrapRes("Number"),
			Handler: makeGenericCompHandler("+"),
		},
		"sub": {
			Args:    collectParams(true, "Number", "Number"),
			RArgs:   true,
			Result:  wrapRes("Number"),
			Handler: makeGenericCompHandler("-"),
		},
		"mult": {
			Args:    collectParams(true, "Number", "Number"),
			RArgs:   true,
			Result:  wrapRes("Number"),
			Handler: makeGenericCompHandler("*"),
		},
		"div": {
			Args:    collectParams(true, "Number", "Number"),
			RArgs:   true,
			Result:  wrapRes("Number"),
			Handler: makeGenericCompHandler("/"),
		},

		// @todo better negation?
		"like": {
			Args:    collectParams(true, "String", "String"),
			Result:  wrapRes("Boolean"),
			Handler: makeGenericCompHandler("LIKE"),
		},
		"nlike": {
			Args:    collectParams(true, "String", "String"),
			Result:  wrapRes("Boolean"),
			Handler: makeGenericCompHandler("NOT LIKE"),
		},

		// functions
		// - aggregation
		"count": {
			Args:    collectParams(false, "Any"),
			Result:  wrapRes("Number"),
			Handler: makeGenericAggFncHandler("COUNT"),
		},
		"sum": {
			Args:    collectParams(true, "Any"),
			Result:  wrapRes("Number"),
			Handler: makeGenericAggFncHandler("SUM"),
		},
		"max": {
			Args:    collectParams(true, "Any"),
			Result:  wrapRes("Number"),
			Handler: makeGenericAggFncHandler("MAX"),
		},
		"min": {
			Args:    collectParams(true, "Any"),
			Result:  wrapRes("Number"),
			Handler: makeGenericAggFncHandler("MIN"),
		},
		"avg": {
			Args:    collectParams(true, "Any"),
			Result:  wrapRes("Number"),
			Handler: makeGenericAggFncHandler("AVG"),
		},

		// - filtering
		"now": {
			Result:  wrapRes("DateTime"),
			Handler: makeGenericFilterFncHandler("NOW"),
		},
		"quarter": {
			Args:    collectParams(true, "DateTime"),
			Result:  wrapRes("Number"),
			Handler: makeGenericFilterFncHandler("QUARTER"),
		},
		"year": {
			Args:    collectParams(true, "DateTime"),
			Result:  wrapRes("Number"),
			Handler: makeGenericFilterFncHandler("YEAR"),
		},
		"date": {
			Args:    collectParams(true, "DateTime"),
			Result:  wrapRes("Number"),
			Handler: makeGenericFilterFncHandler("DATE"),
		},

		// generic stuff
		"null": {
			Result: wrapRes("Null"),
			Handler: func(aa ...FormattedASTArgs) (string, []interface{}, error) {
				return "NULL", nil, nil
			},
		},

		// - typecast
		"float": {
			Args:    collectParams(true, "Any"),
			Result:  wrapRes("Float"),
			Handler: makeGenericTypecastHandler("FLOAT"),
		},
		"int": {
			Args:    collectParams(true, "Any"),
			Result:  wrapRes("Float"),
			Handler: makeGenericTypecastHandler("INT"),
		},
	}
)

func newASTFormatter(n *qlng.ASTNode, custom ...ASTFormatterFn) *astTransformer {
	return &astTransformer{
		placeholders: true,
		root:         n,
		custom:       custom,
	}
}

func (t *astTransformer) SetPlaceholder(use bool) {
	t.placeholders = use
}

// ToSql conforms the struct to squirrel allowing trivial RDBMS use
func (t *astTransformer) ToSql() (string, []interface{}, error) {
	return t.toSql(t.root)
}

func (t *astTransformer) toSql(n *qlng.ASTNode) (string, []interface{}, error) {
	// Leaf edge-cases
	switch {
	case n.Symbol != "":
		return t.handleSymbol(n)
	case n.Value != nil:
		return t.handleValue(n)
	}

	// Process arguments for the op.
	args := make([]FormattedASTArgs, len(n.Args))
	for i, a := range n.Args {
		s, pp, err := t.toSql(a)
		if err != nil {
			return "", nil, err
		}

		args[i] = FormattedASTArgs{
			S:    s,
			Args: pp,
		}
	}

	// Custom handlers take precedence
	for _, c := range t.custom {
		if c == nil {
			continue
		}

		if ok, s, args, err := c(n, args...); ok {
			return s, args, err
		}
	}

	if n.Ref == "group" {
		return bracketHandler(args...)
	}

	// Default handlers
	if e, ok := sqlExprRegistry[n.Ref]; !ok {
		return "", nil, fmt.Errorf("unknown expression: handler not defined: %s", n.Ref)
	} else {
		s, args, err := e.Handler(args...)
		return s, args, err
	}
}

func wrapIntfs(vv ...interface{}) []interface{} {
	return vv
}

func mergeIntfs(aa ...[]interface{}) (out []interface{}) {
	out = make([]interface{}, 0, 10)
	for _, a := range aa {
		out = append(out, a...)
	}

	return out
}

// Analyze analyzes the AST and returns the resulting type and any errors
func (t *astTransformer) Analyze(symbolIndex map[string]string) (string, error) {
	return t.analyze(t.root, symbolIndex)
}

func (t *astTransformer) analyze(n *qlng.ASTNode, symbolIndex map[string]string) (string, error) {
	// Leaf edge-cases
	// - symbol
	if n.Symbol != "" {
		sy, ok := symbolIndex[n.Symbol]
		if !ok {
			return "", fmt.Errorf("unknown symbol %s", n.Symbol)
		}
		return sy, nil
	}
	// - plain value
	if n.Value != nil {
		return n.Value.V.Type(), nil
	}

	// expr. validity
	args := make([]FormattedASTArgs, len(n.Args))
	for i, a := range n.Args {
		t, err := t.analyze(a, symbolIndex)
		if err != nil {
			return "", err
		}

		args[i] = FormattedASTArgs{
			ResultType: t,
		}
	}

	// groups proxy the type from their inner ops
	if n.Ref == "group" {
		if len(args) != 1 {
			return "", fmt.Errorf("a group must have a single root operation")
		}
		return args[0].ResultType, nil
	}

	// custom handlers should have no affect on the input/output of the handler
	e, ok := sqlExprRegistry[n.Ref]
	if !ok {
		return "", fmt.Errorf("unknown expression: handler not defined: %s", n.Ref)
	}

	// - check arg count
	pc, opc := countParams(e)
	if !e.RArgs {
		if len(args) < pc-opc || len(args) > pc {
			if opc > 0 {
				return "", fmt.Errorf("%s: expecting %d + %d arguments, got %d", n.Ref, pc-opc, opc, len(args))
			}
			return "", fmt.Errorf("%s: expecting %d arguments, got %d", n.Ref, pc, len(args))
		}
	} else {
		if len(args)%pc != 0 {
			return "", fmt.Errorf("%s: expecting multiple of %d arguments, got %d", n.Ref, len(e.Args), len(args))
		}
	}

	// - check arg types
	for i, p := range args {
		a := e.Args[i%len(e.Args)].Type
		b := p.ResultType

		if a == "Any" {
			continue
		}

		// Number allows all numeric types
		if a == "Number" && isNumber(b) {
			continue
		}

		if a == b {
			continue
		}

		return "", fmt.Errorf("argument type missmatch: expecting %s, got %s for argument %d", a, b, i)
	}

	return e.Result.Type, nil
}

// // // // // // // // // // // // // // // // // // // // // // // // //

func isNumber(t string) bool {
	return t == "Number" || t == "Integer" || t == "UnsignedInteger" || t == "Float"
}

func countParams(e exprHandler) (all int, opt int) {
	all = len(e.Args)

	for _, p := range e.Args {
		if !p.Required {
			opt++
		}
	}

	return
}

func collectParams(req bool, ss ...string) argSet {
	pp := make(argSet, len(ss))
	for i, t := range ss {
		pp[i] = wrapParam(req, t)
	}
	return pp
}

func wrapParam(req bool, t string) *arg {
	return &arg{
		Required: req,
		Type:     t,
	}
}

func wrapRes(t string) *arg {
	return wrapParam(false, t)
}
