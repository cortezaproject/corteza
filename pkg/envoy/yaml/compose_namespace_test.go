package yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/node"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestComposeNamespace_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*composeNamespace, error) {
			w := &composeNamespace{}
			return w, yaml.Unmarshal([]byte(src), w)
		}
	)

	t.Run("empty", func(t *testing.T) {
		req := require.New(t)

		w, err := parseString(``)
		req.NoError(err)
		req.NotNil(w)
		req.Nil(w.res)
	})

	t.Run("simple name", func(t *testing.T) {
		req := require.New(t)

		w, err := parseString(`{ name: Test }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.NotEmpty(w.res.Name)
		req.True(w.res.Enabled)
	})

	t.Run("disabled", func(t *testing.T) {
		req := require.New(t)

		w, err := parseString(`{ enabled: false }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.False(w.res.Enabled)
	})

	t.Run("compose namespace 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_namespace_1")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.namespaces, 1)
		req.NotNil(doc.compose.namespaces[0])
		req.Equal("CRM", doc.compose.namespaces[0].res.Name)
		req.Equal("crm", doc.compose.namespaces[0].res.Slug)
		req.NotNil(doc.compose.namespaces[0].rbac)
		req.NotEmpty(doc.compose.namespaces[0].rbac.rules)
	})

}

func TestComposeNamespace_MarshalEnvoy(t *testing.T) {
	var (
		req = require.New(t)

		wrap    = composeNamespace{res: &types.Namespace{ID: 42}}
		nn, err = wrap.MarshalEnvoy()
	)

	req.NoError(err)
	req.Len(nn, 1)
	req.IsType(&node.ComposeNamespace{}, nn[0])
	req.Equal(uint64(42), nn[0].(*node.ComposeNamespace).Res.ID)
}
