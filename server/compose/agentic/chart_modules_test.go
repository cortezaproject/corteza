package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/require"
)

// checkChartReports and checkChartModules split the work: the first can only
// see that a moduleID is non-zero, the second whether it is really there. These
// pin the half that needs no service.
func TestCheckChartReports_ModuleID(t *testing.T) {
	report := func(modID uint64) []*cmpTypes.ChartConfigReport {
		return []*cmpTypes.ChartConfigReport{{
			ModuleID: modID,
			Metrics:  []map[string]any{{"field": "count", "type": "bar"}},
			Dimensions: []map[string]any{
				{"field": "rarity", "modifier": chartNoGrouping},
			},
		}}
	}

	t.Run("a zero moduleID is refused with the fix named", func(t *testing.T) {
		err := checkChartReports(report(0))
		require.Error(t, err)
		require.Contains(t, err.Error(), "compose_module_lookup")
	})

	t.Run("a non-zero moduleID passes this check", func(t *testing.T) {
		require.NoError(t, checkChartReports(report(12345)))
	})
}

// TL;DR: a report with no module to check is left to checkChartReports.
// Example: checkChartModules must not double-report the zero case, nor panic on
// a nil report entry.
func TestCheckChartModules_SkipsWhatItCannotJudge(t *testing.T) {
	// namespaceID is irrelevant: nothing here reaches a service.
	require.NoError(t, checkChartModules(t.Context(), 1, nil))
	require.NoError(t, checkChartModules(t.Context(), 1, []*cmpTypes.ChartConfigReport{nil}))
	require.NoError(t, checkChartModules(t.Context(), 1, []*cmpTypes.ChartConfigReport{{ModuleID: 0}}))
}
