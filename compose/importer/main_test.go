package importer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
)

var (
	ns = &types.Namespace{
		ID:      1000000,
		Name:    "Test",
		Slug:    "testing",
		Enabled: true,
	}

	// Add namespace to the stack, make sure importer can find it
	pi *permissions.Importer
	st *settings.Importer

	imp *Importer
)

func TestMain(m *testing.M) {
	resetMocks()
	os.Exit(m.Run())
}

func resetMocks() {
	// whitelist = nil, anything can be added
	pi = permissions.NewImporter(service.AccessControl(nil).Whitelist())
	st = settings.NewImporter()

	imp = NewImporter(nil, nil, nil, nil, pi, st)

	// namespaces does not get initialized in the standard flow
	// if namespace finder is not present
	imp.namespaces = NewNamespaceImporter(imp)

}

func impFixTester(t *testing.T, name string, tester interface{}) {
	t.Run(name, func(t *testing.T) {
		// We're not calling reset mocks BEFORE calling tester()
		// because we want to have an option to set it up as we want
		defer resetMocks()

		var aux interface{}
		req := require.New(t)
		f, err := os.Open(fmt.Sprintf("testdata/%s.yaml", name))
		req.NoError(err)
		req.NoError(yaml.NewDecoder(f).Decode(&aux))
		req.NotNil(aux)

		if reqError, ok := tester.(error); ok {
			req.EqualError(imp.GetNamespaceImporter().Cast(ns.Slug, aux), reqError.Error())
			return
		} else if _, ok := tester.(func(*testing.T, *Importer)); ok {
			// Cast with compose root importer
			req.NoError(imp.Cast(aux))
		} else {
			// Cast with namespace importer
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
		case func(*testing.T, *Record):
			tester(t, imp.GetRecordImporter(ns.Slug))
		case func(*testing.T, *Importer):
			tester(t, imp)
		default:
			panic("unsupported tester function signature")
		}
	})
}
