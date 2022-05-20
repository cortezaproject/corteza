package ql

import (
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	exprHandler struct {
		Handler func(...exp.Expression) exp.Expression
	}
)

var (
	ref2exp = map[string]*exprHandler{
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
		//
		//// - filtering
		//"now": {
		//	Result:  wrapRes("DateTime"),
		//	Handler: makeGenericFilterFncHandler("NOW"),
		//},
		//"quarter": {
		//	Args:    collectParams(true, "DateTime"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericFilterFncHandler("QUARTER"),
		//},
		//"year": {
		//	Args:    collectParams(true, "DateTime"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericFilterFncHandler("YEAR"),
		//},
		//"month": {
		//	Args:    collectParams(true, "DateTime"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericFilterFncHandler("MONTH"),
		//},
		//"date": {
		//	Args:    collectParams(true, "DateTime"),
		//	Result:  wrapRes("Number"),
		//	Handler: makeGenericFilterFncHandler("DAY"),
		//},
	}
)
