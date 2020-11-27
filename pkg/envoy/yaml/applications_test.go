package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestApplication_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*application, error) {
			w := &application{}
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

	t.Run("application 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("application_1")
		req.NoError(err)
		req.NotNil(doc)
		req.Len(doc.applications, 2)
		req.NotNil(doc.applications[0])
		req.NotNil(doc.applications[1])
		req.Equal("one", doc.applications[0].res.Name)
		req.Equal("two", doc.applications[1].res.Name)
	})
}
