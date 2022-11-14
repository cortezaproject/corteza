package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRole_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*role, error) {
			w := &role{}
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

	t.Run("role 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("role_1")
		req.NoError(err)
		req.NotNil(doc)
		req.Len(doc.roles, 2)
		req.NotNil(doc.roles[1])
		req.Equal("Role2", doc.roles[1].res.Name)
		req.Equal("r2", doc.roles[1].res.Handle)
	})
}
