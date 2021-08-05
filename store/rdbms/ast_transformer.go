package rdbms

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/qlng"
)

type (
	// ASTFormatterFn function signature for custom AST formatting
	ASTFormatterFn func(n *qlng.ASTNode, aa ...FormattedASTArgs) (bool, string, []interface{}, error)
	HandlerSig     func(aa ...FormattedASTArgs) (string, []interface{}, error)

	astTransformer struct {
		root   *qlng.ASTNode
		custom []ASTFormatterFn
	}

	// FormattedASTArgs is a temporary storage for already handled nested arguments
	FormattedASTArgs struct {
		S    string
		Args []interface{}
	}
)

var (
	// @todo IS and IS NOT; should this be calculated with eq operators?
	sqlExprRegistry = map[string]HandlerSig{
		// operators
		// - bool
		"and": makeGenericBoolHandler("AND"),
		"or":  makeGenericBoolHandler("OR"),
		"xor": makeGenericBoolHandler("XOR"),

		// - comp.
		"eq":   makeGenericCompHandler("="),
		"ne":   makeGenericCompHandler("!="),
		"lt":   makeGenericCompHandler("<"),
		"le":   makeGenericCompHandler("<="),
		"gt":   makeGenericCompHandler(">"),
		"ge":   makeGenericCompHandler(">="),
		"add":  makeGenericCompHandler("+"),
		"sub":  makeGenericCompHandler("-"),
		"mult": makeGenericCompHandler("*"),
		"div":  makeGenericCompHandler("/"),

		// @todo better negation?
		"ptrn":  makeGenericCompHandler("LIKE"),
		"nptrn": makeGenericCompHandler("NOT LIKE"),

		// functions
		// - aggregation
		"count": makeGenericAggFncHandler("COUNT"),
		"sum":   makeGenericAggFncHandler("SUM"),
		"max":   makeGenericAggFncHandler("MAX"),
		"min":   makeGenericAggFncHandler("MIN"),
		"avg":   makeGenericAggFncHandler("AVG"),

		// - filtering
		"now":     makeGenericFilterFncHandler("NOW"),
		"quarter": makeGenericFilterFncHandler("QUARTER"),
		"year":    makeGenericFilterFncHandler("YEAR"),
		"date":    makeGenericFilterFncHandler("DATE"),

		// generic stuff
		"group": makeGenericBracketHandler("(", ")"),
		"null": func(aa ...FormattedASTArgs) (string, []interface{}, error) {
			return "NULL", nil, nil
		},
	}
)

func newASTFormatter(n *qlng.ASTNode, custom ...ASTFormatterFn) squirrel.Sqlizer {
	return &astTransformer{
		root:   n,
		custom: custom,
	}
}

// ToSql conforms the struct to squirrel allowing trivial RDBMS use
func (t *astTransformer) ToSql() (string, []interface{}, error) {
	return t.format(t.root)
}

func (t *astTransformer) format(n *qlng.ASTNode) (string, []interface{}, error) {
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
		s, pp, err := t.format(a)
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

	// Default handlers
	if e, ok := sqlExprRegistry[n.Ref]; !ok {
		return "", nil, fmt.Errorf("unknown expression: handler not defined: %s", n.Ref)
	} else {
		s, args, err := e(args...)
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

// // // // // // // // // // // // // // // // // // // // // // // // //
// Default handlers

func (t *astTransformer) handleSymbol(n *qlng.ASTNode) (string, []interface{}, error) {
	return n.Symbol, nil, nil
}

func (t *astTransformer) handleValue(n *qlng.ASTNode) (string, []interface{}, error) {
	return "?", wrapIntfs(n.Value.Value), nil
}

func makeGenericBoolHandler(op string) HandlerSig {
	return func(aa ...FormattedASTArgs) (out string, args []interface{}, err error) {
		outPts := make([]string, len(aa))
		args = make([]interface{}, 0, 10)
		for i, a := range aa {
			outPts[i] = a.S
			args = append(args, a.Args...)
		}

		out = strings.Join(outPts, " "+op+" ")
		return
	}
}

func makeGenericBracketHandler(bb ...string) HandlerSig {
	return func(aa ...FormattedASTArgs) (out string, args []interface{}, err error) {
		if len(aa) != 1 {
			err = fmt.Errorf("expecting 1 argument, got %d", len(aa))
			return
		}

		out = fmt.Sprintf("%s%s%s", bb[0], aa[0].S, bb[1])
		args = aa[0].Args
		return
	}
}

func makeGenericCompHandler(comp string) HandlerSig {
	return func(aa ...FormattedASTArgs) (out string, args []interface{}, err error) {
		if len(aa) != 2 {
			err = fmt.Errorf("expecting 2 arguments, got %d", len(aa))
			return
		}

		out = fmt.Sprintf("%s %s %s", aa[0].S, comp, aa[1].S)
		args = mergeIntfs(aa[0].Args, aa[1].Args)
		return
	}
}

func makeGenericAggFncHandler(fnc string) HandlerSig {
	return func(aa ...FormattedASTArgs) (out string, args []interface{}, err error) {
		if fnc == "count" && len(aa) == 0 {
			out = "COUNT(*)"
			return
		}

		if len(aa) != 1 {
			err = fmt.Errorf("expecting 1 argument, got %d", len(aa))
			return
		}

		out = fmt.Sprintf("%s(%s)", fnc, aa[0].S)
		args = aa[0].Args
		return
	}
}

func makeGenericFilterFncHandler(fnc string) HandlerSig {
	return func(aa ...FormattedASTArgs) (out string, args []interface{}, err error) {
		if len(aa) == 0 {
			return fmt.Sprintf("%s()", fnc), nil, nil
		}

		args = make([]interface{}, 0, len(aa))
		auxArgs := make([]string, len(aa))
		for i, a := range aa {
			auxArgs[i] = a.S
			args = append(args, a.Args...)
		}

		out = fmt.Sprintf("%s(%s)", fnc, strings.Join(auxArgs, ", "))
		return
	}
}
