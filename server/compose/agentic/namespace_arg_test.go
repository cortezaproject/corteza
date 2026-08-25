package agentic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TL;DR: a missing "namespace" argument says so.
// Example: the argument was handed straight to FindByAny, so an absent one came
// back as "namespace lookup failed: invalid ID" — which reads as though a
// namespace was supplied and rejected, and invites a model to invent an ID
// rather than pass the namespace it forgot.
func TestLookupNamespaceArg_MissingArgument(t *testing.T) {
	for name, args := range map[string]map[string]any{
		"absent":          {},
		"nil":             {"namespace": nil},
		"empty string":    {"namespace": ""},
		"other args only": {"module": "card", "recordID": "1"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := lookupNamespaceArg(t.Context(), args)
			require.Error(t, err)
			require.Contains(t, err.Error(), "namespace is required")
			require.NotContains(t, err.Error(), "invalid ID")
		})
	}
}
