package ql

// Squirrel Sqlizer interface implementators for all ast node types
// This helps us to throw columns into squirrel's select builder

import (
	"fmt"

	"gopkg.in/Masterminds/squirrel.v1"
)

// ToSql concatenates outputs and arguments from all nodes
func (nn ASTNodes) ToSql() (out string, args []interface{}, err error) {
	var _out string
	var _args []interface{}

	for _, s := range nn {
		if _out, _args, err = s.ToSql(); err != nil {
			return
		} else {
			out = out + _out

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
		op = "COLLATE utf8_general_ci " + n.Kind
	}

	return " " + op + " ", nil, nil
}

func (n String) ToSql() (string, []interface{}, error) {
	return "?", []interface{}{n.Value}, nil
}

func (n Number) ToSql() (string, []interface{}, error) {
	return n.Value, nil, nil
}
