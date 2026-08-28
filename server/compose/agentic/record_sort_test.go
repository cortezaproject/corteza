package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/require"
)

func TestApplyRecordSort(t *testing.T) {
	mod := &cmpTypes.Module{Fields: cmpTypes.ModuleFieldSet{
		&cmpTypes.ModuleField{Name: "applied_on"},
		&cmpTypes.ModuleField{Name: "stage"},
	}}

	t.Run("an empty sort is not an error and orders nothing", func(t *testing.T) {
		f := &cmpTypes.RecordFilter{}
		require.NoError(t, applyRecordSort(f, mod, "   "))
		require.Empty(t, f.Sort)
	})

	t.Run("a module field sorts, keeping its direction", func(t *testing.T) {
		f := &cmpTypes.RecordFilter{}
		require.NoError(t, applyRecordSort(f, mod, "applied_on DESC"))
		require.Len(t, f.Sort, 1)
		require.Equal(t, "applied_on", f.Sort[0].Column)
		require.True(t, f.Sort[0].Descending)
	})

	t.Run("a record column sorts even though the module has no such field", func(t *testing.T) {
		f := &cmpTypes.RecordFilter{}
		require.NoError(t, applyRecordSort(f, mod, "createdAt DESC"))
		require.Equal(t, "createdAt", f.Sort[0].Column)
	})

	t.Run("recordID is taken as the store's ID", func(t *testing.T) {
		f := &cmpTypes.RecordFilter{}
		require.NoError(t, applyRecordSort(f, mod, "recordID ASC"))
		require.Equal(t, "ID", f.Sort[0].Column)
	})

	t.Run("an unknown column is refused, naming what exists", func(t *testing.T) {
		f := &cmpTypes.RecordFilter{}
		err := applyRecordSort(f, mod, "no_such_field DESC")
		require.Error(t, err)
		require.Contains(t, err.Error(), "no_such_field")
		require.Contains(t, err.Error(), "applied_on")
	})

	t.Run("every column is checked, not just the first", func(t *testing.T) {
		f := &cmpTypes.RecordFilter{}
		err := applyRecordSort(f, mod, "applied_on DESC, no_such_field ASC")
		require.Error(t, err)
		require.Contains(t, err.Error(), "no_such_field")
	})
}
