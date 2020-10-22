package yaml

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestComposePage_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*composePage, error) {
			w := &composePage{}
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

		w, err := parseString(`{ title: Test }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.NotEmpty(w.res.Title)
	})

	t.Run("doc 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_page_1")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.pages, 1)
		req.Len(doc.compose.pages[0].pages, 1)
		req.Equal(3, len(doc.compose.pages[0].res.Blocks))
		req.NotNil(doc.compose.pages[0].rbacRules)
		req.NotEmpty(doc.compose.pages[0].rbacRules.rules)
	})
}
