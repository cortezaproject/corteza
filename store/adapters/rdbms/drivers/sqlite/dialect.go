package sqlite

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	sqliteDialect struct{}
)

var (
	_ drivers.Dialect = &sqliteDialect{}

	dialect            = &sqliteDialect{}
	goquDialectWrapper = goqu.Dialect("sqlite3")
)

func init() {
	d := sqlite3.DialectOptions()

	// https://github.com/doug-martin/goqu/v9/pull/330
	d.TruncateClause = []byte("DELETE FROM")

	// Overriding vanila SQLite dialect
	goqu.RegisterDialect("sqlite3", d)
}

func Dialect() *sqliteDialect {
	return dialect
}

func (sqliteDialect) GOQU() goqu.DialectWrapper {
	return goquDialectWrapper
}

func (sqliteDialect) DeepIdentJSON(ident exp.IdentifierExpression, pp ...any) (exp.LiteralExpression, error) {
	return drivers.DeepIdentJSON(ident, pp...), nil
}

func (d sqliteDialect) TableCodec(m *dal.Model) drivers.TableCodec {
	return drivers.NewTableCodec(m, d)
}

func (d sqliteDialect) TypeWrap(t dal.Type) drivers.Type {
	// Any exception to general type-wrap implementation in the drivers package
	// should be placed here
	return drivers.TypeWrap(t)
}

func (sqliteDialect) AttributeCast(attr *dal.Attribute, val exp.LiteralExpression) (exp.LiteralExpression, error) {
	return drivers.AttributeCast(attr, val)
}

func (sqliteDialect) NativeColumnType(t dal.Type) (ct *ddl.ColumnType, err error) {
	ct = &ddl.ColumnType{
		Null: t.IsNullable(),
	}

	switch c := t.(type) {
	case *dal.TypeID, *dal.TypeRef:
		ct.Name = "BIGINT"

	case *dal.TypeTimestamp:
		ct.Name = "TIMESTAMP"

	case *dal.TypeTime:
		ct.Name = "TIMESTAMP"

	case *dal.TypeDate:
		ct.Name = "DATE"

	case *dal.TypeNumber:
		ct.Name = "NUMERIC"
		// @todo precision, scale?

	case *dal.TypeText:
		if c.Length > 0 {
			ct.Name = fmt.Sprintf("VARCHAR(%d)", c.Length)
		} else {
			ct.Name = "TEXT"
		}

	case *dal.TypeJSON:
		ct.Name = "TEXT"

	case *dal.TypeGeometry:
		ct.Name = "TEXT"

	case *dal.TypeBlob:
		ct.Name = "BLOB"

	case *dal.TypeBoolean:
		ct.Name = "BOOLEAN"

	case *dal.TypeUUID:
		ct.Name = "CHAR(36)"

	default:
		return nil, fmt.Errorf("unsupported column type: %s ", c.Type())
	}

	return
}

func (sqliteDialect) ExprHandler(n *ql.ASTNode, args ...exp.Expression) (exp.Expression, error) {
	return ref2exp.RefHandler(n, args...)
}
