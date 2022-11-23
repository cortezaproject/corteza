package drivers

import (
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	Nuances struct {
		// HavingClauseMustUseAlias
		// For example, Postgres and SQLite require
		// use of aliases inside HAVING clause and
		// MySQL does not.
		HavingClauseMustUseAlias bool

		// TwoStepUpsert allows support for databases which don't have an upsert
		// or we simply couldn't figure out how to make work.
		//
		// TwoStepUpsert uses the context from the update statement to figure out
		// if it needs to do an insert.
		TwoStepUpsert bool
	}

	Dialect interface {
		// Nuances returns dialect nuances
		// subtle differences between RDBMS implementations that
		// should be handled on common code
		Nuances() Nuances

		// GOQU returns goqu's dialect wrapper struct
		GOQU() goqu.DialectWrapper

		JsonQuote(exp.Expression) exp.Expression

		// JsonExtract returns expression that returns a value from  inside JSON document
		//
		// Use this when you want use JSON encoded value
		JsonExtract(exp.Expression, ...any) (exp.Expression, error)

		// JsonExtractUnquote returns expression that returns a value from  inside JSON document:
		//
		// Use this when you want to use unencoded value!
		JsonExtractUnquote(exp.Expression, ...any) (exp.Expression, error)

		// JsonArrayContains generates expression JSON array containment check expression
		//
		// Literal values need to be JSON docs!
		//
		// @todo recheck if we really need JsonArrayContains on Dialect interface
		JsonArrayContains(needle, haystack exp.Expression) (exp.Expression, error)

		// AttributeCast prepares complex SQL expression that verifies
		// arbitrary string value in the db and casts it to b used in
		// comparison or soring expression
		AttributeCast(*dal.Attribute, exp.Expression) (exp.Expression, error)

		// TableCodec returns table codec (encodes & decodes data to/from db table)
		TableCodec(*dal.Model) TableCodec

		// TypeWrap returns driver's type implementation for a particular attribute type
		TypeWrap(dal.Type) Type

		QuoteIdent(string) string

		// AttributeToColumn converts attribute to column defunition
		AttributeToColumn(*dal.Attribute) (*ddl.Column, error)

		// ExprHandler returns driver specific expression handling
		ExprHandler(*ql.ASTNode, ...exp.Expression) (exp.Expression, error)

		// OrderedExpression returns compatible expression for ordering
		//
		// Database understand order modifiers differently. For example, MySQL does not know
		// about NULLS FIRST/LAST. Drivers should gracefully handle this.
		OrderedExpression(exp.Expression, exp.SortDirection, exp.NullSortType) exp.OrderedExpression
	}
)

func init() {
	goqu.SetDefaultPrepared(true)
}

func IndexFieldModifiers(attr *dal.Attribute, quoteIdent func(i string) string, mm ...dal.IndexFieldModifier) (string, error) {
	var (
		modifier string
		out      = quoteIdent(attr.StoreIdent())
	)

	for _, m := range mm {
		switch m {
		case dal.IndexFieldModifierLower:
			modifier = "LOWER"
		default:
			return "", fmt.Errorf("unknown index field modifier: %s", m)
		}

		out = fmt.Sprintf("%s(%s)", modifier, out)
	}

	return out, nil
}

func OpHandlerIn(d Dialect, n *ql.ASTNode, args ...exp.Expression) (expr exp.Expression, err error) {
	return opHandlerIn(d, n, false, args...)
}

func OpHandlerNotIn(d Dialect, n *ql.ASTNode, args ...exp.Expression) (expr exp.Expression, err error) {
	return opHandlerIn(d, n, true, args...)
}

func opHandlerIn(d Dialect, n *ql.ASTNode, negate bool, args ...exp.Expression) (expr exp.Expression, err error) {
	if len(n.Args) == 2 && n.Args[1] != nil && n.Args[1].Meta["dal.Attribute"] != nil && n.Args[1].Meta["dal.Attribute"].(*dal.Attribute).MultiValue {
		// if right-side argument is multi-value attribute,
		// then we need to adjust the arguments a bit:
		// 1) left side, if it is a value, is encoded as JSON
		// 2)            if ref we access JSON encoded value
		//
		// right side, access JSON encoded array of values.
		//
		//
		//
		//
		for a := range n.Args {
			left := a == 0

			switch {
			case n.Args[a].Meta != nil && n.Args[a].Meta["dal.Attribute"] != nil:
				// symbol, ident probably...
				var (
					attr       = n.Args[a].Meta["dal.Attribute"].(*dal.Attribute)
					model      = n.Args[a].Meta["dal.Model"].(*dal.Model)
					storeIdent = exp.NewIdentifierExpression(
						"",
						model.Ident,
						attr.StoreIdent(),
					)

					_, isJSON = attr.Store.(*dal.CodecRecordValueSetJSON)
				)

				if attr.MultiValue {
					if left {
						return nil, fmt.Errorf("multi-value attribute %s cannot be used as left-side argument of IN operator", attr.Ident)
					}

					args[a], err = d.JsonExtract(storeIdent, attr.Ident)
				} else {
					if !left {
						return nil, fmt.Errorf("single-value attribute %s cannot be used as right-side argument of IN operator", attr.Ident)
					}

					if isJSON {
						args[a], err = d.JsonExtract(storeIdent, attr.Ident, 0)
					} else if attr.Type.Type() == dal.AttributeTypeBoolean {
						// SQLite converts boolean to integer but JSON stores boolean as boolean
						args[a] = exp.NewCaseExpression().
							When(exp.NewBooleanExpression(exp.EqOp, args[a], LiteralTRUE), exp.NewLiteralExpression(`'true'`)).
							When(exp.NewBooleanExpression(exp.EqOp, args[a], LiteralFALSE), exp.NewLiteralExpression(`'false'`)).
							Else(LiteralNULL)
					} else {
						args[a] = d.JsonQuote(args[a])
					}
				}

				if err != nil {
					return nil, err
				}

			case a == 0 && n.Args[a].Value != nil:
				// for 1st arg only, when value
				var jsonDoc []byte
				jsonDoc, err = json.Marshal(n.Args[a].Value.V.Get())
				if err != nil {
					return nil, err
				}

				// encode it as json
				args[a] = exp.NewLiteralExpression("?", string(jsonDoc))
			}
		}

		expr, err = d.JsonArrayContains(args[0], args[1])
		if err != nil {
			return
		}

		if negate {
			expr = exp.NewLiteralExpression("NOT ?", expr)
		}

		return
	}

	return nil, fmt.Errorf("unsupported IN operator")
}
