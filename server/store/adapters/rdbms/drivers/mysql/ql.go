package mysql

import (
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	ref2exp = ql.ExprHandlerMap{
		"std": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("STD", args[0])
			},
		},
	}.ExprHandlers()
)
