package rbac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermuteResource(t *testing.T) {
	tcc := []struct {
		in  string
		out []string
	}{{
		in:  "xx/1/2",
		out: []string{"xx/1/2", "xx/1/*", "xx/*/*"},
	}, {
		in:  "xx/1/*",
		out: []string{"xx/1/*", "xx/*/*"},
	}, {
		in:  "xx/*/*",
		out: []string{"xx/*/*"},
	}, {
		in:  "xx",
		out: []string{"xx"},
	}}

	req := require.New(t)
	for _, tc := range tcc {
		t.Run(tc.in, func(t *testing.T) {
			req.Equal(tc.out, permuteResource(tc.in))
		})
	}
}
