package sqlite

import (
	"github.com/doug-martin/goqu/v9"
	dialect "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

func init() {
	d := dialect.DialectOptions()

	// https://github.com/doug-martin/goqu/v9/pull/330
	d.TruncateClause = []byte("DELETE FROM")

	// Overriding vanila SQLite dialect
	goqu.RegisterDialect("sqlite3", d)
}
