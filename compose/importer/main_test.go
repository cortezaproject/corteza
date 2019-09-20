package importer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/permissions"
)

type (
	namespaceMock struct{ set types.NamespaceSet }
	moduleMock    struct{ set types.ModuleSet }
	chartMock     struct{ set types.ChartSet }
	pageMock      struct{ set types.PageSet }
)

var (
	ns = &types.Namespace{
		ID:      1000000,
		Name:    "Test",
		Slug:    "test",
		Enabled: true,
	}

	// Add namespace to the stack, make sure importer can find it
	namespaces = &namespaceMock{set: types.NamespaceSet{ns}}
	modules    = &moduleMock{}
	charts     = &chartMock{}
	pages      = &pageMock{}

	// whitelist = nil, anything can be added
	pi = permissions.NewImporter(service.AccessControl(nil).Whitelist())

	imp = NewImporter(namespaces, modules, charts, pages, pi)
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func (mock *namespaceMock) FindByHandle(slug string) (o *types.Namespace, err error) {
	oo, err := mock.set.Filter(func(o *types.Namespace) (b bool, e error) {
		return o.Slug == slug, nil
	})

	if len(oo) > 0 {
		return oo[0], nil
	} else {
		return nil, repository.ErrNamespaceNotFound
	}
}

func (mock *moduleMock) FindByHandle(namespaceID uint64, handle string) (o *types.Module, err error) {
	oo, err := mock.set.Filter(func(o *types.Module) (b bool, e error) {
		return o.Handle == handle && o.NamespaceID == namespaceID, nil
	})

	if len(oo) > 0 {
		return oo[0], nil
	} else {
		return nil, repository.ErrModuleNotFound
	}
}

func (mock *chartMock) FindByHandle(namespaceID uint64, handle string) (o *types.Chart, err error) {
	oo, err := mock.set.Filter(func(o *types.Chart) (b bool, e error) {
		return o.Handle == handle && o.NamespaceID == namespaceID, nil
	})

	if len(oo) > 0 {
		return oo[0], nil
	} else {
		return nil, repository.ErrChartNotFound
	}
}

func (mock *pageMock) FindByHandle(namespaceID uint64, handle string) (o *types.Page, err error) {
	oo, err := mock.set.Filter(func(o *types.Page) (b bool, e error) {
		return o.Handle == handle && o.NamespaceID == namespaceID, nil
	})

	if len(oo) > 0 {
		return oo[0], nil
	} else {
		return nil, repository.ErrPageNotFound
	}
}

func impFixTester(t *testing.T, name string, tester interface{}) {
	t.Run(name, func(t *testing.T) {
		var aux interface{}
		req := require.New(t)
		f, err := os.Open(fmt.Sprintf("testdata/%s.yaml", name))
		req.NoError(err)
		req.NoError(yaml.NewDecoder(f).Decode(&aux))
		req.NotNil(aux)

		if reqError, ok := tester.(error); ok {
			req.EqualError(imp.GetNamespaceImporter().Cast(ns.Slug, aux), reqError.Error())
			return
		} else {
			req.NoError(imp.GetNamespaceImporter().Cast(ns.Slug, aux))
		}

		switch tester := tester.(type) {
		case func(*testing.T, *Namespace):
			tester(t, imp.GetNamespaceImporter())
		case func(*testing.T, *Module):
			tester(t, imp.GetModuleImporter(ns.Slug))
		case func(*testing.T, *Chart):
			tester(t, imp.GetChartImporter(ns.Slug))
		case func(*testing.T, *Page):
			tester(t, imp.GetPageImporter(ns.Slug))
		case func(*testing.T, *Importer):
			tester(t, imp)
		default:
			panic("unsupported tester function signature")
		}
	})
}
