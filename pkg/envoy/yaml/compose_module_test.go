package yaml

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestComposeModule_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*composeModule, error) {
			w := &composeModule{}
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

	t.Run("field with default value", func(t *testing.T) {
		req := require.New(t)

		w, err := parseString(`{ fields: { one: { default: foo }, two: { default: [ foo, bar ] } } }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.Len(w.res.Fields, 2)
		req.Len(w.res.Fields[0].DefaultValue, 1)
		req.Equal("foo", w.res.Fields[0].DefaultValue[0].Value)
		req.Len(w.res.Fields[1].DefaultValue, 2)
		req.Equal("foo", w.res.Fields[1].DefaultValue[0].Value)
		req.Equal("bar", w.res.Fields[1].DefaultValue[1].Value)
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
		req.NotNil(doc.compose.modules[0].rbac)
		req.NotEmpty(doc.compose.modules[0].rbac.rules)
	})
}
