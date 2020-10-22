package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestComposeNamespace_UnmarshalYAML(t *testing.T) {
	var (
		req = require.New(t)

		parseString = func(src string) (*ComposeNamespace, error) {
			w := &ComposeNamespace{}
			return w, yaml.Unmarshal([]byte(src), w)
		}

		parseDocument = func(i int) (*Document, error) {
			doc := &Document{}
			f, err := os.Open(fmt.Sprintf("testdata/compose_namespace_%d.yaml", i))
			if err != nil {
				return nil, err
			}

			return doc, yaml.NewDecoder(f).Decode(doc)
		}
	)

	t.Run("empty", func(t *testing.T) {
		w, err := parseString(``)
		req.NoError(err)
		req.NotNil(w)
		req.Nil(w.res)
	})

	t.Run("simple name", func(t *testing.T) {
		w, err := parseString(`{ name: Test }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.NotEmpty(w.res.Name)
		req.True(w.res.Enabled)
	})

	t.Run("disabled", func(t *testing.T) {
		w, err := parseString(`{ enabled: false }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.False(w.res.Enabled)
	})

	t.Run("compose namespace 1", func(t *testing.T) {
		doc, err := parseDocument(1)
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.namespaces, 1)
		req.NotNil(doc.compose.namespaces[0])
		req.Equal("CRM", doc.compose.namespaces[0].res.Name)
		req.Equal("crm", doc.compose.namespaces[0].res.Slug)
	})

}

func TestComposeNamespace_MarshalEnvoy(t *testing.T) {
	var (
		req = require.New(t)

		wrap    = ComposeNamespace{res: &types.Namespace{ID: 42}}
		nn, err = wrap.MarshalEnvoy()
	)

	req.NoError(err)
	req.NotEmpty(nn)
	req.Equal(uint64(42), nn[0].(*envoy.ComposeNamespaceNode).Ns.ID)
}
