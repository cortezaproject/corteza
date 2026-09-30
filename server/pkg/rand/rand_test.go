package rand

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytes(t *testing.T) {
	seen := map[string]bool{}

	for i := 0; i < 100; i++ {
		b := string(Bytes(32))
		require.Len(t, b, 32)
		require.Equal(t, "", strings.Trim(b, letterBytes), "only letters and digits")
		require.False(t, seen[b], "must not repeat")
		seen[b] = true
	}
}

func TestPassword(t *testing.T) {
	for i := 0; i < 100; i++ {
		p := Password(12)
		require.Len(t, p, 12)
		require.True(t, strings.ContainsAny(p, letterDigits), "has a digit")
		require.True(t, strings.ContainsAny(p, letterSpecials), "has a special")
		require.Equal(t, "", strings.Trim(p, letterBytes+letterSpecials))
	}
}

func TestIntn(t *testing.T) {
	hit := map[int]bool{}

	for i := 0; i < 1000; i++ {
		v := Intn(6)
		require.GreaterOrEqual(t, v, 0)
		require.Less(t, v, 6)
		hit[v] = true
	}

	require.Len(t, hit, 6, "every value must come up")
}
