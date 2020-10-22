package yaml

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestComposeChart_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*composeChart, error) {
			w := &composeChart{}
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

		doc, err := parseDocument("compose_chart_1")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.charts, 2)
		req.Equal(1, len(doc.compose.charts[0].res.Config.Reports))
		req.NotNil(doc.compose.charts[1].rbacRules)
		req.NotEmpty(doc.compose.charts[1].rbacRules.rules)
	})
}
