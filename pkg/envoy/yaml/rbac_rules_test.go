package yaml

import (
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRBACRules_UnmarshalYAML(t *testing.T) {
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

	t.Run("compose rbac 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_rbac_1")
		req.NoError(err)
		req.NotNil(doc)

		req.NotNil(doc.compose)
		req.Len(doc.compose.Modules, 1)
		req.NotNil(doc.compose.Modules[0].rbac)
		req.NotEmpty(doc.compose.Modules[0].rbac)
	})

}

func TestRBACRules_MarshalEnvoy(t *testing.T) {
	var (
		req = require.New(t)

		wrap    = composeNamespace{res: &types.Namespace{ID: 42}}
		nn, err = wrap.MarshalEnvoy()
	)

	req.NoError(err)
	req.Len(nn, 1)
	// req.IsType(&node.ComposeNamespace{}, nn[0])
	// req.Equal(uint64(42), nn[0].(*node.ComposeNamespace).Res.ID)
}
