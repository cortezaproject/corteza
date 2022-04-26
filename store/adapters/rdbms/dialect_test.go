package rdbms

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/require"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

// test deep ident expression generator
func Test_DeepIdentJSON(t *testing.T) {
	var (
		pre  = `SELECT `
		post = ` FROM "test"`

		cc = []struct {
			input []interface{}
			sql   string
			args  []interface{}
		}{
			{
				input: []interface{}{"one"},
				sql:   `"one"`,
				args:  []interface{}{},
			},
			{
				input: []interface{}{"one", "two"},
				sql:   `"one"->'two'`,
				args:  []interface{}{},
			},
			{
				input: []interface{}{"one", 2, "three"},
				sql:   `"one"->2->'three'`,
				args:  []interface{}{},
			},
			{
				input: []interface{}{"one", "two", 3},
				sql:   `"one"->'two'->3`,
				args:  []interface{}{},
			},
		}
	)

	for _, c := range cc {
		t.Run(c.sql, func(t *testing.T) {
			var (
				r = require.New(t)
			)

			sql, args, err := goqu.Dialect("postgres").Select(DeepIdentJSON(c.input[0].(string), c.input[1:]...)).From("test").ToSQL()
			r.NoError(err)
			r.Equal(pre+c.sql+post, sql)
			r.Equal(c.args, args)
		})
	}
}

// test deep ident expression generator
func Test_JsonPath(t *testing.T) {

	var (
		cc = []struct {
			input []interface{}
			path  string
		}{
			{
				input: []interface{}{"two"},
				path:  `$.two`,
			},
			{
				input: []interface{}{2, "three"},
				path:  `$[2].three`,
			},
			{
				input: []interface{}{"two", 3},
				path:  `$.two[3]`,
			},
		}
	)

	for _, c := range cc {
		t.Run(c.path, func(t *testing.T) {
			require.Equal(t, c.path, JsonPath(c.input...))
		})
	}
}
