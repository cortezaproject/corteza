package rdbms

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

type (
	squirrelConcatExpr struct {
		parts []string
		args  []interface{}
		err   error
	}
)

func SquirrelConcatExpr(args ...interface{}) squirrel.Sqlizer {
	var w = new(squirrelConcatExpr)

	for _, a := range args {
		if w.err != nil {
			break
		}

		switch o := a.(type) {
		case string:
			w.parts = append(w.parts, o)
		case squirrel.Sqlizer:
			p, a, err := o.ToSql()
			w.parts = append(w.parts, p)
			w.args = append(w.args, a...)
			w.err = err
		}
	}

	return w
}

func (w *squirrelConcatExpr) ToSql() (string, []interface{}, error) {
	return strings.Join(w.parts, ""), w.args, w.err
}
