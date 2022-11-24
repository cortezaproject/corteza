package mssql

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"

	"github.com/cortezaproject/corteza/server/pkg/cast2"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/dialect/sqlserver"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	mssqlDialect struct{}
)

var (
	_ drivers.Dialect = &mssqlDialect{}

	dialect            = &mssqlDialect{}
	goquDialectWrapper = goqu.Dialect("sqlserver")
	quoteIdent         = string(sqlserver.DialectOptions().QuoteRune)

	nuances = drivers.Nuances{
		HavingClauseMustUseAlias: true,
		TwoStepUpsert:            true,
	}
)

func init() {
	d := sqlserver.DialectOptions()

	// MSSQL doesn't support the classic boolean constants but we still need to
	// boolean expressions.
	// Disallowing boolean datatype on goqu level is too strict and it prevents most
	// of the queries to work.
	d.BooleanDataTypeSupported = true

	// Use 1/0 as an alternative to booleans
	d.True = []byte("1")
	d.False = []byte("0")

	// d.CastFragment = []byte("TRY_CONVERT")

	// Overriding vanila dialect
	goqu.RegisterDialect("sqlserver", d)
}

func Dialect() *mssqlDialect {
	return dialect
}

func (mssqlDialect) Nuances() drivers.Nuances {
	return nuances
}

func (mssqlDialect) GOQU() goqu.DialectWrapper  { return goquDialectWrapper }
func (mssqlDialect) QuoteIdent(i string) string { return quoteIdent + i + quoteIdent }

func (d mssqlDialect) IndexFieldModifiers(attr *dal.Attribute, mm ...dal.IndexFieldModifier) (string, error) {
	return drivers.IndexFieldModifiers(attr, d.QuoteIdent, mm...)
}

func (d mssqlDialect) JsonQuote(expr exp.Expression) exp.Expression {
	return exp.NewSQLFunctionExpression(
		"JSON_VALUE",
		expr,
		exp.NewLiteralExpression("'$[0]'"),
	)
}

func (d mssqlDialect) JsonExtract(jsonDoc exp.Expression, pp ...any) (path exp.Expression, err error) {
	if path, err = jsonPathExpr(pp...); err != nil {
		return
	} else {
		return exp.NewSQLFunctionExpression("JSON_QUERY", jsonDoc, path), nil
	}
}

func (d mssqlDialect) JsonExtractUnquote(jsonDoc exp.Expression, pp ...any) (path exp.Expression, err error) {
	if path, err = jsonPathExpr(pp...); err != nil {
		return
	} else {
		return exp.NewSQLFunctionExpression("JSON_VALUE", jsonDoc, path), nil
	}
}

// JsonArrayContains prepares sqlserver compatible comparison of value (or ident) and JSON array
//
// # literal value = multi-value field / plain
// # multi-value field = single-value field / plain
// JSON_CONTAINS(v, JSON_EXTRACT(needle, '$.f3'), '$.f2')
//
// # single-value field = multi-value field / plain
// # multi-value field = single-value field / plain
// JSON_CONTAINS(v, '"needle"', '$.f2')
//
// This approach is not optimal, but it is the only way to make it work
func (d mssqlDialect) JsonArrayContains(needle, haystack exp.Expression) (_ exp.Expression, err error) {
	// @todo untested
	return exp.NewLiteralExpression(
		"? IN ?",
		needle,
		exp.NewSQLFunctionExpression("OPENJSON", haystack),
	), nil
	// return exp.NewSQLFunctionExpression(

	// 	"JSON_CONTAINS",
	// 	haystack, needle), nil

	// return exp.NewSQLFunctionExpression("JSON_CONTAINS", haystack, needle), nil
}

func (d mssqlDialect) TableCodec(m *dal.Model) drivers.TableCodec {
	return drivers.NewTableCodec(m, d)
}

func (d mssqlDialect) TypeWrap(dt dal.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	switch c := dt.(type) {
	case *dal.TypeTimestamp:
		return &drivers.TypeTimestamp{&dal.TypeTimestamp{
			Nullable: c.Nullable,

			// sqlserver does not support timezone
			Timezone: false,

			// sqlserver does not support precision
			Precision: 0,
		}}
	}

	return drivers.TypeWrap(dt)
}

