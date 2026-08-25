package toolkit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TL;DR: a required string of nothing but whitespace counts as missing.
// Example: compose_namespace_create with name "   " passed the required check
// and stored a namespace whose name renders as blank — nobody can find it, and
// in a listing it looks like a row that failed to load rather than one someone
// created.
func TestReqStr_Whitespace(t *testing.T) {
	for _, blank := range []string{"", " ", "   ", "\t", "\n", " \t\n "} {
		_, err := ReqStr(map[string]any{"name": blank}, "name")
		require.Error(t, err, "%q must read as missing", blank)
		require.Contains(t, err.Error(), "name is required")
	}

	t.Run("a missing key is still missing", func(t *testing.T) {
		_, err := ReqStr(map[string]any{}, "name")
		require.Error(t, err)
	})

	t.Run("a non-string is still missing", func(t *testing.T) {
		_, err := ReqStr(map[string]any{"name": 42}, "name")
		require.Error(t, err)
	})

	// Rejecting a name is this function's business; tidying one is not, so a
	// value that passes comes back exactly as it was sent.
	t.Run("a real value is returned untouched", func(t *testing.T) {
		for _, in := range []string{"Leads", " Leads ", "a", "  padded name  "} {
			got, err := ReqStr(map[string]any{"name": in}, "name")
			require.NoError(t, err)
			require.Equal(t, in, got)
		}
	})
}
