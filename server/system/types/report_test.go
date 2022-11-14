package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQLWrapParsing(t *testing.T) {
	tcc := []struct {
		name string
		in   string
		out  string
		err  bool
	}{
		{
			name: "ast node valid",
			in:   `{"ref": "or", "args": [{"value": {"@type": "String","@value": "Maria"}}, {"value": {"@type": "String","@value": "Maria"}}]}`,
			out:  `or("Maria", "Maria")`,
		},
		{
			name: "ast raw expr valid",
			in:   `{"raw": "a || b"}`,
			out:  `or(a, b)`,
		},
		{
			name: "raw expr valid",
			in:   `"c || d"`,
			out:  `or(c, d)`,
		},
		{
			name: "ast node fnc valid",
			in:   `{"ref": "year", "args": [{"ref": "now"}]}`,
			out:  `year(now())`,
		},
		{
			name: "expr fnc valid",
			in:   `"year(now())"`,
			out:  `year(now())`,
		},
		{
			name: "raw expr empty",
			in:   `""`,
			// this is what Stringer returns
			out: `<nil>`,
		},
		{
			name: "ast node empty",
			in:   `{}`,
			// this is what Stringer returns
			out: `<nil>`,
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			w := &ReportFilterExpr{}

			err := json.Unmarshal([]byte(tc.in), &w)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.out, w.ASTNode.String())
			}
		})
	}
}
