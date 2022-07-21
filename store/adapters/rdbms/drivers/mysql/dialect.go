package mysql

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	mysqlDialect struct{}
)

var (
	_ drivers.Dialect = &mysqlDialect{}

	dialect            = &mysqlDialect{}
	goquDialectWrapper = goqu.Dialect("mysql")
)

func Dialect() *mysqlDialect {
	return dialect
}

func (mysqlDialect) GOQU() goqu.DialectWrapper {
	return goquDialectWrapper
}

func (mysqlDialect) DeepIdentJSON(ident exp.IdentifierExpression, pp ...any) (exp.LiteralExpression, error) {
	return JSONPath(ident, pp...)
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
func (mysqlDialect) AttributeCast(attr *dal.Attribute, val exp.LiteralExpression) (exp.LiteralExpression, error) {
	var (
		c exp.CastExpression
	)

	switch attr.Type.(type) {

	case *dal.TypeNumber:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(drivers.CheckNumber), val).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "DECIMAL(65,10)")

	case *dal.TypeTimestamp:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(drivers.CheckFullISO8061), val).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "DATETIME")

	case *dal.TypeBoolean:
		ce := exp.NewCaseExpression().
			When(val.In(drivers.LiteralTRUE, exp.NewLiteralExpression(`'true'`)), drivers.LiteralTRUE).
			When(val.In(drivers.LiteralFALSE, exp.NewLiteralExpression(`'false'`)), drivers.LiteralFALSE).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "SIGNED")

	default:
		return drivers.AttributeCast(attr, val)

	}

	return exp.NewLiteralExpression("?", c), nil
}

func (mysqlDialect) NativeColumnType(ct ddl.ColumnType) string {
	return columnTypeTranslator(ct)
}

func (mysqlDialect) ExprHandler(n *ql.ASTNode, args ...exp.Expression) (exp.Expression, error) {
	return ql.DefaultRefHandler(n, args...)
}

func JSONPath(ident exp.IdentifierExpression, pp ...any) (exp.LiteralExpression, error) {
	var (
		sql strings.Builder
	)

	sql.WriteString(`?->>'$`)

	for _, p := range pp {
		switch path := p.(type) {
		case string:
			sql.WriteString(".")
			sql.WriteString(strings.ReplaceAll(path, "'", `\'`))
		case int:
			sql.WriteString("[")
			sql.WriteString(strconv.Itoa(path))
			sql.WriteString("]")
		default:
			return nil, fmt.Errorf("unexpected path part (%q) type: %T", p, p)
		}
	}

	sql.WriteString(`'`)
	return exp.NewLiteralExpression(sql.String(), ident), nil
}
