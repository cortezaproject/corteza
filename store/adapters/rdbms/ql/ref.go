package ql

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	ExprHandlerMap map[string]*ExprHandler

	ExprHandler struct {
		Handler  func(...exp.Expression) exp.Expression
		HandlerE func(...exp.Expression) (exp.Expression, error)
	}
)

var (
	ref2exp = ExprHandlerMap{
		// keywords
		"null": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("NULL")
			},
		},
		"nnull": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("NOT NULL")
			},
		},

		// operators
		"not": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("(NOT ?)", args[0])
			},
		},

		// - bool
		"and": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewExpressionList(exp.AndType, args...)
			},
		},
		"or": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewExpressionList(exp.OrType, args...)
			},
		},

		// - comp.
		"eq": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.EqOp, args[0], args[1])
			},
		},
		"ne": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.NeqOp, args[0], args[1])
			},
		},
		"lt": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.LtOp, args[0], args[1])
			},
		},
		"le": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.LteOp, args[0], args[1])
			},
		},
		"gt": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.GtOp, args[0], args[1])
			},
		},
		"ge": {
			//Handler: makeGenericCompHandler(">="),
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.GteOp, args[0], args[1])
			},
		},

		// - math
		"add": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("? + ?", args[0], args[1])
			},
		},
		"sub": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("? - ?", args[0], args[1])
			},
		},
		"mult": {
			Handler: func(args ...exp.Expression) exp.Expression {
				// Handling COUNT(*) scenario:
				// ql parser interprets * as a multiplier
				//
				// in situations where there are 0 arguments we'll just
				// return star
				if len(args) == 0 {
					return exp.Star()
				}

				return exp.NewLiteralExpression("? * ?", args[0], args[1])
			},
		},
		"div": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("? / ?", args[0], args[1])
			},
		},

		// - strings
		"concat": {
			Handler: func(args ...exp.Expression) exp.Expression {
				aa := make([]any, len(args))
				for a := range args {
					aa[a] = args[a]
				}

				return exp.NewSQLFunctionExpression("CONCAT", aa...)
			},
		},

		"interval": {
			Handler: func(args ...exp.Expression) exp.Expression {

				// The problem here is similar to the one with PgSQL... but more complicated...
				// The interval duration identifiers are keywords and can't be bound via value placeholders.
				// This forces us to interpolate the thing when building the literal expression below.
				//
				// Since goqu doesn't let me get the raw value of the expression... I was forced
				// to do this calamity
				_, aa, err := goqu.Select(args[0]).ToSQL()
				if err != nil {
					// This error should never occur
					panic(err)
				}
				intv := aa[0].(string)

				return exp.NewLiteralExpression(fmt.Sprintf("INTERVAL ? %s", intv), args[1])
			},
		},

		"date_add": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DATE_ADD", args[0], args[1])
			},
		},

		"date_sub": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DATE_SUB", args[0], args[1])
			},
		},

		// @todo better negation?
		"like": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.ILikeOp, args[0], args[1])
			},
		},
		"nlike": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.NotILikeOp, args[0], args[1])
			},
		},

		"is": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.IsOp, args[0], args[1])
			},
		},

		"nis": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewBooleanExpression(exp.IsNotOp, args[0], args[1])
			},
		},

		"group": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewLiteralExpression("(?)", args[0])
			},
		},

		// - filtering
		"date_format": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DATE_FORMAT", args[0], args[1])
			},
		},
		"now": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("NOW")
			},
		},
		"quarter": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("QUARTER", args[0])
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
		"date": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DATE", args[0])
			},
		},
		"count": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("COUNT", args[0])
			},
		},
		"sum": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("SUM", args[0])
			},
		},
		"avg": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("AVG", args[0])
			},
		},
		"min": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("MIN", args[0])
			},
		},
		"max": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("MAX", args[0])
			},
		},
	}
)

func (ee ExprHandlerMap) ExprHandlers() (out ExprHandlerMap) {
	// Set the default values; make sure to copy the global map
	out = make(map[string]*ExprHandler, len(ref2exp))
	for name, expr := range ref2exp {
		out[name] = expr
	}

	// Overwrite with customs
	for name, expr := range ee {
		out[name] = expr
	}
	return
}

func (ee ExprHandlerMap) RefHandler(n *ql.ASTNode, args ...exp.Expression) (exp.Expression, error) {
	r := strings.ToLower(n.Ref)

	if ee[r] == nil {
		return nil, fmt.Errorf("unknown ref %q", r)
	}

	if ee[r].Handler != nil {
		return ee[r].Handler(args...), nil
	}

	return ee[r].HandlerE(args...)
}
