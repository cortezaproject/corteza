package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/assert"
)

func draftSet(pairs ...string) cmpTypes.RecordValueSet {
	out := cmpTypes.RecordValueSet{}
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, &cmpTypes.RecordValue{Name: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMergeDraftReplacesOnlyTheNamedFields(t *testing.T) {
	current := draftSet("title", "Old", "tags", "a", "tags", "b", "stage", "lead")
	proposed := draftSet("title", "New")

	merged := mergeDraft(current, proposed, []string{"title", "tags"})

	assert.Equal(t, draftSet("stage", "lead", "title", "New"), merged,
		"a named field with no proposed value is cleared; an unnamed one keeps the record's value")
}

func TestDraftValuesGroupsMultiValueFields(t *testing.T) {
	mod := &cmpTypes.Module{Fields: cmpTypes.ModuleFieldSet{
		{Name: "title", Kind: "String"},
		{Name: "tags", Kind: "String", Multi: true},
	}}

	got := draftValues(mod, draftSet("title", "Kickoff", "tags", "a", "tags", "b", "unknown", "x"))

	assert.Equal(t, map[string]any{"title": "Kickoff", "tags": []string{"a", "b"}, "unknown": "x"}, got)
}

func TestChoicesAreOrderedByLabelAndFallBackToTheID(t *testing.T) {
	c := choicesOf([]uint64{3, 1, 2}, map[string]string{"1": "Zoe", "3": "Adam"}, true)

	assert.Equal(t, []fieldChoice{{ID: "2", Label: "2"}, {ID: "3", Label: "Adam"}, {ID: "1", Label: "Zoe"}}, c.Items)
	assert.True(t, c.More)
}

func TestRecordFormViewCarriesRequiredAndOptions(t *testing.T) {
	mod := &cmpTypes.Module{Name: "Deals", Fields: cmpTypes.ModuleFieldSet{
		{Name: "amount", Kind: "Number", Required: true, Options: cmpTypes.ModuleFieldOptions{"prefix": "€ "}},
	}}

	v := recordFormViewOf(mod, nil)

	assert.Equal(t, []formField{{
		Name: "amount", Kind: "Number", Required: true, Options: cmpTypes.ModuleFieldOptions{"prefix": "€ "},
	}}, v.Fields)
}
