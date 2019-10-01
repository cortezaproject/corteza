package importer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/system/service"
)

var (
	pi *permissions.Importer

	imp *Importer
)

func TestMain(m *testing.M) {
	resetMocks()
	os.Exit(m.Run())
}

func resetMocks() {
	// whitelist = nil, anything can be added
	pi = permissions.NewImporter(service.AccessControl(nil).Whitelist())

	imp = NewImporter(
		pi,
		NewChannelImport(pi, nil),
	)
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
			req.EqualError(imp.Cast(aux), reqError.Error())
			return
		} else {
			req.NoError(imp.Cast(aux))
		}

		switch tester := tester.(type) {
		case func(*testing.T, *Channel):
			tester(t, imp.channels)
		case func(*testing.T, *Importer):
			tester(t, imp)
		default:
			panic("unsupported tester function signature")
		}
	})
}
