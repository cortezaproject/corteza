package rbac

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResourceMatch(t *testing.T) {
	var (
		tcc = []struct {
			m string
			r string
			e bool
		}{
			{"a:b:c", "a:b:c", true},
			{"a:b:*", "a:b:c", true},
			{"a:*:*", "a:b:c", true},
			{"*:*:*", "a:b:c", true},
			{"a:*:*", "1:2:3", false},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.m, func(t *testing.T) {
			require.Equal(t, tc.e, matchResource(tc.m, tc.r))
		})
	}
}
