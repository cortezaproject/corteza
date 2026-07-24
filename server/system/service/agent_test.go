package service

import (
	"testing"

	"github.com/crusttech/human/server/system/agentic/tcl"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestPrepareTCL(t *testing.T) {
	enabled, disabled := true, false

	t.Run("nil flag leaves behavior untouched", func(t *testing.T) {
		b := types.AgentBehavior{}
		prepareTCL(&b)
		require.Nil(t, b.TreatyCLEnabled)
		require.Empty(t, b.TreatyCLArticles)
	})

	t.Run("explicit false leaves articles untouched", func(t *testing.T) {
		b := types.AgentBehavior{TreatyCLEnabled: &disabled, TreatyCLArticles: []string{"custom"}}
		prepareTCL(&b)
		require.Equal(t, []string{"custom"}, b.TreatyCLArticles)
	})

	t.Run("enabled with no articles populates defaults", func(t *testing.T) {
		b := types.AgentBehavior{TreatyCLEnabled: &enabled}
		prepareTCL(&b)
		require.Equal(t, tcl.DefaultArticleIDs(), b.TreatyCLArticles)
	})

	t.Run("enabled with articles merges hardwired back in", func(t *testing.T) {
		b := types.AgentBehavior{TreatyCLEnabled: &enabled, TreatyCLArticles: []string{"custom"}}
		prepareTCL(&b)
		require.Equal(t, tcl.MergeWithHardwired([]string{"custom"}), b.TreatyCLArticles)
	})
}
