package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sanitizeLink(t *testing.T) {
	type (
		tt struct {
			name   string
			link   string
			expect string
		}
	)

	tcc := []tt{
		{
			name:   `empty link`,
			link:   ``,
			expect: `//`,
		},
		{
			name:   `Example URL with query`,
			link:   `https://example.url/query`,
			expect: `//example.url/query`,
		},
		{
			name:   `URL with additional js`,
			link:   `javascript:window.alert('foobar')`,
			expect: `//javascript:window.alert(foobar)`,
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			require.New(t).Equal(tc.expect, sanitizeLink(tc.link))
		})
	}
}
