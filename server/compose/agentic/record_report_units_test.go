package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/require"
)

func unitsModule() *cmpTypes.Module {
	return &cmpTypes.Module{
		Fields: cmpTypes.ModuleFieldSet{
			{Name: "line_value", Label: "Total value", Kind: "Number",
				Options: cmpTypes.ModuleFieldOptions{"prefix": "$ ", "format": "0,0.00"}},
			{Name: "change_pct", Label: "Change %", Kind: "Number",
				Options: cmpTypes.ModuleFieldOptions{"suffix": " %"}},
			{Name: "quantity", Label: "Quantity", Kind: "Number",
				Options: cmpTypes.ModuleFieldOptions{"format": "0,0"}},
		},
	}
}

// TL;DR: an aggregated field's prefix/suffix travels with the number.
// Example: the agent summed a field prefixed "$ " and reported the total in
// euros. Nothing in the answer said what the unit was, and the module lookup
// that would have said so is a call it had no reason to make.
func TestMetricUnits(t *testing.T) {
	mod := unitsModule()

	t.Run("a prefixed field reports its prefix", func(t *testing.T) {
		u := metricUnits(mod, "SUM(line_value) AS total")
		require.Contains(t, u, "line_value")
		require.Equal(t, "$ ", u["line_value"].(map[string]any)["prefix"])
		require.Equal(t, "Total value", u["line_value"].(map[string]any)["label"])
	})

	t.Run("a suffixed field reports its suffix", func(t *testing.T) {
		u := metricUnits(mod, "AVG(change_pct)")
		require.Equal(t, " %", u["change_pct"].(map[string]any)["suffix"])
	})

	// A plain number has no unit to state, and listing it would only invite one
	// to be invented.
	t.Run("a field with neither is left out", func(t *testing.T) {
		require.Empty(t, metricUnits(mod, "SUM(quantity)"))
	})

	t.Run("several metrics at once", func(t *testing.T) {
		u := metricUnits(mod, "SUM(line_value) AS total, SUM(quantity) AS copies, AVG(change_pct)")
		require.Len(t, u, 2)
		require.Contains(t, u, "line_value")
		require.Contains(t, u, "change_pct")
	})

	t.Run("an unknown field is not invented", func(t *testing.T) {
		require.Empty(t, metricUnits(mod, "SUM(no_such_field)"))
	})

	t.Run("no metrics at all", func(t *testing.T) {
		require.Empty(t, metricUnits(mod, ""))
	})

	t.Run("whitespace inside the call is tolerated", func(t *testing.T) {
		require.Contains(t, metricUnits(mod, "SUM( line_value )"), "line_value")
	})
}
