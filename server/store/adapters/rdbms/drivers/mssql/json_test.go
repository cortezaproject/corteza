package mssql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// test JSON path generator
func Test_jsonPath(t *testing.T) {
	var (
		cc = []struct {
			input []any
			path  string
		}{
			{
				input: []any{"one"},
				path:  `$.one`,
			},
			{
				input: []any{"one", 2, "three"},
				path:  `$.one[2].three`,
			},
			{
				input: []any{"tw'o"},
				path:  `$.tw''o`,
			},
		}
	)

	for _, c := range cc {
		t.Run(c.path, func(t *testing.T) {
			var (
				r = require.New(t)
			)

			path, err := jsonPath(c.input...)
			r.NoError(err)
			r.Equal(c.path, path)
		})
	}
}
