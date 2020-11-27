package yaml

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
		req.Len(doc.compose.Modules, 3)
		req.Equal(30, len(doc.compose.Modules[0].res.Fields))
		req.Equal(21, len(doc.compose.Modules[1].res.Fields))
		req.Equal(23, len(doc.compose.Modules[2].res.Fields))
		req.NotNil(doc.compose.Modules[0].rbac)
		req.NotEmpty(doc.compose.Modules[0].rbac)
	})

	t.Run("doc rbac", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_module_rbac")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)

		mod := doc.compose.Modules[0]
		req.Len(mod.rbac, 2)
		a := mod.rbac[0]
		b := mod.rbac[1]
		req.Equal(a.res.Operation, rbac.Operation("read"))
		req.Equal(a.res.Access, rbac.Allow)
		req.Equal(b.res.Operation, rbac.Operation("delete"))
		req.Equal(b.res.Access, rbac.Deny)
	})
}
