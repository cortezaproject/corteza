package postgres

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

// test deep ident expression generator
func Test_DeepIdentJSON(t *testing.T) {
	var (
		pre  = `SELECT `
		post = ` FROM "test"`

		cc = []struct {
			input  []interface{}
			sql    string
			asJSON bool
			args   []interface{}
		}{
			{
				input: []interface{}{"one"},
				sql:   `"one"`,
				args:  []interface{}{},
			},
			{
				input: []interface{}{"one", "two"},
				sql:   `"one"->>'two'`,
				args:  []interface{}{},
			},
			{
				input: []interface{}{"one", 2, "three"},
				sql:   `"one"->2->>'three'`,
				args:  []interface{}{},
			},
			{
				input: []interface{}{"one", "two", 3},
				sql:   `"one"->'two'->>3`,
				args:  []interface{}{},
			},
			{
				input:  []interface{}{"one"},
				asJSON: true,
				sql:    `"one"`,
				args:   []interface{}{},
			},
			{
				input:  []interface{}{"one", "two"},
				asJSON: true,
				sql:    `"one"->'two'`,
				args:   []interface{}{},
			},
			{
				input:  []interface{}{"one", 2, "three"},
				asJSON: true,
				sql:    `"one"->2->'three'`,
				args:   []interface{}{},
			},
			{
				input:  []interface{}{"one", "two", 3},
				asJSON: true,
				sql:    `"one"->'two'->3`,
				args:   []interface{}{},
			},
		}

		conv = func(asJSON bool, pp ...any) exp.SQLExpression {
			return goqu.Dialect("postgres").
				Select(DeepIdentJSON(asJSON, exp.ParseIdentifier(pp[0].(string)), pp[1:]...)).From("test")
		}
	)

	for _, c := range cc {
		t.Run(c.sql, func(t *testing.T) {
			var (
				r = require.New(t)
			)

			sql, args, err := conv(c.asJSON, c.input...).ToSQL()
			r.NoError(err)
			r.Equal(pre+c.sql+post, sql)
			r.Equal(c.args, args)
		})
	}
}