// AttributeCast for sqlserver
//
// https://dev.sqlserver.com/doc/refman/8.0/en/cast-functions.html#function_cast
func (mssqlDialect) AttributeCast(attr *dal.Attribute, val exp.Expression) (expr exp.Expression, err error) {

	switch attr.Type.(type) {

	case *dal.TypeText:
		expr = exp.NewCastExpression(val, "VARCHAR(MAX)")

	case *dal.TypeBoolean:
		return val, nil

	default:
		return attributeCast(attr, val)

	}

	return
}

func (mssqlDialect) AttributeToColumn(attr *dal.Attribute) (col *ddl.Column, err error) {
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
		col.Type.Name = "DATETIME"
		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *TypeTime:
		col.Type.Name = "TIME"
		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *dal.TypeDate:
		col.Type.Name = "DATE"
		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *dal.TypeNumber:
		col.Type.Name = "DECIMAL"
		// @todo precision, scale?
		col.Default = ddl.DefaultNumber(t.HasDefault, t.Precision, t.DefaultValue)

	case *dal.TypeText:
		if t.Length > 0 {
			col.Type.Name = fmt.Sprintf("VARCHAR(%d)", t.Length)
		} else {
			col.Type.Name = "VARCHAR(MAX)"
		}

		if t.HasDefault {
			// @todo use proper quote type
			col.Default = fmt.Sprintf("%q", t.DefaultValue)
		}

	case *dal.TypeJSON:
		col.Type.Name = "VARCHAR(MAX)"

	case *dal.TypeGeometry:
		col.Type.Name = "VARCHAR(MAX)"

	case *dal.TypeBlob:
		col.Type.Name = "VARBINARY(MAX)"

	case *dal.TypeBoolean:
		col.Type.Name = "BIT"

	case *dal.TypeUUID:
		col.Type.Name = "CHAR(36)"

	default:
		return nil, fmt.Errorf("unsupported column type: %s ", t.Type())
	}

	return
}

// @todo untested
func (d mssqlDialect) ExprHandler(n *ql.ASTNode, args ...exp.Expression) (expr exp.Expression, err error) {
	switch ref := strings.ToLower(n.Ref); ref {
	case "concat":
		return exp.NewLiteralExpression("?"+strings.Repeat(" || ?", len(args)-1), cast2.Anys(args...)...), nil

	case "in":
		return drivers.OpHandlerIn(d, n, args...)

	case "nin":
		return drivers.OpHandlerNotIn(d, n, args...)

	}

	return ref2exp.RefHandler(n, args...)
}

func (d mssqlDialect) OrderedExpression(expr exp.Expression, dir exp.SortDirection, _ exp.NullSortType) exp.OrderedExpression {
	return exp.NewOrderedExpression(expr, dir, exp.NoNullsSortType)
}

// attributeCast uses mssql's TRY_CONVERT to avoid the need for regex value validation
// since mssql doesn't support regexes, nor can we use the extended LIKE.
func attributeCast(attr *dal.Attribute, val exp.Expression) (exp.Expression, error) {
	switch attr.Type.(type) {
	case *dal.TypeID, *dal.TypeRef:
		return exp.NewLiteralExpression("TRY_CONVERT(BIGINT,?)", val), nil

	case *dal.TypeNumber:
		return exp.NewLiteralExpression("TRY_CONVERT(NUMERIC,?)", val), nil

	case *dal.TypeTimestamp:
		return exp.NewLiteralExpression("TRY_CONVERT(TIMESTAMPTZ,?)", val), nil

	case *dal.TypeDate:
		return exp.NewLiteralExpression("TRY_CONVERT(DATE,?)", val), nil

	case *dal.TypeTime:
		return exp.NewLiteralExpression("TRY_CONVERT(TIMETZ,?)", val), nil

	case *dal.
		TypeBoolean:
		return exp.NewLiteralExpression("TRY_CONVERT(BIT,?)", val), nil

	default:
		return val, nil
	}
}
