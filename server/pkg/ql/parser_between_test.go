package ql

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_betweenAnd(t *testing.T) {
	tcc := []struct {
		in, out string
	}{
		{
			"a BETWEEN 1 AND 10",
			`{"ref":"between","args":[{"symbol":"a"},{"value":{"@type":"Integer","@value":1}},{"value":{"@type":"Integer","@value":10}}]}`,
		},
		{
			"a NOT BETWEEN 1 AND 10 AND b = 2",
			`{"ref":"and","args":[{"ref":"nbetween","args":[{"symbol":"a"},{"value":{"@type":"Integer","@value":1}},{"value":{"@type":"Integer","@value":10}}]},{"ref":"eq","args":[{"symbol":"b"},{"value":{"@type":"Integer","@value":2}}]}]}`,
		},
		{
			"b = 2 AND a BETWEEN 1 AND 10",
			`{"ref":"and","args":[{"ref":"eq","args":[{"symbol":"b"},{"value":{"@type":"Integer","@value":2}}]},{"ref":"between","args":[{"symbol":"a"},{"value":{"@type":"Integer","@value":1}},{"value":{"@type":"Integer","@value":10}}]}]}`,
		},
		{
			// the form the webapp writes stays as it was
			"(a BETWEEN 1 10)",
			`{"ref":"between","args":[{"symbol":"a"},{"value":{"@type":"Integer","@value":1}},{"value":{"@type":"Integer","@value":10}}]}`,
		},
	}

	for _, tc := range tcc {
		t.Run(tc.in, func(t *testing.T) {
			n, err := NewParser().Parse(tc.in)
			require.NoError(t, err)

			out, err := json.Marshal(n)
			require.NoError(t, err)
			require.JSONEq(t, tc.out, string(out))
		})
	}
}
