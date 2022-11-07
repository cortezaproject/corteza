package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/spf13/cast"
	"strings"
)

type (
	postgresDialect struct{}
)

var (
	_ drivers.Dialect = &postgresDialect{}

	dialect            = &postgresDialect{}
	goquDialectWrapper = goqu.Dialect("postgres")
	quoteIdent         = string(postgres.DialectOptions().QuoteRune)

	nuances = drivers.Nuances{
		HavingClauseMustUseAlias: true,
	}
)

func Dialect() *postgresDialect {
	return dialect
}

func (postgresDialect) Nuances() drivers.Nuances {
	return nuances
}

func (postgresDialect) GOQU() goqu.DialectWrapper  { return goquDialectWrapper }
func (postgresDialect) QuoteIdent(i string) string { return quoteIdent + i + quoteIdent }

func (d postgresDialect) IndexFieldModifiers(attr *dal.Attribute, mm ...dal.IndexFieldModifier) (string, error) {
	return drivers.IndexFieldModifiers(attr, d.QuoteIdent, mm...)
}

func (d postgresDialect) JsonQuote(expr exp.Expression) exp.Expression {
	return exp.NewSQLFunctionExpression("TO_JSON", expr)
}

func (d postgresDialect) JsonExtract(ident exp.Expression, pp ...any) (exp.Expression, error) {
	return DeepIdentJSON(true, ident, pp...), nil
}

func (d postgresDialect) JsonExtractUnquote(ident exp.Expression, pp ...any) (exp.Expression, error) {
	return DeepIdentJSON(false, ident, pp...), nil
}

// JsonArrayContains prepares postgresql compatible comparison of value and JSON array
//
// literal value = multi-value field / plain
// 'value' <@ (v->'f0')::JSONB
//
// single-value field = multi-value field / plain
// v->'f1'->0 <@ (v->'f0')::JSONB
func (d postgresDialect) JsonArrayContains(needle, haystack exp.Expression) (exp.Expression, error) {
	return exp.NewLiteralExpression("(?)::JSONB <@ (?)::JSONB", needle, haystack), nil
}

func (d postgresDialect) TableCodec(m *dal.Model) drivers.TableCodec {
	return drivers.NewTableCodec(m, d)
}

func (d postgresDialect) TypeWrap(dt dal.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	switch c := dt.(type) {
	case *dal.TypeTime:
		return &TypeTime{c}
	}

	return drivers.TypeWrap(dt)
}

func (postgresDialect) AttributeCast(attr *dal.Attribute, val exp.Expression) (exp.Expression, error) {
	var (
		c exp.Expression
	)

	switch attr.Type.(type) {
	case *dal.TypeText:
		c = exp.NewCastExpression(val, "TEXT")

	case *dal.TypeBoolean:
		// convert to text first
		c = exp.NewCastExpression(val, "TEXT")

		// compare to text repr of true
		c = exp.NewBooleanExpression(exp.EqOp, c, exp.NewLiteralExpression(`true::TEXT`))

	default:
		return drivers.AttributeCast(attr, val)

	}

	return exp.NewLiteralExpression("?", c), nil
}

