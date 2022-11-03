package tests

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeepIdentJSON(t *testing.T) {
	const (
		tbl = "test_json_path_test"
	)

	var (
		req = require.New(t)

		count = func(val string) int {
			var (
				out = struct {
					Count int `db:"count"`
				}{}

				diJSON exp.Expression
				err    error
			)
			diJSON, err = conn.dialect.DeepIdentJSON(
				exp.NewIdentifierExpression("", tbl, "c"),
				"a", "b", "c",
			)
			req.NoError(err)

			query := conn.dialect.GOQU().
				Select(goqu.COUNT(goqu.Star()).As("count")).
				From(tbl).
				Where(exp.NewLiteralExpression("?", diJSON).Eq(val))

			err = conn.store.QueryOne(ctx, query, &out)
			req.NoError(err)

			return out.Count
		}
	)

	err := makeTableWithJsonColumn(tbl)
	if err != nil {
		t.Fatalf("can not create table: %v", err)
	}

	insert := conn.dialect.GOQU().
		Insert(tbl).
		Cols("c").
		Vals([]any{`{"a": {"b": {"c": "match"}}}`})

	req.NoError(conn.store.Exec(ctx, insert))

	req.Equal(1, count("match"))
	req.Equal(0, count("nope"))
}

func makeTableWithJsonColumn(tbl string) (err error) {
	if err = exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tbl)); err != nil {
		return
	}

	switch {
	case conn.isSQLite, conn.isPostgres:
		return exec(fmt.Sprintf(`CREATE TABLE %s (c JSONB)`, tbl))

	case conn.isMySQL:
		return exec(fmt.Sprintf(`CREATE TABLE %s (c TEXT)`, tbl))

	default:
		return fmt.Errorf("unsupported driver: %q", conn.config.DriverName)
	}
}
