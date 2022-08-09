package tests

import (
	"context"
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

	eachDB(t, func(t *testing.T, c *conn) error {
		var (
			req = require.New(t)
			ctx = context.Background()

			count = func(val string) int {
				var (
					out = struct {
						Count int `db:"count"`
					}{}

					diJSON exp.Expression
					err    error
				)
				diJSON, err = c.dialect.DeepIdentJSON(
					exp.NewIdentifierExpression("", tbl, "c"),
					"a", "b", "c",
				)
				req.NoError(err)

				query := c.dialect.GOQU().
					Select(goqu.COUNT(goqu.Star()).As("count")).
					From(tbl).
					Where(exp.NewLiteralExpression("?", diJSON).Eq(val))

				err = c.store.QueryOne(ctx, query, &out)
				req.NoError(err)

				return out.Count
			}
		)

		err := makeTableWithJsonColumn(c, tbl)
		if err != nil {
			t.Fatalf("can not create table: %v", err)
		}

		insert := c.dialect.GOQU().
			Insert(tbl).
			Cols("c").
			Vals([]any{`{"a": {"b": {"c": "match"}}}`})

		if err = c.store.Exec(ctx, insert); err != nil {
			return err
		}

		req.Equal(1, count("match"))
		req.Equal(0, count("nope"))

		return nil
	})
}

func makeTableWithJsonColumn(c *conn, tbl string) (err error) {
	if err = exec(c, fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tbl)); err != nil {
		return
	}

	switch {
	case c.isSQLite, c.isPostgres:
		return exec(c, fmt.Sprintf(`CREATE TABLE %s (c JSONB)`, tbl))

	case c.isMySQL:
		return exec(c, fmt.Sprintf(`CREATE TABLE %s (c TEXT)`, tbl))

	default:
		return fmt.Errorf("unsupported driver: %q", c.config.DriverName)
	}
}
