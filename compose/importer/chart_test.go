package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

func TestChartImport_CastSet(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		yaml     string
		charts   types.ChartSet
	}{
		{name: "empty", yaml: ``},
		{name: "empty map", yaml: `{}`},
		{name: "empty slice", yaml: `[]`},
		{
			name: "single",
			yaml: `
example:
  name: Test Example
  allow:
    everyone: edit
`,
			charts: []*types.Chart{
				{
					Handle: "example",
					Name:   "Test Example",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imp := &ChartImport{
				permissions: &permissions.Importer{},
			}

			aux, err := importer.ParseYAML([]byte(tt.yaml))
			require.NoError(t, err)
			require.NoError(t, imp.CastSet(aux))
			require.Equal(t, tt.charts, imp.set)
		})
	}
}
