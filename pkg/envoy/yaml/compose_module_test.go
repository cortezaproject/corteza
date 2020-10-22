package yaml

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestComposeModule_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*ComposeModule, error) {
			w := &ComposeModule{}
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
	})

	t.Run("doc 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_module_1")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.modules, 3)
		req.Equal(30, len(doc.compose.modules[0].res.Fields))
		req.Equal(21, len(doc.compose.modules[1].res.Fields))
		req.Equal(23, len(doc.compose.modules[2].res.Fields))
		req.NotNil(doc.compose.modules[0].rbacRules)
		req.NotEmpty(doc.compose.modules[0].rbacRules.rules)
	})
}
