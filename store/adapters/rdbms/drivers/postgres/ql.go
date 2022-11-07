package postgres

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	ref2exp = ql.ExprHandlerMap{
		"concat": {
			Handler: func(args ...exp.Expression) exp.Expression {
				// need to force text type on all arguments
				aa := make([]any, len(args))
				for a := range args {
					aa[a] = exp.NewCastExpression(exp.NewLiteralExpression("?", args[a]), "TEXT")
				}

				return exp.NewSQLFunctionExpression("CONCAT", aa...)
			},
		},

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
		"day": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("EXTRACT",
					exp.NewLiteralExpression("DAY FROM ?", args[0]),
				)
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

		"interval": {
			Handler: func(args ...exp.Expression) exp.Expression {
				// The problem here is that PGSQL, to my findings, doesn't have functions to add/sub dates
				// like MySQL for example.
				//
				// We need to construct an expression in the lines of `INTERVAL 'N UNIT'` which
				// then becomes, for example, d + INTERVAL 'N UNIT'.
				//
				// The problem #2 is that we can't just use value placeholders in string literals
				// nor is there a 2 arg function to make an interval. There is a make_interval function
				// but that one won't do.
				//
				// So...
				// (?  || 'S') makes the interval label a plural because MySQL uses singular and that
				// is what QL supports. PgSQL uses plural, such as years, months, and days.
				//
				// The rest of the expression is just to construct the string which can then be casted to INTERVAL
				// which can then be used in date math, which is done with regular math operators.
				return exp.NewLiteralExpression("(?::INTEGER || ' ' || (?  || 'S'))::INTERVAL", args[1], args[0])
			},
		},

		"date_add": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("(? + ?)", args[0], args[1])
			},
		},

		"date_sub": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("(? - ?)", args[0], args[1])
			},
		},

		// functions currently unsupported in PostgreSQL store backend
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
