package sqlite

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/cast2"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/doug-martin/goqu/v9/sqlgen"
)

type (
	sqliteDialect struct{}
)

var (
	_ drivers.Dialect = &sqliteDialect{}

	dialect            = &sqliteDialect{}
	goquDialectWrapper = goqu.Dialect("sqlite3")
	goquDialectOptions *sqlgen.SQLDialectOptions
	quoteIdent         = string(sqlite3.DialectOptions().QuoteRune)

	nuances = drivers.Nuances{
		HavingClauseMustUseAlias: true,
	}
)

func init() {
	goquDialectOptions = sqlite3.DialectOptions()

	// https://github.com/doug-martin/goqu/v9/pull/330
	goquDialectOptions.TruncateClause = []byte("DELETE FROM")

	// Overriding vanila SQLite dialect
	goqu.RegisterDialect("sqlite3", goquDialectOptions)
}

func Dialect() *sqliteDialect {
	return dialect
}

func (sqliteDialect) Nuances() drivers.Nuances {
	return nuances
}

func (sqliteDialect) GOQU() goqu.DialectWrapper                 { return goquDialectWrapper }
func (sqliteDialect) DialectOptions() *sqlgen.SQLDialectOptions { return goquDialectOptions }
func (sqliteDialect) QuoteIdent(i string) string                { return quoteIdent + i + quoteIdent }

func (d sqliteDialect) IndexFieldModifiers(attr *dal.Attribute, mm ...dal.IndexFieldModifier) (string, error) {
	return drivers.IndexFieldModifiers(attr, d.QuoteIdent, mm...)
}

func (d sqliteDialect) JsonQuote(expr exp.Expression) exp.Expression {
	return exp.NewSQLFunctionExpression("JSON_QUOTE", expr)
}

func (sqliteDialect) JsonExtract(ident exp.Expression, pp ...any) (exp.Expression, error) {
	return DeepIdentJSON(true, ident, pp...), nil
}

func (sqliteDialect) JsonExtractUnquote(ident exp.Expression, pp ...any) (exp.Expression, error) {
	return DeepIdentJSON(false, ident, pp...), nil
}

// JsonArrayContains prepares SQLite compatible comparison of value and JSON array
//
// # literal value = multi-value field / plain
// # multi-value field = single-value field / plain
//
// 'aaa' in (select value from json_each(v->'f2'))
//
// # single-value field = multi-value field / plain
// # multi-value field = single-value field / plain
// json_extract(v, '$.f3[0]') in (select value from json_each(v->'f2'));
//
// Unfortunately SQLite converts boolean values into 0 and 1 when decoding from
// JSON and we need a special handler for that.
func (sqliteDialect) JsonArrayContains(needle, haystack exp.Expression) (exp.Expression, error) {
	// @todo should be implemented using native SQLite capabilties and
	//       not through custom JSON_ARRAY_CONTAINS function
	return exp.NewLiteralExpression("JSON_ARRAY_CONTAINS(?, ?)", needle, haystack), nil
}

func (d sqliteDialect) TableCodec(m *dal.Model) drivers.TableCodec {
	return drivers.NewTableCodec(m, d)
}

func (d sqliteDialect) TypeWrap(dt dal.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here

	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	switch c := dt.(type) {
	case *dal.TypeDate:
		return &TypeDate{c}
	case *dal.TypeTimestamp:
		return &TypeTimestamp{c}
	}

	return drivers.TypeWrap(dt)
}

func (sqliteDialect) AttributeCast(attr *dal.Attribute, val exp.Expression) (exp.Expression, error) {
	var (
		c exp.Expression
	)

	switch t := attr.Type.(type) {
	case *dal.TypeTimestamp:
		// we have way of handling TZ here (not sure how to do that with strftime() function)
		// if we use %Z at the end strftime will always return NULL
		if t.Precision == 0 {
			c = exp.NewSQLFunctionExpression("strftime", "%Y-%m-%dT%H:%M:%S", val)
		} else {
			c = exp.NewSQLFunctionExpression("strftime", "%Y-%m-%dT%H:%M:%f", val)
		}

	case *dal.TypeTime:
		// we have way of handling TZ here (not sure how to do that with strftime() function)
		// if we use %Z at the end strftime will always return NULL
		if t.Precision == 0 {
			c = exp.NewSQLFunctionExpression("strftime", "%H:%M:%S", val)
		} else {
			c = exp.NewSQLFunctionExpression("strftime", "%H:%M:%f", val)
		}

	case *dal.TypeDate:
		c = exp.NewSQLFunctionExpression("strftime", "%Y-%m-%d", val)

	case *dal.TypeNumber:
		c = exp.NewCastExpression(val, "NUMERIC")

	default:
		return drivers.AttributeCast(attr, val)
	}

	return exp.NewLiteralExpression("?", c), nil

}

func (sqliteDialect) AttributeToColumn(attr *dal.Attribute) (col *ddl.Column, err error) {

	col = &ddl.Column{
		Ident:   attr.StoreIdent(),
		Comment: attr.Label,
		Type: &ddl.ColumnType{
			Null: attr.Type.IsNullable(),
		},
	}

	switch t := attr.Type.(type) {
	case *dal.TypeID, *dal.TypeRef:
		col.Type.Name = "BIGINT"

	case *dal.TypeTimestamp:
		col.Type.Name = "TIMESTAMP"
		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *dal.TypeTime:
		col.Type.Name = "TIMESTAMP"
		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *dal.TypeDate:
		col.Type.Name = "TEXT"
		col.Default = ddl.DefaultValueCurrentTimestamp(t.DefaultCurrentTimestamp)

	case *dal.TypeNumber:
		col.Type.Name = "NUMERIC"
		// @todo precision, scale? x
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
		col.Type.Name = "TEXT"
		if col.Default, err = ddl.DefaultJSON(t.HasDefault, t.DefaultValue); err != nil {
			return nil, err
		}

	case *dal.TypeGeometry:
		col.Type.Name = "TEXT"

	case *dal.TypeBlob:
		col.Type.Name = "BLOB"

	case *dal.TypeBoolean:
		col.Type.Name = "BOOLEAN"
		col.Default = ddl.DefaultBoolean(t.HasDefault, t.DefaultValue)

	case *dal.TypeUUID:
		col.Type.Name = "CHAR(36)"

	case *dal.TypeEnum:
		col.Type.Name = "TEXT"
	default:
		return nil, fmt.Errorf("unsupported column type: %s ", t.Type())
	}

	return
}

func (d sqliteDialect) ExprHandler(n *ql.ASTNode, args ...exp.Expression) (expr exp.Expression, err error) {
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

func (d sqliteDialect) OrderedExpression(expr exp.Expression, dir exp.SortDirection, nst exp.NullSortType) exp.OrderedExpression {
	return exp.NewOrderedExpression(expr, dir, exp.NoNullsSortType)
}
