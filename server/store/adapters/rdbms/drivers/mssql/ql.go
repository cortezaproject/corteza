package mssql

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	ref2exp = ql.ExprHandlerMap{
		"interval": {
			HandlerE: func(e ...exp.Expression) (exp.Expression, error) {
				return nil, fmt.Errorf("@todo not implemented")
			},
		},

		"date_add": {
			HandlerE: func(e ...exp.Expression) (exp.Expression, error) {
				return nil, fmt.Errorf("@todo not implemented")
			},
		},

		"date_sub": {
			HandlerE: func(e ...exp.Expression) (exp.Expression, error) {
				return nil, fmt.Errorf("@todo not implemented")
			},
		},

		"date_format": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression(
					"FORMAT",
					args[0],
					translateDateFormatParam(args[1]),
				)
			},
		},

		"now": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("GETDATE")
			},
		},
		"date": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("CONVERT(date,?)", args[0])
			},
		},

		"quarter": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("DATEPART(QUARTER, ?)", args[0])
			},
		},
		"year": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("YEAR", args[0])
			},
		},
		"month": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("MONTH", args[0])
			},
		},
		"day": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DAY", args[0])
			},
		},

		"timestamp": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("CONVERT(DATETIME, ?)", args[0])
			},
		},
		"time": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("CONVERT(TIME, ?)", args[0])
			},
		},
	}.ExprHandlers()

	// @todo consider
	// supportedSubstitutions = map[string]bool{
	// 	"d": true,
	// 	"H": true,
	// 	"j": true,
	// 	"m": true,
	// 	"M": true,
	// 	"S": true,
	// 	"w": true,
	// 	"W": true,
	// 	"Y": true,
	// 	"%": true,
	// }
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

		//  - this is day of month from 01-31
		`%d`, `dd`,
		//  - this is the month number from 01-12 !! %m is from 0
		`%c`, `MM`,
		//  - this is the month number from 01-12 !! %m is from 0
		`%m`, `MM`,

		//  - month name abbreviated
		`%b`, `MMM`,
		//  - this is the month spelled out
		`%M`, `MMMM`,

		//  - this is the year with two digits
		`%y`, `yy`,
		//  - this is the year with four digits
		`%Y`, `yyyy`,

		//  - this is the hour from 01-12
		`%l`, `hh`,
		//  - this is the hour from 00-23
		`%k`, `HH`,

		//  - this is the minute from 00-59
		`%i`, `mm`,
		//  - this is the second from 00-59
		`%s`, `ss`,
		//  - this is the second from 00-59
		`%S`, `ss`,

		//  - this shows either AM or PM
		`%r`, `tt`,
	).Replace(s)
}
