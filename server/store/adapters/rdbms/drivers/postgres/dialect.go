package postgres

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/doug-martin/goqu/v9/sqlgen"
	"github.com/spf13/cast"
)

type (
	postgresDialect struct{}
)

var (
	_ drivers.Dialect = &postgresDialect{}

	dialect            = &postgresDialect{}
	goquDialectWrapper = goqu.Dialect("postgres")
	goquDialectOptions = postgres.DialectOptions()
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

func (postgresDialect) GOQU() goqu.DialectWrapper                 { return goquDialectWrapper }
func (postgresDialect) DialectOptions() *sqlgen.SQLDialectOptions { return goquDialectOptions }
func (postgresDialect) QuoteIdent(i string) string                { return quoteIdent + i + quoteIdent }

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

func (postgresDialect) AttributeCast(attr *dal.Attribute, val exp.Expression) (expr exp.Expression, err error) {
	switch attr.Type.(type) {
	case *dal.TypeText:
		expr = exp.NewCastExpression(val, "TEXT")

	case *dal.TypeBoolean:
		// convert to text first
		expr = exp.NewCastExpression(val, "TEXT")

		// compare to text representation of true
		expr = exp.NewBooleanExpression(exp.EqOp, expr, exp.NewLiteralExpression(`true::TEXT`))

	default:
		return drivers.AttributeCast(attr, val)

	}

	return
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
		col.Type.Name = "NUMERIC"
		col.Default = ddl.DefaultID(t.HasDefault, t.DefaultValue)
	case *dal.TypeRef:
		col.Type.Name = "NUMERIC"
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

	case *dal.TypeEnum:
		col.Type.Name = "TEXT"

	default:
		return nil, fmt.Errorf("unsupported column type: %s ", t.Type())
	}

	return
}

func (postgresDialect) ColumnFits(target, assert *ddl.Column) bool {
	targetType, targetName, targetMeta := ddl.ParseColumnTypes(target)
	assertType, assertName, assertMeta := ddl.ParseColumnTypes(assert)

	// If everything matches up perfectly use that
	if assertType == targetType {
		return true
	}

	// See if we can guess it
	// [the type of the target column][what types fit the target col. type]
	matches := map[string]map[string]bool{
		"numeric": {
			"text":    true,
			"varchar": true,
			"bigint":  true,
			// @note this isn't entirely correct, but record revision lapsus made us
			"integer": true,
		},
		"integer": {
			"text":    true,
			"varchar": true,
			"bigint":  true,
			"numeric": true,
		},
		"timestamp": {
			"text":    true,
			"varchar": true,

			"timestamptz": true,
		},
		"timestamptz": {
			"text":    true,
			"varchar": true,
		},
		"time": {
			"text":    true,
			"varchar": true,

			"timetz": true,
		},
		"timetz": {
			"text":    true,
			"varchar": true,
		},
		"date": {
			"text":    true,
			"varchar": true,
		},
		"text": {},
		"varchar": {
			"text": true,
		},
		"jsonb": {},
		"bytea": {},
		"boolean": {
			"numeric": true,
		},
		"uuid": {
			"text":    true,
			"varchar": true,
		},
	}

	baseMatch := assertName == targetName || matches[assertName][targetName]

	// Special cases
	switch {
	case assertName == "varchar" && targetName == "varchar":
		// Check varchar size
		for i := len(assertMeta); i < 1; i++ {
			assertMeta = append(assertMeta, "0")
		}
		for i := len(targetMeta); i < 1; i++ {
			targetMeta = append(targetMeta, "0")
		}

		return baseMatch && cast.ToInt(assertMeta[0]) <= cast.ToInt(targetMeta[0])

	case assertName == "numeric" && targetName == "numeric":
		// Check numeric size and precision
		for i := len(assertMeta); i < 2; i++ {
			assertMeta = append(assertMeta, "0")
		}
		for i := len(targetMeta); i < 2; i++ {
			targetMeta = append(targetMeta, "0")
		}

		return baseMatch && cast.ToInt(assertMeta[0]) <= cast.ToInt(targetMeta[0]) && cast.ToInt(assertMeta[1]) <= cast.ToInt(targetMeta[1])
	}

	return baseMatch
}

func (d postgresDialect) ExprHandler(n *ql.ASTNode, args ...exp.Expression) (expr exp.Expression, err error) {
	switch ref := strings.ToLower(n.Ref); ref {
	case "concat":
		return castColumnDataToText("CONCAT", args...)

	case "in":
		return drivers.OpHandlerIn(d, n, args...)

	case "nin":
		return drivers.OpHandlerNotIn(d, n, args...)

	case "like", "nlike":
		if dalType, ok := n.Args[0].Meta["dal.Attribute"].(*dal.Attribute); ok {
			col, err := d.AttributeToColumn(dalType)
			if err != nil {
				return nil, err
			}

			// if the type is id (numeric) data type, then cast it to text
			if col.Type.Name == "NUMERIC" {
				op := "LIKE"
				if ref == "nlike" {
					op = "NOT LIKE"
				}
				return castColumnDataToText(op, args...)
			}
		}
	}

	return ref2exp.RefHandler(n, args...)
}

func (d postgresDialect) ValHandler(n *ql.ASTNode) (out exp.Expression, err error) {
	switch v := n.Value.V.(type) {
	case *expr.Boolean:
		// value handling for boolean is different in postgres
		// this was done to support value parsing for (IS $3) statement
		if cast.ToBool(v.Get()) {
			out = drivers.LiteralTRUE
		} else {
			out = drivers.LiteralFALSE
		}
	default:
		out = exp.NewLiteralExpression("?", n.Value.V.Get())
	}

	return
}

func (d postgresDialect) OrderedExpression(expr exp.Expression, dir exp.SortDirection, nst exp.NullSortType) exp.OrderedExpression {
	return exp.NewOrderedExpression(expr, dir, nst)
}

// castColumnDataToText converts column data to text.
func castColumnDataToText(op string, args ...exp.Expression) (expr exp.Expression, err error) {
	aa := make([]any, len(args))
	for a := range args {
		aa[a] = exp.NewCastExpression(exp.NewLiteralExpression("?", args[a]), "TEXT")
	}

	return exp.NewSQLFunctionExpression(op, aa...), nil
}
