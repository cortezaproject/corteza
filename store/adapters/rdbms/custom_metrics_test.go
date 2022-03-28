package rdbms

import (
	"strings"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/require"
)

func Test_timestampStats(t *testing.T) {
	tests := []struct {
		fields []string
		sql    string
	}{
		{
			fields: []string{"a"},
			sql:    `SELECT COUNT(*), SUM(CASE  WHEN ("a_at" IS NOT NULL) THEN 1 ELSE 0 END) AS "a", SUM(CASE  WHEN ("a_at" IS NULL) THEN 1 ELSE 0 END) AS "valid"`,
		},
		{
			fields: []string{"a", "b"},
			sql:    `SELECT COUNT(*), SUM(CASE  WHEN ("a_at" IS NOT NULL) THEN 1 ELSE 0 END) AS "a", SUM(CASE  WHEN ("b_at" IS NOT NULL) THEN 1 ELSE 0 END) AS "b", SUM(CASE  WHEN (("a_at" IS NULL) AND ("b_at" IS NULL)) THEN 1 ELSE 0 END) AS "valid"`,
		},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.fields, "+"), func(t *testing.T) {
			var (
				req         = require.New(t)
				ee          = timestampStatExpr(tt.fields...)
				sql, _, err = goqu.Select(ee...).ToSQL()
			)

			req.NoError(err)
			req.Equal(tt.sql, sql)
		})
	}
}
