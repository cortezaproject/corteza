package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/require"
)

func selectModule() *cmpTypes.Module {
	return &cmpTypes.Module{
		ID:          1,
		NamespaceID: 1,
		Fields: cmpTypes.ModuleFieldSet{
			&cmpTypes.ModuleField{Name: "stage", Kind: "Select", Options: cmpTypes.ModuleFieldOptions{
				"options": []any{
					map[string]any{"value": "applied", "text": "Applied"},
					map[string]any{"value": "offer", "text": "Offer"},
					map[string]any{"value": "rejected", "text": "Rejected"},
				},
			}},
			&cmpTypes.ModuleField{Name: "rejection_reason", Kind: "Select", Options: cmpTypes.ModuleFieldOptions{
				"options": []any{
					map[string]any{"value": "skills", "text": "Skills mismatch"},
					map[string]any{"value": "compensation", "text": "Compensation"},
				},
			}},
			&cmpTypes.ModuleField{Name: "title", Kind: "String"},
			&cmpTypes.ModuleField{Name: "candidate", Kind: "Record", Options: cmpTypes.ModuleFieldOptions{"moduleID": "2"}},
			&cmpTypes.ModuleField{Name: "freeform", Kind: "Select"},
		},
	}
}

func Test_checkSelectValues(t *testing.T) {
	tcc := []struct {
		name string
		expr string
		errs string
	}{
		{"declared value passes", "stage = 'offer'", ""},
		{"the label is refused", "stage = 'Offer'", `"Offer" is not a value`},
		{"the label is named as the near match", "stage = 'Offer'", `did you mean "offer", the value shown as "Offer"`},
		{"the legal values are listed", "stage = 'Offer'", "Values: applied, offer, rejected"},
		{"a second Select field is covered too", "rejection_reason = 'Skills mismatch'", `did you mean "skills"`},
		{"wrong case alone is refused", "stage = 'OFFER'", `did you mean "offer"`},
		{"inequality is checked", "stage != 'Offer'", `"Offer" is not a value`},
		{"angle inequality is checked", "stage <> 'Offer'", `"Offer" is not a value`},
		{"a value nothing resembles gets no suggestion", "stage = 'zzz'", "Values: applied, offer, rejected"},
		{"the empty value is a real state", "stage = ''", ""},
		{"a prefilter interpolation is left alone", "stage = '${record.values.stage}'", ""},
		{"LIKE is a pattern, not a value", "stage LIKE '%Offer%'", ""},
		{"a non-Select field is not judged", "title = 'Offer'", ""},
		{"a Select with no declared options is not judged", "freeform = 'anything'", ""},
		{"a dotted reference path is left to resolveRefPaths", "candidate.stage = 'Offer'", ""},
		{"a term inside a literal is not a term", "title = 'stage = ''Offer'''", ""},
		{"the check reaches into a conjunction", "is_active = true AND stage = 'Offer'", `"Offer" is not a value`},
		{"a good term beside a bad one still fails", "stage = 'offer' OR stage = 'Rejected'", `"Rejected" is not a value`},
		{"every term good passes", "stage = 'offer' OR stage = 'rejected'", ""},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			err := checkSelectValues(selectModule(), "record list", tc.expr)
			if tc.errs == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.errs)
			require.Contains(t, err.Error(), "record list failed")
		})
	}
}

func Test_checkSelectValues_noModule(t *testing.T) {
	require.NoError(t, checkSelectValues(nil, "record list", "stage = 'Offer'"))
}

func Test_checkBlockFilters(t *testing.T) {
	mod := selectModule()
	load := func(id uint64) *cmpTypes.Module {
		if id == 1 {
			return mod
		}
		return nil
	}

	t.Run("a RecordList prefilter is checked", func(t *testing.T) {
		err := checkBlockFilters(load, "block \"Offers\" (RecordList)", map[string]any{
			"moduleID":  "1",
			"prefilter": "stage = 'Offer'",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), `block "Offers" (RecordList) failed`)
		require.Contains(t, err.Error(), `did you mean "offer"`)
	})

	t.Run("a Metric tile's own filter is checked", func(t *testing.T) {
		err := checkBlockFilters(load, "block \"Count\" (Metric)", map[string]any{
			"metrics": []any{
				map[string]any{"moduleID": "1", "filter": "stage = 'Offer'"},
			},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), `"Offer" is not a value`)
	})

	t.Run("a module the loader cannot resolve is not judged", func(t *testing.T) {
		require.NoError(t, checkBlockFilters(load, "block", map[string]any{
			"moduleID":  "99",
			"prefilter": "stage = 'Offer'",
		}))
	})

	t.Run("a filter with no module beside it is not judged", func(t *testing.T) {
		require.NoError(t, checkBlockFilters(load, "block", map[string]any{
			"prefilter": "stage = 'Offer'",
		}))
	})

	t.Run("a correct prefilter passes", func(t *testing.T) {
		require.NoError(t, checkBlockFilters(load, "block", map[string]any{
			"moduleID":  "1",
			"prefilter": "stage = 'offer'",
		}))
	})
}

func Test_selectOptionTexts_shapes(t *testing.T) {
	t.Run("a bare list of strings is value and text alike", func(t *testing.T) {
		f := &cmpTypes.ModuleField{Kind: "Select", Options: cmpTypes.ModuleFieldOptions{
			"options": []any{"a", "b"},
		}}
		require.Equal(t, map[string]string{"a": "a", "b": "b"}, selectOptionTexts(f))
	})

	t.Run("a typed string slice reads the same", func(t *testing.T) {
		f := &cmpTypes.ModuleField{Kind: "Select", Options: cmpTypes.ModuleFieldOptions{
			"options": []string{"a"},
		}}
		require.Equal(t, map[string]string{"a": "a"}, selectOptionTexts(f))
	})

	t.Run("no options at all is nothing to check against", func(t *testing.T) {
		require.Nil(t, selectOptionTexts(&cmpTypes.ModuleField{Kind: "Select"}))
	})
}
