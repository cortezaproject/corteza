package mcpkit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNearestArgument(t *testing.T) {
	names := []string{"filter", "limit", "module", "namespace", "pageCursor", "recordID", "sort"}

	cases := map[string]string{
		"fliter":     "filter", // transposed
		"Filter":     "filter", // case only
		"filte":      "filter", // truncated
		"pagecursor": "pageCursor",
		"lmit":       "limit",
		"zzzzzzzzzz": "", // nothing close enough to guess at
	}

	for given, want := range cases {
		t.Run(given, func(t *testing.T) {
			require.Equal(t, want, nearestArgument(given, names))
		})
	}
}

func TestUnknownArgumentError(t *testing.T) {
	declared := map[string]any{"filter": nil, "limit": nil, "sort": nil}

	t.Run("one unknown name is quoted and the near match offered", func(t *testing.T) {
		err := unknownArgumentError("compose_record_lookup", []string{"fliter"}, declared)
		require.Contains(t, err.Error(), `does not take a parameter named "fliter"`)
		require.Contains(t, err.Error(), `did you mean "filter"`)
		require.Contains(t, err.Error(), "It takes: filter, limit, sort")
	})

	t.Run("several unknown names are reported together", func(t *testing.T) {
		err := unknownArgumentError("t", []string{"aaa", "bbb"}, declared)
		require.Contains(t, err.Error(), `does not take the parameters "aaa", "bbb"`)
	})

	t.Run("no near match means no suggestion", func(t *testing.T) {
		err := unknownArgumentError("t", []string{"zzzzzzzzzz"}, declared)
		require.NotContains(t, err.Error(), "did you mean")
	})
}

func TestEditDistance(t *testing.T) {
	require.Equal(t, 0, editDistance("filter", "filter"))
	require.Equal(t, 1, editDistance("filte", "filter"))
	require.Equal(t, 2, editDistance("fliter", "filter"))
	require.Equal(t, 6, editDistance("", "filter"))
}
