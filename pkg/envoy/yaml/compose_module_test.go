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
		req.NotNil(doc.compose.Modules[0].res.Fields.FindByName("AccountSource"))
		req.Equal(21, len(doc.compose.Modules[1].res.Fields))
		req.NotNil(doc.compose.Modules[1].res.Fields.FindByName("NumberOfConvertedLeads"))
		req.Equal(23, len(doc.compose.Modules[2].res.Fields))
		req.NotNil(doc.compose.Modules[0].rbac)
		req.NotEmpty(doc.compose.Modules[0].rbac)
	})

	t.Run("doc 2", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_module_2")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.Modules, 1)
		req.Len(doc.compose.Modules[0].res.Fields, 2)

		full := doc.compose.Modules[0].res.Fields.FindByName("full")
		vaKV := doc.compose.Modules[0].res.Fields.FindByName("validatorKV")

		req.Equal("a > b", full.Expressions.ValueExpr)
		req.Len(full.Expressions.Sanitizers, 1)
		req.Contains(full.Expressions.Sanitizers, "trim(value)")
		req.Len(full.Expressions.Validators, 1)
		req.Equal(`a == ""`, full.Expressions.Validators[0].Test)
		req.Equal("Value should not be empty", full.Expressions.Validators[0].Error)
		req.True(full.Expressions.DisableDefaultValidators)

		req.Len(vaKV.Expressions.Validators, 1)
		req.Equal(`value == ""`, vaKV.Expressions.Validators[0].Test)
		req.Equal("Value should be filled", vaKV.Expressions.Validators[0].Error)

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
		req.Equal(a.res.Operation, "read")
		req.Equal(a.res.Access, rbac.Allow)
		req.Equal(b.res.Operation, "delete")
		req.Equal(b.res.Access, rbac.Deny)
	})
}

func TestComposeModule_Shaping(t *testing.T) {
	t.Run("1 key; list mapping", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_module_s_1")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.Modules, 1)
		m1 := doc.compose.Modules[0]

		tp := m1.recTpl
		req.NotNil(tp)
		req.Equal("mod1.csv", tp.Source)
		req.Equal([]string{"id"}, tp.Key)

		req.Len(tp.Mapping, 2)
		req.Equal(uint(0), tp.Mapping[0].Index)
		req.Equal("f1", tp.Mapping[0].Field)
		req.Equal(uint(3), tp.Mapping[1].Index)
		req.Equal("f2", tp.Mapping[1].Field)
	})

	t.Run("2 keys", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_module_s_2")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.Modules, 1)
		m1 := doc.compose.Modules[0]

		tp := m1.recTpl
		req.NotNil(tp)
		req.Equal([]string{"f1", "f2"}, tp.Key)
	})

	t.Run("complex mapping", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_module_s_3")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.Modules, 1)
		m1 := doc.compose.Modules[0]

		tp := m1.recTpl
		req.Len(tp.Mapping, 2)
		req.Equal(uint(0), tp.Mapping[0].Index)
		req.Equal("c1", tp.Mapping[0].Cell)
		req.Equal("f1", tp.Mapping[0].Field)
		req.Equal(uint(0), tp.Mapping[1].Index)
		req.Equal("c2", tp.Mapping[1].Cell)
		req.Equal("f2", tp.Mapping[1].Field)
	})
}
