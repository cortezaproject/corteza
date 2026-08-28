package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/assert"
)

func TestLabelFieldOf(t *testing.T) {
	mod := &cmpTypes.Module{Fields: cmpTypes.ModuleFieldSet{
		{Name: "code"},
		{Name: "title"},
	}}

	t.Run("named field wins", func(t *testing.T) {
		assert.Equal(t, "title", labelFieldOf(mod, "title").Name)
	})

	// The viewer falls back to the module's first field rather than showing a
	// raw ID, and this has to agree with it or the same record reads two ways.
	t.Run("falls back to the first field", func(t *testing.T) {
		assert.Equal(t, "code", labelFieldOf(mod, "").Name)
		assert.Equal(t, "code", labelFieldOf(mod, "no_such_field").Name)
	})

	t.Run("no fields, no label", func(t *testing.T) {
		assert.Nil(t, labelFieldOf(&cmpTypes.Module{}, "title"))
		assert.Nil(t, labelFieldOf(nil, "title"))
	})
}

// refLabels reaches the record and user services, so what is unit-testable is
// the part before that: it must not touch them at all when there is nothing to
// resolve. A nil module or an empty set reaching the store is a call made for
// every record read on a module with no references.
func TestRefLabelsSkipsEmptyWork(t *testing.T) {
	assert.Nil(t, refLabels(nil, nil, nil))
	assert.Nil(t, refLabels(nil, &cmpTypes.Module{}, nil))
	assert.Nil(t, refLabels(nil, &cmpTypes.Module{}, cmpTypes.RecordSet{}))
}

// nestedLabels reaches the module and record services, so what is unit-testable
// is that it decides not to before it gets there: a label field that is not
// itself a reference, or one pointing nowhere, must cost no query at all. The
// resolving half is exercised against a live module.
func TestNestedLabelsSkipsEmptyWork(t *testing.T) {
	recs := cmpTypes.RecordSet{
		{ID: 1, Values: cmpTypes.RecordValueSet{{Name: "candidate", Value: "7"}}},
	}

	t.Run("a plain label field has no second level", func(t *testing.T) {
		lf := &cmpTypes.ModuleField{Name: "full_name", Kind: "String"}
		assert.Nil(t, nestedLabels(nil, 1, lf, "full_name", recs))
	})

	t.Run("a reference naming no module resolves nothing", func(t *testing.T) {
		lf := &cmpTypes.ModuleField{Name: "candidate", Kind: "Record"}
		assert.Nil(t, nestedLabels(nil, 1, lf, "full_name", recs))
	})

	t.Run("no records means nothing to look up", func(t *testing.T) {
		lf := &cmpTypes.ModuleField{Name: "candidate", Kind: "Record", Options: cmpTypes.ModuleFieldOptions{"moduleID": "2"}}
		assert.Nil(t, nestedLabels(nil, 1, lf, "full_name", nil))
	})

	t.Run("records carrying no reference mean nothing to look up", func(t *testing.T) {
		lf := &cmpTypes.ModuleField{Name: "candidate", Kind: "Record", Options: cmpTypes.ModuleFieldOptions{"moduleID": "2"}}
		empty := cmpTypes.RecordSet{{ID: 1, Values: cmpTypes.RecordValueSet{{Name: "candidate", Value: "0"}}}}
		assert.Nil(t, nestedLabels(nil, 1, lf, "full_name", empty))
	})
}
