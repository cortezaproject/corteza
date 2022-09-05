package postgres

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	postgresDialect struct{}
)

var (
	_ drivers.Dialect = &postgresDialect{}

	dialect            = &postgresDialect{}
	goquDialectWrapper = goqu.Dialect("postgres")
)

func Dialect() *postgresDialect {
	return dialect
}

func (postgresDialect) GOQU() goqu.DialectWrapper {
	return goquDialectWrapper
}

func (postgresDialect) DeepIdentJSON(ident exp.IdentifierExpression, pp ...any) (exp.LiteralExpression, error) {
	return drivers.DeepIdentJSON(ident, pp...), nil
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

func (postgresDialect) AttributeCast(attr *dal.Attribute, val exp.LiteralExpression) (exp.LiteralExpression, error) {
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

func (postgresDialect) NativeColumnType(t dal.Type) (ct *ddl.ColumnType, err error) {
	ct = &ddl.ColumnType{
		Null: t.IsNullable(),
	}

	switch c := t.(type) {
	case *dal.TypeID, *dal.TypeRef:
		ct.Name = "BIGINT"

	case *dal.TypeTimestamp:
		ct.Name = "TIMESTAMP"

		if c.Timezone {
			ct.Name += "TZ"
		}
		ct.Name += fmt.Sprintf("(%d)", c.Precision)

	case *TypeTime:
		ct.Name = "TIME"

		if c.Timezone {
			ct.Name += "TZ"
		}
		ct.Name += fmt.Sprintf("(%d)", c.Precision)

	case *dal.TypeDate:
		ct.Name = "DATE"

	case *dal.TypeNumber:
		ct.Name = "NUMERIC"
		// @todo precision, scale?

	case *dal.TypeText:
		if c.Length > 0 {
			// VARCHAR(0) is useless
			ct.Name = fmt.Sprintf("VARCHAR(%d)", c.Length)
		} else {
			ct.Name = "TEXT"
		}

	case *dal.TypeJSON:
		ct.Name = "JSONB"

	case *dal.TypeGeometry:
		// @todo geometry type
		ct.Name = "JSONB"

	case *dal.TypeBlob:
		ct.Name = "BYTEA"

	case *dal.TypeBoolean:
		ct.Name = "BOOLEAN"

	case *dal.TypeUUID:
		ct.Name = "UUID"

	default:
		return nil, fmt.Errorf("unsupported column type: %s ", c.Type())
	}

	return
}

func (postgresDialect) ExprHandler(n *ql.ASTNode, args ...exp.Expression) (exp.Expression, error) {
	return ref2exp.RefHandler(n, args...)
}
