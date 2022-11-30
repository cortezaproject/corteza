package mysql

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/doug-martin/goqu/v9/sqlgen"
)

type (
	mysqlDialect struct{}
)

var (
	_ drivers.Dialect = &mysqlDialect{}

	dialect            = &mysqlDialect{}
	goquDialectWrapper = goqu.Dialect("mysql")
	goquDialectOptions = mysql.DialectOptions()
	quoteIdent         = string(mysql.DialectOptions().QuoteRune)

	nuances = drivers.Nuances{
		HavingClauseMustUseAlias: false,
	}
)

func Dialect() *mysqlDialect {
	return dialect
}

func (mysqlDialect) Nuances() drivers.Nuances {
	return nuances
}

func (mysqlDialect) GOQU() goqu.DialectWrapper                 { return goquDialectWrapper }
func (mysqlDialect) DialectOptions() *sqlgen.SQLDialectOptions { return goquDialectOptions }
func (mysqlDialect) QuoteIdent(i string) string                { return quoteIdent + i + quoteIdent }

func (d mysqlDialect) IndexFieldModifiers(attr *dal.Attribute, mm ...dal.IndexFieldModifier) (string, error) {
	return drivers.IndexFieldModifiers(attr, d.QuoteIdent, mm...)
}

func (d mysqlDialect) JsonQuote(expr exp.Expression) exp.Expression {
	return exp.NewSQLFunctionExpression(
		"JSON_EXTRACT",
		exp.NewSQLFunctionExpression("JSON_ARRAY", expr),
		exp.NewLiteralExpression("'$[0]'"),
	)
}

func (d mysqlDialect) JsonExtract(jsonDoc exp.Expression, pp ...any) (path exp.Expression, err error) {
	if path, err = jsonPathExpr(pp...); err != nil {
		return
	} else {
		return exp.NewSQLFunctionExpression("JSON_EXTRACT", jsonDoc, path), nil
	}
}

func (d mysqlDialect) JsonExtractUnquote(jsonDoc exp.Expression, pp ...any) (_ exp.Expression, err error) {
	if jsonDoc, err = d.JsonExtract(jsonDoc, pp...); err != nil {
		return
	} else {
		return exp.NewSQLFunctionExpression("JSON_UNQUOTE", jsonDoc), nil
	}
}

// JsonArrayContains prepares MySQL compatible comparison of value (or ident) and JSON array
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
func (d mysqlDialect) JsonArrayContains(needle, haystack exp.Expression) (_ exp.Expression, err error) {
	return exp.NewSQLFunctionExpression("JSON_CONTAINS", haystack, needle), nil
}

func (d mysqlDialect) TableCodec(m *dal.Model) drivers.TableCodec {
	return drivers.NewTableCodec(m, d)
}

func (d mysqlDialect) TypeWrap(dt dal.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	switch c := dt.(type) {
	case *dal.TypeTimestamp:
		return &drivers.TypeTimestamp{&dal.TypeTimestamp{
			Nullable: c.Nullable,

			// mysql does not support timezone
			Timezone: false,

			// mysql does not support precision
			Precision: 0,
		}}
	}

	return drivers.TypeWrap(dt)
}

// AttributeCast for mySQL
//
// https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html#function_cast
func (mysqlDialect) AttributeCast(attr *dal.Attribute, val exp.Expression) (exp.Expression, error) {
	var (
		c exp.CastExpression
	)

	switch attr.Type.(type) {

	case *dal.TypeNumber:
		ce := exp.NewCaseExpression().
			When(drivers.RegexpLike(drivers.CheckNumber, val), val).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "DECIMAL(65,10)")

	case *dal.TypeTimestamp:
		ce := exp.NewCaseExpression().
			When(drivers.RegexpLike(drivers.CheckFullISO8061, val), val).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "DATETIME")

	case *dal.TypeBoolean:
		c = exp.NewCastExpression(drivers.BooleanCheck(val), "SIGNED")

	case *dal.TypeID, *dal.TypeRef:
		ce := exp.NewCaseExpression().
			When(drivers.RegexpLike(drivers.CheckID, val), val).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "UNSIGNED")

	case *dal.TypeTime:
		ce := exp.NewCaseExpression().
			When(drivers.RegexpLike(drivers.CheckTimeISO8061, val), val).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "TIME")

	default:
		return drivers.AttributeCast(attr, val)

	}

	return exp.NewLiteralExpression("?", c), nil
}

func (mysqlDialect) AttributeToColumn(attr *dal.Attribute) (col *ddl.Column, err error) {
	col = &ddl.Column{
		Ident:   attr.StoreIdent(),
		Comment: attr.Label,
		Type: &ddl.ColumnType{
			Null: attr.Type.IsNullable(),
		},
	}

	switch t := attr.Type.(type) {
	case *dal.TypeID:
		col.Type.Name = "BIGINT UNSIGNED"
		col.Default = ddl.DefaultID(t.HasDefault, t.DefaultValue)
	case *dal.TypeRef:
		col.Type.Name = "BIGINT UNSIGNED"
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
			col.Type.Name = "TEXT"
		}

		if t.HasDefault {
			// @todo use proper quote type
			col.Default = fmt.Sprintf("%q", t.DefaultValue)
		}

	case *dal.TypeJSON:
		col.Type.Name = "JSON"

	case *dal.TypeGeometry:
		col.Type.Name = "JSON"

	case *dal.TypeBlob:
		col.Type.Name = "BLOB"

	case *dal.TypeBoolean:
		col.Type.Name = "TINYINT(1)"

	case *dal.TypeUUID:
		col.Type.Name = "CHAR(36)"

	case *dal.TypeEnum:
		col.Type.Name = "TEXT"

	default:
		return nil, fmt.Errorf("unsupported column type: %s ", t.Type())
	}

	return
}

func (d mysqlDialect) ExprHandler(n *ql.ASTNode, args ...exp.Expression) (expr exp.Expression, err error) {
	switch ref := strings.ToLower(n.Ref); ref {
	case "in":
		return drivers.OpHandlerIn(d, n, args...)

	case "nin":
		return drivers.OpHandlerNotIn(d, n, args...)
	}

	return ql.DefaultRefHandler(n, args...)
}

func (d mysqlDialect) OrderedExpression(expr exp.Expression, dir exp.SortDirection, _ exp.NullSortType) exp.OrderedExpression {
	return exp.NewOrderedExpression(expr, dir, exp.NoNullsSortType)
}