func (postgresDialect) AttributeToColumn(attr *dal.Attribute) (col *ddl.Column, err error) {
	col = &ddl.Column{
		Ident:   attr.StoreIdent(),
		Comment: attr.Label,
		Type: &ddl.ColumnType{
			Null: attr.Type.IsNullable(),
		},
	}

	switch t := attr.Type.(type) {
	case *dal.TypeID:
		col.Type.Name = "BIGINT"
		col.Default = ddl.DefaultID(t.HasDefault, t.DefaultValue)
	case *dal.TypeRef:
		col.Type.Name = "BIGINT"
		col.Default = ddl.DefaultID(t.HasDefault, t.DefaultValue)

	case *dal.TypeTimestamp:
		col.Type.Name = "TIMESTAMP"

		if t.Timezone {
			col.Type.Name += "TZ"
		}

		if t.Precision >= 0 {
			col.Type.Name += fmt.Sprintf("(%d)", t.Precision)
		}

		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *TypeTime:
		col.Type.Name = "TIME"

		if t.Timezone {
			col.Type.Name += "TZ"
		}

		if t.Precision >= 0 {
			col.Type.Name += fmt.Sprintf("(%d)", t.Precision)
		}

		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *dal.TypeDate:
		col.Type.Name = "DATE"
		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *dal.TypeNumber:
		if numType := cast.ToString(t.Meta["rdbms:type"]); numType != "" {
			col.Type.Name = numType
			col.Default = ddl.DefaultNumber(t.HasDefault, 0, t.DefaultValue)
			break
		}

		col.Type.Name = "NUMERIC"

		switch {
		case t.Precision > 0 && t.Scale > 0:
			col.Type.Name += fmt.Sprintf("(%d, %d)", t.Precision, t.Scale)
		case t.Precision > 0:
			col.Type.Name += fmt.Sprintf("(%d)", t.Precision)
		}

		col.Default = ddl.DefaultNumber(t.HasDefault, t.Precision, t.DefaultValue)

	case *dal.TypeText:
		if t.Length > 0 {
			col.Type.Name = fmt.Sprintf("VARCHAR(%d)", t.Length)
		} else {
			col.Type.Name = "TEXT"
		}

		if t.HasDefault {
			// @todo use proper quote type
			col.Default = fmt.Sprintf("%q", t.DefaultValue)
		}

	case *dal.TypeJSON:
		col.Type.Name = "JSONB"
		if col.Default, err = ddl.DefaultJSON(t.HasDefault, t.DefaultValue); err != nil {
			return nil, err
		}

	case *dal.TypeGeometry:
		// @todo geometry type
		col.Type.Name = "JSONB"

	case *dal.TypeBlob:
		col.Type.Name = "BYTEA"

	case *dal.TypeBoolean:
		col.Type.Name = "BOOLEAN"
		col.Default = ddl.DefaultBoolean(t.HasDefault, t.DefaultValue)

	case *dal.TypeUUID:
		col.Type.Name = "UUID"

	default:
		return nil, fmt.Errorf("unsupported column type: %s ", t.Type())
	}

	return
}

func (d postgresDialect) ExprHandler(n *ql.ASTNode, args ...exp.Expression) (expr exp.Expression, err error) {
	switch ref := strings.ToLower(n.Ref); ref {
	case "concat":
		// need to force text type on all arguments
		aa := make([]any, len(args))
		for a := range args {
			aa[a] = exp.NewCastExpression(exp.NewLiteralExpression("?", args[a]), "TEXT")
		}

		return exp.NewSQLFunctionExpression("CONCAT", aa...), nil

	case "in":
		return drivers.OpHandlerIn(d, n, args...)

	case "nin":
		return drivers.OpHandlerNotIn(d, n, args...)

	}

	return ref2exp.RefHandler(n, args...)
}

func (d postgresDialect) handleInOperator(n *ql.ASTNode, negate bool, args ...exp.Expression) (expr exp.Expression, err error) {
	if n.Args[1] != nil && n.Args[1].Meta["dal.Attribute"] != nil && n.Args[1].Meta["dal.Attribute"].(*dal.Attribute).MultiValue {
		// if right-side argument is multi-value attribute,
		// then we need to adjust the arguments a bit:
		//  left side, if it is a value, is encoded as JSON
		//             if ref we access JSON encoded value
		//
		//  right side, access JSON encoded array of values.
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
							When(exp.NewBooleanExpression(exp.EqOp, args[a], drivers.LiteralTRUE), exp.NewLiteralExpression(`'true'`)).
							When(exp.NewBooleanExpression(exp.EqOp, args[a], drivers.LiteralFALSE), exp.NewLiteralExpression(`'false'`)).
							Else(drivers.LiteralNULL)
					} else {
						args[a] = exp.NewSQLFunctionExpression("TO_JSON", args[a])
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
	}

	return
}
