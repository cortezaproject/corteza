package importer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/permissions"
)

type (
	namespaceMock struct{ set types.NamespaceSet }
	moduleMock    struct{ set types.ModuleSet }
	chartMock     struct{ set types.ChartSet }
	pageMock      struct{ set types.PageSet }
)

func (mock *namespaceMock) FindByHandle(handle string) (o *types.Namespace, err error) {
	return
}

func (mock *moduleMock) FindByHandle(namespaceID uint64, handle string) (o *types.Module, err error) {
	return
}

func (mock *chartMock) FindByHandle(namespaceID uint64, handle string) (o *types.Chart, err error) {
	return
}

func (mock *pageMock) FindByHandle(namespaceID uint64, handle string) (o *types.Page, err error) {
	return
}

func TestChartImport_CastSet(t *testing.T) {
	var (
		namespace = &namespaceMock{}
		module    = &moduleMock{}
		chart     = &chartMock{}
		page      = &pageMock{}

		pi = permissions.NewImporter(nil)

		imp = NewImporter(namespace, module, chart, page, pi)

		ns = &types.Namespace{
			ID:      1000000,
			Name:    "Test",
			Slug:    "Test",
			Enabled: true,
		}

		impFixTester = func(t *testing.T, name string, fn func(*testing.T, *Chart)) {
			t.Run(name, func(t *testing.T) {
				var aux interface{}
				req := require.New(t)
				f, err := os.Open(fmt.Sprintf("testdata/%s.yaml", name))
				req.NoError(err)
				req.NoError(yaml.NewDecoder(f).Decode(&aux))
				req.NotNil(aux)
				ci := NewChartImporter(imp, ns)
				req.NoError(ci.CastSet(aux))
				fn(t, ci)
			})
		}
	)

	impFixTester(t, "chart_full_slice", func(t *testing.T, chart *Chart) {
		req := require.New(t)
		req.Len(chart.set, 2)
	})

	impFixTester(t, "chart_full", func(t *testing.T, chart *Chart) {
		req := require.New(t)
		req.Len(chart.set, 2)
		req.Equal(chart.set[0].Handle, "chart1")
		req.Equal(chart.set[0].Name, "chart 1")
	})
}
