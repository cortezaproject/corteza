package dal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSensitivityLevelIncluded(t *testing.T) {
	t.Run("does not include; empty", func(t *testing.T) {
		ll := sensitivityLevelIndex{}

		require.False(t, ll.includes(1))
	})

	t.Run("does not include; not in there", func(t *testing.T) {
		ll := SensitivityLevelIndex(SensitivityLevel{
			Handle: "a",
			ID:     0,
			Level:  1,
		})

		require.False(t, ll.includes(1))
	})

	t.Run("includes", func(t *testing.T) {
		ll := SensitivityLevelIndex(SensitivityLevel{
			Handle: "a",
			ID:     1,
			Level:  1,
		})

		require.True(t, ll.includes(1))
	})
}

func TestSensitivityLevelSubset(t *testing.T) {
	commonLevels := SensitivityLevelSet{{
		ID:     1,
		Handle: "a",
		Level:  1,
	}, {
		ID:     2,
		Handle: "b",
		Level:  2,
	}, {
		ID:     3,
		Handle: "c",
		Level:  3,
	}}

	cases := []struct {
		name string
		a    uint64
		b    uint64
		out  bool
	}{{
		name: "lower is zero",
		a:    0,
		b:    1,
		out:  true,
	}, {
		name: "both zero",
		a:    0,
		b:    0,
		out:  true,
	}, {
		name: "upper zero, lower not",
		a:    1,
		b:    0,
		out:  false,
	}, {
		name: "lower less",
		a:    1,
		b:    2,
		out:  true,
	}, {
		name: "equal",
		a:    2,
		b:    2,
		out:  true,
	}, {
		name: "lower higher",
		a:    3,
		b:    2,
		out:  false,
	}}

	ix := SensitivityLevelIndex(commonLevels...)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.out, ix.isSubset(c.a, c.b))
		})
	}
}
