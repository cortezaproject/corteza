package postgres

import (
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
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
	goquDialectWrapper = goqu.Dialect("postgres")
)

func Dialect() *mysqlDialect {
	return dialect
}

func (mysqlDialect) GOQU() goqu.DialectWrapper {
	return goquDialectWrapper
}

func (mysqlDialect) DeepIdentJSON(ident exp.IdentifierExpression, pp ...any) (exp.LiteralExpression, error) {
	return drivers.DeepIdentJSON(ident, pp...), nil
}

func (d mysqlDialect) TableCodec(m *dal.Model) drivers.TableCodec {
	return drivers.NewTableCodec(m, d)
}

func (d mysqlDialect) TypeWrap(dt dal.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	switch c := dt.(type) {
	case *dal.TypeTime:
		return &TypeTime{c}
	}

	return drivers.TypeWrap(dt)
}

func (mysqlDialect) AttributeCast(attr *dal.Attribute, val exp.LiteralExpression) (exp.LiteralExpression, error) {
	var (
		c exp.CastExpression
	)

	switch attr.Type.(type) {
	case *dal.TypeBoolean:
		// we need to be strictly dealing with strings here!
		// 1) postgresql's JSON op ->> (last one) returns any JSON value as string
		//    so booleans are casted to 'true' & 'false'
		// 2) postgresql will complain about true == 'true' expressions
		ce := exp.NewCaseExpression().
			When(val.In(exp.NewLiteralExpression(`'true'`)), drivers.LiteralTRUE).
			When(val.In(exp.NewLiteralExpression(`'false'`)), drivers.LiteralFALSE).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "BOOLEAN")

	default:
		return drivers.AttributeCast(attr, val)

	}

	return exp.NewLiteralExpression("?", c), nil
}

func (mysqlDialect) NativeColumnType(ct ddl.ColumnType) string {
	return ddl.ColumnTypeTranslator(ct)
}
