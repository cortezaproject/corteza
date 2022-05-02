package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	dialect struct{}
)

var (
	goquDialectWrapper = goqu.Dialect("mysql")

	_ drivers.Dialect = &dialect{}
)

func Dialect() *dialect {
	return &dialect{}
}

func (dialect) GOQU() goqu.DialectWrapper {
	return goquDialectWrapper
}

func (dialect) DeepIdentJSON(ident exp.IdentifierExpression, pp ...any) (exp.LiteralExpression, error) {
	return JSONPath(ident, pp...)
}

func (d dialect) TableCodec(m *data.Model) drivers.TableCodec {
	return drivers.NewTableCodec(m, d)
}

func (d dialect) TypeWrap(dt data.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	switch c := dt.(type) {
	case *data.TypeTimestamp:
		return &drivers.TypeTimestamp{&data.TypeTimestamp{
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
func (dialect) AttributeCast(attr *data.Attribute, val exp.LiteralExpression) (exp.LiteralExpression, error) {
	var (
		c exp.CastExpression
	)

	switch attr.Type.(type) {

	case *data.TypeNumber:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(drivers.CheckNumber), val).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "DECIMAL(65,10)")

	case *data.TypeTimestamp:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(drivers.CheckFullISO8061), val).
			Else(drivers.LiteralNULL)

		c = exp.NewCastExpression(ce, "DATETIME")

	case *data.TypeBoolean:
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

func JSONPath(ident exp.IdentifierExpression, pp ...any) (exp.LiteralExpression, error) {
	var (
		sql strings.Builder
	)

	sql.WriteString(`?->>"$`)

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

	sql.WriteString(`"`)

	return exp.NewLiteralExpression(sql.String(), ident), nil
}
