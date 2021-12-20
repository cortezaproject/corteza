package ql

// Squirrel Sqlizer interface implementators for all ast node types
// This helps us to throw columns into squirrel's select builder

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

// ToSql concatenates outputs and arguments from all nodes
func (nn ASTNodes) ToSql() (out string, args []interface{}, err error) {
	var _out string
	var _args []interface{}

	for _, s := range nn {
		if _out, _args, err = s.ToSql(); err != nil {
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

// ToSql concatenates outputs and arguments from all nodes, comma delimited
func (nn ASTSet) ToSql() (out string, args []interface{}, err error) {
	var _out string
	var _args []interface{}

	for i, s := range nn {
		if _out, _args, err = s.ToSql(); err != nil {
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

// ToSql returns column alias expression or output of underlying expression's ToSql()
func (n Column) ToSql() (string, []interface{}, error) {
	if n.Alias != "" {
		return squirrel.Alias(n.Expr, n.Alias).ToSql()
	} else {
		return n.Expr.ToSql()
	}
}

func (n Ident) ToSql() (string, []interface{}, error) {
	return n.Value, n.Args, nil
}

func (n LNull) ToSql() (string, []interface{}, error) {
	return "NULL", nil, nil
}

func (n LBoolean) ToSql() (string, []interface{}, error) {
	if n.Value {
		return "TRUE", nil, nil
	} else {
		return "FALSE", nil, nil
	}
}

func (n Function) ToSql() (string, []interface{}, error) {
	if paramsSql, args, err := n.Arguments.ToSql(); err != nil {
		return "", nil, err
	} else {
		return fmt.Sprintf("%s(%s)", n.Name, paramsSql), args, nil
	}

}

func (n Keyword) ToSql() (string, []interface{}, error) {
	return n.Keyword, nil, nil
}

func (n Interval) ToSql() (string, []interface{}, error) {
	return fmt.Sprintf("INTERVAL ? %s", n.Unit), []interface{}{n.Value}, nil
}

func (n Operator) ToSql() (string, []interface{}, error) {
	var op = n.Kind

	switch n.Kind {
	case "LIKE", "NOT LIKE":
		// Make sure we are doing case insensitive search
		op = QueryEncoder.CaseInsensitiveLike(n.Kind == "NOT LIKE")
	}

	return " " + op + " ", nil, nil
}

func (n LString) ToSql() (string, []interface{}, error) {
	return "?", []interface{}{n.Value}, nil
}

func (n LNumber) ToSql() (string, []interface{}, error) {
	return n.Value, nil, nil
}

func (n NodeF) ToSql() (string, []interface{}, error) {
	var (
		// used for sprintf to complete the base expression
		fArgs []interface{}

		// collection of al args from ToSql() that
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

		if fa, aa, err := s.ToSql(); err != nil {
			return "", nil, err
		} else {
			fArgs = append(fArgs, fa)
			adtArgs = append(adtArgs, aa...)
		}
	}

	return fmt.Sprintf(n.Expr, fArgs...), adtArgs, nil
}
