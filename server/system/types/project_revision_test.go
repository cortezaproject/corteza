package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// A module the revision ADDED gets a mapping too — diffModules emits one so the
// plan lists it — but there is nothing behind it to read. Publishing that
// mapping asked the importer for a source table that does not exist in the old
// namespace and failed the entire publish, so the deployment plan's own
// suggestions could not be published back.
func TestModuleMappingHasSource(t *testing.T) {
	require.False(t, ModuleMapping{Module: "added"}.HasSource(),
		"a mapping with no fields at all is an added module")

	require.False(t, ModuleMapping{Module: "added", Fields: []ModuleFieldMapping{
		{TargetField: "title", Op: "default"},
		{TargetField: "body", Op: "default", Value: "x"},
	}}.HasSource(), "defaults invent values, they do not read a source")

	require.True(t, ModuleMapping{Module: "kept", Fields: []ModuleFieldMapping{
		{TargetField: "title", Op: "default"},
		{SourceField: "amount", TargetField: "amount", Op: "copy"},
	}}.HasSource(), "one field to copy is enough to need the migration")
}
