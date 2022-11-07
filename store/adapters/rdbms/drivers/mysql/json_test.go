package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// test deep ident expression generator
func Test_DeepIdentJSON(t *testing.T) {
	var (
		cc = []struct {
			input []interface{}
			sql   string
		}{
			{
				input: []interface{}{"one"},
				sql:   `?->>"$.one"`,
			},
			{
				input: []interface{}{"one", "two"},
				sql:   `?->>"$.one.two"`,
			},
			{
				input: []interface{}{"one", 2, "three"},
				sql:   `?->>"$.one[2].three"`,
			},
			{
				input: []interface{}{"one", "two", 3},
				sql:   `?->>"$.one.two[3]"`,
			},
		}
	)

	for _, c := range cc {
		t.Run(c.sql, func(t *testing.T) {
			var (
				r      = require.New(t)
				l, err = JSONPath(nil, c.input...)
			)

			r.NoError(err)
			r.Equal(c.sql, l.Literal())
		})
	}
}
