package yaml

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
		req.Len(doc.compose.Charts, 2)
		req.Equal(1, len(doc.compose.Charts[0].res.Config.Reports))
		req.NotNil(doc.compose.Charts[1].rbac)
		req.NotEmpty(doc.compose.Charts[1].rbac)
	})

	t.Run("doc rbac", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_chart_rbac")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)

		ch := doc.compose.Charts[0]
		req.Len(ch.rbac, 2)
		a := ch.rbac[0]
		b := ch.rbac[1]
		req.Equal(a.res.Operation, "read")
		req.Equal(a.res.Access, rbac.Allow)
		req.Equal(b.res.Operation, "delete")
		req.Equal(b.res.Access, rbac.Deny)
	})
}
