package ql

import (
	"fmt"
)

// ToSQL concatenates outputs and arguments from all nodes
func (nn ASTNodes) ToSQL() (out string, args []interface{}, err error) {
	var _out string
	var _args []interface{}

	for _, s := range nn {
		if _out, _args, err = s.ToSQL(); err != nil {
			return
		} else {
			if _, ok := s.(ASTNodes); ok {
				// Nested nodes should be wrapped
				// Example: ((A) AND (B))
				out = out + "(" + _out + ")"
			} else {
				out = out + _out
			}

			args = append(args, _args...)
		}
	}

	return out, args, err
}

// ToSQL concatenates outputs and arguments from all nodes, comma delimited
func (nn ASTSet) ToSQL() (out string, args []interface{}, err error) {
	var _out string
	var _args []interface{}

	for i, s := range nn {
		if _out, _args, err = s.ToSQL(); err != nil {
			return
		} else {
			if i > 0 {
				out = out + ", "
			}

			out = out + _out

			args = append(args, _args...)
		}
	}

	return out, args, err
}

// ToSQL returns column alias expression or output of underlying expression's ToSQL()
func (n Column) ToSQL() (string, []interface{}, error) {
	if n.Alias != "" {
		panic("@todo reimplement this")
		//return squirrel.Alias(n.Expr, n.Alias).ToSQL()
	} else {
		return n.Expr.ToSQL()
	}
}

func (n Ident) ToSQL() (string, []interface{}, error) {
	return n.Value, n.Args, nil
}

func (n LNull) ToSQL() (string, []interface{}, error) {
	return "NULL", nil, nil
}

func (n LBoolean) ToSQL() (string, []interface{}, error) {
	if n.Value {
		return "TRUE", nil, nil
	} else {
		return "FALSE", nil, nil
	}
}

func (n Function) ToSQL() (string, []interface{}, error) {
	if paramsSql, args, err := n.Arguments.ToSQL(); err != nil {
		return "", nil, err
	} else {
		return fmt.Sprintf("%s(%s)", n.Name, paramsSql), args, nil
	}

}

func (n Keyword) ToSQL() (string, []interface{}, error) {
	return n.Keyword, nil, nil
}

func (n Interval) ToSQL() (string, []interface{}, error) {
	return fmt.Sprintf("INTERVAL ? %s", n.Unit), []interface{}{n.Value}, nil
}

func (n Operator) ToSQL() (string, []interface{}, error) {
	var op = n.Kind

	switch n.Kind {
	case "LIKE", "NOT LIKE":
		// Make sure we are doing case insensitive search
		op = QueryEncoder.CaseInsensitiveLike(n.Kind == "NOT LIKE")
	}

	return " " + op + " ", nil, nil
}

func (n LString) ToSQL() (string, []interface{}, error) {
	return "?", []interface{}{n.Value}, nil
}

func (n LNumber) ToSQL() (string, []interface{}, error) {
	return n.Value, nil, nil
}

func (n NodeF) ToSQL() (string, []interface{}, error) {
	var (
		// used for sprintf to complete the base expression
		fArgs []interface{}

		// collection of al args from ToSQL() that
		// are passed on to the caller
		adtArgs []interface{}
	)

	for i, s := range n.Arguments {
		// When provided, apply the replacer over the arguments of the node.
		// We can skip any node other then LString as thats the only one we can apply it to (currently)
		if n.replacer != nil {
			if c, ok := s.(LString); ok {
				c.Value = n.replacer(c.Value)
				// Updating the originals as we're dealing with values
				n.Arguments[i] = c
				s = c
			}
		}

		if fa, aa, err := s.ToSQL(); err != nil {
			return "", nil, err
		} else {
			fArgs = append(fArgs, fa)
			adtArgs = append(adtArgs, aa...)
		}
	}

	return fmt.Sprintf(n.Expr, fArgs...), adtArgs, nil
}
