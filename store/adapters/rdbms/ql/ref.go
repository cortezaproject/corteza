package ql

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/ql"
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

		// functions
		// - aggregation
		//"count": {
		//	Args:    collectParams(false, "Any"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericAggFncHandler("COUNT"),
		//},
		//"sum": {
		//	Args:    collectParams(true, "Any"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericAggFncHandler("SUM"),
		//},
		//"max": {
		//	Args:    collectParams(true, "Any"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericAggFncHandler("MAX"),
		//},
		//"min": {
		//	Args:    collectParams(true, "Any"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericAggFncHandler("MIN"),
		//},
		//"avg": {
		//	Args:    collectParams(true, "Any"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericAggFncHandler("AVG"),
		//},

		// - filtering
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
		"date": {
			Handler: func(args ...exp.Expression) exp.Expression {
				return exp.NewSQLFunctionExpression("DAY", args[0])
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
