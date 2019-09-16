package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

func TestPageImport_CastSet(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		yaml     string
		pages    types.PageSet
		parents  map[string]string
	}{
		{name: "empty", yaml: ``},
		{name: "empty map", yaml: `{}`},
		{name: "empty slice", yaml: `[]`},
		{
			name: "single",
			yaml: `
example:
  title: Test Example
`,
			pages: []*types.Page{
				{
					Handle: "example",
					Title:  "Test Example",
				},
			},
		},
		{
			name: "nested",
			yaml: `
example:
  title: Test Example
  pages:
    sub1: Sub page 1
    sub2:
      pages:
        subsub1: "SUB-SUB #1"
`,
			parents: map[string]string{
				"sub1":    "example",
				"sub2":    "example",
				"subsub1": "sub2",
			},
			pages: []*types.Page{
				{Handle: "example", Title: "Test Example"},
				{Handle: "sub1", Title: "Sub page 1"},
				{Handle: "sub2", Title: "sub2"},
				{Handle: "subsub1", Title: "SUB-SUB #1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.parents == nil {
				tt.parents = map[string]string{}
			}

			if tt.pages == nil {
				tt.pages = types.PageSet{}
			}

			imp := NewPageImporter(nil, nil)

			aux, err := importer.ParseYAML([]byte(tt.yaml))
			require.NoError(t, err)
			require.NoError(t, imp.CastSet(aux))
			require.Equal(t, tt.pages, imp.set)

			require.Equal(t, tt.parents, imp.parents)
		})
	}
}
