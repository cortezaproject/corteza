package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

func TestModuleImport_CastSet(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		yaml     string
		modules  types.ModuleSet
	}{
		{name: "empty", yaml: ``},
		{name: "empty map", yaml: `{}`},
		{name: "empty slice", yaml: `[]`},
		{
			name: "single",
			yaml: `
example:
  name: Test Example
  fields:
    num: Numeric
    str:
      kind: String
      label: String
      multi: true
      private: false
  allow:
    everyone: edit
`,
			modules: []*types.Module{
				{
					Handle: "example",
					Name:   "Test Example",
					Fields: []*types.ModuleField{
						{Name: "num", Label: "num", Kind: "Numeric"},
						{Name: "str", Label: "String", Kind: "String", Place: 1, Multi: true, Private: false},
					},
				},
			},
		},
		{
			name: "no fields",
			yaml: `
example:
  name: Test Example
  fields:
`,
			modules: []*types.Module{
				{
					Handle: "example",
					Name:   "Test Example",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imp := &ModuleImport{
				permissions: &permissions.Importer{},
			}

			aux, err := importer.ParseYAML([]byte(tt.yaml))
			require.NoError(t, err)
			require.NoError(t, imp.CastSet(aux))
			require.Equal(t, tt.modules, imp.set)
		})
	}
}
