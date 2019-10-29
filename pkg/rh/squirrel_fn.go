package rh

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

type (
	// Waiting for PR to be merged:
	// https://github.com/Masterminds/squirrel/pull/206
	//
	// then we can move to squirrel.Fn(...)
	squirrelFunction struct {
		name  string
		fargs []squirrel.Sqlizer
	}
)

func SquirrelFunction(name string, args ...squirrel.Sqlizer) *squirrelFunction {
	return &squirrelFunction{name: name, fargs: args}
}

func (f squirrelFunction) ToSql() (sql string, args []interface{}, err error) {
	var (
		aSql  string
		aArgs []interface{}
	)

	sql = f.name + "("
	args = make([]interface{}, 0)
	for a := 0; a < len(f.fargs); a++ {
		if a > 0 {
			sql += ", "
		}

		aSql, aArgs, err = f.fargs[a].ToSql()
		if err != nil {
			return
		}

		sql += aSql
		args = append(args, aArgs...)
	}
	sql += ")"

	return
}

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
