package agentic

import (
	"errors"
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/assert"
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

// TL;DR: only an unknown-attribute failure gets the field list appended, and a
// report that cannot reach the module still reports the original failure.
// Example: the DAL says "unknown attribute card.rarity" and nothing else. A
// model given only that guesses again or falls back to reading every record and
// adding them up by hand — which is what this tool exists to stop.
func TestReportError(t *testing.T) {
	ctx := t.Context()

	// DefaultModule is not wired here, so enrichment cannot happen; what must
	// survive is the original cause either way.
	t.Run("an unrelated failure is passed through", func(t *testing.T) {
		err := reportError(ctx, 1, 2, errors.New("connection refused"))
		require.ErrorContains(t, err, "connection refused")
		require.NotContains(t, err.Error(), "Available:")
	})

	t.Run("an unknown attribute still names the attribute", func(t *testing.T) {
		err := reportError(ctx, 1, 2, errors.New("unknown attribute card.rarity"))
		require.ErrorContains(t, err, "card.rarity")
	})

	t.Run("the cause stays wrapped", func(t *testing.T) {
		cause := errors.New("unknown attribute card.rarity")
		require.ErrorIs(t, reportError(ctx, 1, 2, cause), cause)
	})
}

// A total over an expression carries the units of the fields in it. Reading the
// whole parenthesised body as one field name found none, and a wishlist priced
// in dollars was reported back in euros.
func TestMetricUnitsAcrossAnExpression(t *testing.T) {
	mod := &cmpTypes.Module{Fields: cmpTypes.ModuleFieldSet{
		{Name: "quantity_wanted", Kind: "Number", Label: "Wanted"},
		{Name: "current_price", Kind: "Number", Label: "Current price",
			Options: cmpTypes.ModuleFieldOptions{"prefix": "$ "}},
		{Name: "change_pct", Kind: "Number", Label: "Change",
			Options: cmpTypes.ModuleFieldOptions{"suffix": " %"}},
	}}

	units := metricUnits(mod, "SUM(quantity_wanted * current_price) AS total, AVG(change_pct) AS avg")

	assert.Len(t, units, 2, "a field with neither prefix nor suffix has no unit to state")
	assert.Equal(t, "$ ", units["current_price"].(map[string]any)["prefix"])
	assert.Equal(t, " %", units["change_pct"].(map[string]any)["suffix"])
	assert.NotContains(t, units, "quantity_wanted")
}
