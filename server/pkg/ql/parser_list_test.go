package ql

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_lists(t *testing.T) {
	tcc := []struct {
		in, out string
	}{
		{
			"a IN ('1', '2', '3')",
			`{"ref":"in","args":[{"symbol":"a"},{"ref":"list","args":[{"value":{"@type":"String","@value":"1"}},{"value":{"@type":"String","@value":"2"}},{"value":{"@type":"String","@value":"3"}}]}]}`,
		},
		{
			"a NOT IN (1,2) AND b = 'x'",
			`{"ref":"and","args":[{"ref":"nin","args":[{"symbol":"a"},{"ref":"list","args":[{"value":{"@type":"Integer","@value":1}},{"value":{"@type":"Integer","@value":2}}]}]},{"ref":"eq","args":[{"symbol":"b"},{"value":{"@type":"String","@value":"x"}}]}]}`,
		},
		{
			"(a IN (1 , 2)) OR (b IN ( 3, 4 ))",
			`{"ref":"or","args":[{"ref":"in","args":[{"symbol":"a"},{"ref":"list","args":[{"value":{"@type":"Integer","@value":1}},{"value":{"@type":"Integer","@value":2}}]}]},{"ref":"in","args":[{"symbol":"b"},{"ref":"list","args":[{"value":{"@type":"Integer","@value":3}},{"value":{"@type":"Integer","@value":4}}]}]}]}`,
		},
		{
			// a single parenthesised value stays a plain value
			"a IN ('1')",
			`{"ref":"in","args":[{"symbol":"a"},{"value":{"@type":"String","@value":"1"}}]}`,
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

	_, err := NewParser().Parse("a IN (1, 2")
	require.Error(t, err)
}
