package postgres

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	ref2exp = ql.ExprHandlerMap{
		// filtering
		"now": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("NOW")
			},
		},
		"quarter": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("EXTRACT",
					exp.NewLiteralExpression("QUARTER FROM ?", args[0]),
				)
			},
		},
		"year": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("EXTRACT",
					exp.NewLiteralExpression("YEAR FROM ?", args[0]),
				)
			},
		},
		"month": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("EXTRACT",
					exp.NewLiteralExpression("MONTH FROM ?", args[0]),
				)
			},
		},
		"timestamp": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("?::TIMESTAMPTZ", args[0])
			},
		},
		"date": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("?::DATE", args[0])
			},
		},
		"time": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("DATE_TRUNC('second', ?::TIME)::TIME", args[0])
			},
		},

		// @todo replace given argument before constructing sql
		"date_format": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("TO_CHAR",
					exp.NewLiteralExpression("?::TIMESTAMPTZ", args[0]),
					exp.NewLiteralExpression("?::TEXT", translateDateFormatParam(args[1])),
				)
			},
		},

		// functions currently unsupported in PostgreSQL store backend
		//"DATE_ADD": {
		//	Handler: func(args ...exp.Expression) exp.Expression {
		//		return exp.NewLiteralExpression("")
		//	},
		//},
		//"DATE_SUB": {
		//	Handler: func(args ...exp.Expression) exp.Expression {
		//		return exp.NewLiteralExpression("")
		//	},
		//},
		//"STD": {
		//	Handler: func(args ...exp.Expression) exp.Expression {
		//		return exp.NewLiteralExpression("")
		//	},
		//},
	}.ExprHandlers()
)

func translateDateFormatParam(e interface{}) interface{} {
	le, ok := e.(exp.LiteralExpression)
	if !ok {
		return e
	}

	args := le.Args()
	if len(args) > 0 {
		return dateFormatReplacer(fmt.Sprintf("%s", args[0]))
	}

	return e
}

func dateFormatReplacer(s string) string {
	return strings.NewReplacer(
		// @todo Doing ...%dT%H... (for iso timestamp) pgsql doesn't format it correctly
		// so I'm covering this edge case.
		// We should fix this properly when we redo record storage.
		`%dT%H`, `DD"T"HH24`,

		`%a`, `Dy`,
		`%b`, `Mon`,
		`%c`, `FMMM`,
		`%d`, `DD`,
		`%e`, `FMDD`,
		`%f`, `US`,
		`%H`, `HH24`,
		`%h`, `HH12`,
		`%I`, `HH12`,
		`%i`, `MI`,
		`%j`, `DDD`,
		`%k`, `FMHH24`,
		`%l`, `FMHH12`,
		`%M`, `FMMonth`,
		`%m`, `MM`,
		`%p`, `AM`,
		`%r`, `HH12:MI:SS AM`,
		`%S`, `SS`,
		`%s`, `SS`,
		`%T`, `HH24:MI:SS`,
		`%W`, `FMDay`,
		`%Y`, `YYYY`,
		`%y`, `YY`,
		`%%`, `%`,
	).Replace(s)
}
