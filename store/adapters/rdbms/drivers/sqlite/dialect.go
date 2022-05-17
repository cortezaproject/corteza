package sqlite

import (
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	dialect struct{}
)

var (
	goquDialectWrapper = goqu.Dialect("sqlite3")

	_ drivers.Dialect = &dialect{}
)

func init() {
	d := sqlite3.DialectOptions()

	// https://github.com/doug-martin/goqu/v9/pull/330
	d.TruncateClause = []byte("DELETE FROM")

	// Overriding vanila SQLite dialect
	goqu.RegisterDialect("sqlite3", d)
}

func Dialect() *dialect {
	return &dialect{}
}

func (dialect) GOQU() goqu.DialectWrapper {
	return goquDialectWrapper
}

func (dialect) DeepIdentJSON(ident exp.IdentifierExpression, pp ...any) (exp.LiteralExpression, error) {
	return drivers.DeepIdentJSON(ident, pp...), nil
}

func (d dialect) TableCodec(m *dal.Model) drivers.TableCodec {
	return drivers.NewTableCodec(m, d)
}

func (d dialect) TypeWrap(t dal.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	return drivers.TypeWrap(t)
}

func (dialect) AttributeCast(attr *dal.Attribute, val exp.LiteralExpression) (exp.LiteralExpression, error) {
	return drivers.AttributeCast(attr, val)
}
