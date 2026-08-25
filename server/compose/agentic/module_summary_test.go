package agentic

import (
	"encoding/json"
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/require"
)

// TL;DR: the summary carries what reading data needs and drops what storing it
// needs.
// Example: the stored shape repeats moduleID and namespaceID on every field and
// carries each field's ID, timestamps and DAL config — about 500 bytes a field,
// so learning a namespace cost tens of kilobytes of an agent's context.
func TestModuleSummary(t *testing.T) {
	m := &cmpTypes.Module{
		ID: 12, NamespaceID: 100, Handle: "collection_item", Name: "Collection Item",
		Fields: cmpTypes.ModuleFieldSet{
			{Name: "quantity", Kind: "Number", Label: "Quantity", Required: true},
			{Name: "notes", Kind: "String", Label: "notes"},
			{Name: "rarity", Kind: "Select", Label: "Rarity",
				Options: cmpTypes.ModuleFieldOptions{"options": []any{"Rare"}}},
			{Name: "line_value", Kind: "Number", Label: "Total",
				Expressions: cmpTypes.ModuleFieldExpr{ValueExpr: "quantity * unit_value"}},
		},
	}

	out := moduleSummary(m)
	fields := out["fields"].([]map[string]any)
	require.Len(t, fields, 4)

	byName := map[string]map[string]any{}
	for _, f := range fields {
		byName[f["name"].(string)] = f
	}

	t.Run("keeps what a reader needs", func(t *testing.T) {
		require.Equal(t, "Number", byName["quantity"]["kind"])
		require.Equal(t, true, byName["quantity"]["required"])
		require.Contains(t, byName["rarity"], "options")
	})

	// A derived field cannot be written to; a caller that does not know writes a
	// value the next save overwrites.
	t.Run("says which fields are derived", func(t *testing.T) {
		require.Equal(t, "quantity * unit_value", byName["line_value"]["valueExpression"])
		require.NotContains(t, byName["quantity"], "valueExpression")
	})

	t.Run("drops the storage detail", func(t *testing.T) {
		for _, f := range fields {
			for _, gone := range []string{"fieldID", "moduleID", "namespaceID", "createdAt", "config", "expressions"} {
				require.NotContains(t, f, gone)
			}
		}
	})

	// A label that only repeats the name is noise on every field.
	t.Run("omits a label that says nothing", func(t *testing.T) {
		require.NotContains(t, byName["notes"], "label")
		require.Equal(t, "Quantity", byName["quantity"]["label"])
	})

	t.Run("keeps the module identifiers", func(t *testing.T) {
		require.Equal(t, "12", out["moduleID"])
		require.Equal(t, "100", out["namespaceID"])
		require.Equal(t, "collection_item", out["handle"])
	})

	t.Run("is markedly smaller than the stored shape", func(t *testing.T) {
		lean, err := json.Marshal(out)
		require.NoError(t, err)
		full, err := json.Marshal(m)
		require.NoError(t, err)
		require.Less(t, len(lean), len(full)/2, "the point of the projection is the size")
	})
}
