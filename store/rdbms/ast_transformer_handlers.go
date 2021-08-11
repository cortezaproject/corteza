package rdbms

import (
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/qlng"
)

func (t *astTransformer) handleSymbol(n *qlng.ASTNode) (string, []interface{}, error) {
	return n.Symbol, nil, nil
}

func (t *astTransformer) handleValue(n *qlng.ASTNode) (string, []interface{}, error) {
	if t.placeholders {
		return "?", wrapIntfs(n.Value.V.Get()), nil
	}

	vl := n.Value.V.Get()

	// such type casting is ok here, as the typedvalue should already be this
	switch n.Value.V.Type() {
	case "Boolean":
		if vl.(bool) {
			return "TRUE", nil, nil
		}
		return "FALSE", nil, nil

	case "Integer":
		return fmt.Sprintf("%d", vl.(int64)), nil, nil

	case "UnsignedInteger":
		return fmt.Sprintf("%d", vl.(uint64)), nil, nil

	case "Float":
		return fmt.Sprintf("%f", vl.(float64)), nil, nil

	case "String":
		return fmt.Sprintf("'%s'", vl.(string)), nil, nil

	case "DateTime":
		dt := vl.(*time.Time)
		return fmt.Sprintf("'%s'", dt.Format(time.RFC3339)), nil, nil
	}

	return "", nil, fmt.Errorf("unsupported value type: %s", n.Value.V.Type())
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
		if fnc == "COUNT" && len(aa) == 0 {
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

func makeGenericTypecastHandler(t string) HandlerSig {
	return func(aa ...FormattedASTArgs) (out string, args []interface{}, err error) {
		if len(aa) != 1 {
			err = fmt.Errorf("expecting 1 argument, got %d", len(aa))
			return
		}

		out = fmt.Sprintf("CAST(%s AS %s)", aa[0].S, t)
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
