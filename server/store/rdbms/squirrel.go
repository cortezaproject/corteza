package rdbms

import (
	"github.com/Masterminds/squirrel"
)

type (
	squirrelConcatExpr struct {
		f    squirrel.PlaceholderFormat
		args []interface{}
	}
)

func SquirrelConcatExpr(args ...interface{}) *squirrelConcatExpr {
	return &squirrelConcatExpr{args: args, f: squirrel.Question}
}

func (w *squirrelConcatExpr) PlaceholderFormat(f squirrel.PlaceholderFormat) *squirrelConcatExpr {
	w.f = f
	return w
}

func (w *squirrelConcatExpr) ToSql() (sql string, args []interface{}, err error) {
	var (
		partSql  string
		partArgs []interface{}
	)

	for _, a := range w.args {
		switch o := a.(type) {
		case string:
			sql += o
		case squirrel.Sqlizer:
			if partSql, partArgs, err = o.ToSql(); err != nil {
				return
			}

			sql += partSql
			args = append(args, partArgs...)
		}
	}

	if sql, err = w.f.ReplacePlaceholders(sql); err != nil {
		return
	}

	return
}
