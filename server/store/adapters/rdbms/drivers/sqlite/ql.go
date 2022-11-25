package sqlite

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	ref2exp = ql.ExprHandlerMap{
		// filtering
		"now": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DATE",
					exp.NewLiteralExpression("'NOW'"),
				)
			},
		},
		"quarter": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("(CAST(STRFTIME('%m', ?) AS INTEGER) + 2) / 3", args[0])
			},
		},
		"year": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("STRFTIME",
					exp.NewLiteralExpression("'%Y'"),
					args[0],
				)
			},
		},
		"month": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("STRFTIME",
					exp.NewLiteralExpression("'%m'"),
					args[0],
				)
			},
		},
		"date": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("STRFTIME",
					exp.NewLiteralExpression("'%Y-%m-%dT00:00:00Z'"),
					args[0],
				)
			},
		},
		"day": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("STRFTIME",
					exp.NewLiteralExpression("'%d'"),
					args[0],
				)
			},
		},
		"week": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("STRFTIME",
					exp.NewLiteralExpression("'%W'"),
					args[0],
				)
			},
		},
		"datetime": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DATETIME", args[0])
			},
		},
		"timestamp": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DATETIME", args[0])
			},
		},
		"date_format": {
			HandlerE: func(args ...exp.Expression) (exp.Expression, error) {
				format, err := supportedDateFormatParams(args[1])
				if err != nil {
					return nil, err
				}

				return exp.NewSQLFunctionExpression("STRFTIME",
					format,
					args[0],
				), nil
			},
		},

		// functions currently unsupported in SQLite store backend
		// "DATE_ADD": {
		//	Handler: func(args ...exp.Expression) exp.Expression {
		//		return exp.NewLiteralExpression("")
		//	},
		// },
		// "DATE_SUB": {
		//	Handler: func(args ...exp.Expression) exp.Expression {
		//		return exp.NewLiteralExpression("")
		//	},
		// },
		// "STD": {
		//	Handler: func(args ...exp.Expression) exp.Expression {
		//		return exp.NewLiteralExpression("")
		//	},
		// },
	}.ExprHandlers()

	supportedSubstitutions = map[string]bool{
		"d": true,
		"H": true,
		"j": true,
		"m": true,
		"M": true,
		"S": true,
		"w": true,
		"W": true,
		"Y": true,
		"%": true,
	}
)

func supportedDateFormatParams(e interface{}) (interface{}, error) {
	le, ok := e.(exp.LiteralExpression)
	if !ok {
		return e, fmt.Errorf("unknown date format")
	}

	var format string
	args := le.Args()
	if len(args) > 0 {
		format = dateFormatReplacer(fmt.Sprintf("%s", args[0]))
	} else {
		return e, fmt.Errorf("date format not found")
	}

	r := regexp.MustCompile(`%(?P<sub>.)`)

	for _, m := range r.FindAllStringSubmatch(format, -1) {
		if len(m) == 0 {
			continue
		}

		if _, ok := supportedSubstitutions[m[1]]; !ok {
			return e, fmt.Errorf("format substitution not supported: %%%s", m[1])
		}
	}

	return format, nil
}

func dateFormatReplacer(format string) string {
	return strings.NewReplacer(
		`%i`, `%M`,
		`%U`, `%W`,
	).Replace(format)
}
