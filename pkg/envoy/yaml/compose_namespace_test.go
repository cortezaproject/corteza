package yaml

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestComposeNamespace_UnmarshalYAML(t *testing.T) {
	var (
		req = require.New(t)

		string = func(src string) (*ComposeNamespace, error) {
			w := &ComposeNamespace{}
			return w, yaml.Unmarshal([]byte(src), w)
		}

		file = func(i int) (*Document, error) {
			doc := &Document{}
			f, err := os.Open(fmt.Sprintf("testdata/compose_namespace_%d.yaml", i))
			if err != nil {
				return nil, err
			}

			return doc, yaml.NewDecoder(f).Decode(doc)
		}
	)

	t.Run("empty", func(t *testing.T) {
		w, err := string(``)
		req.NoError(err)
		req.NotNil(w)
		req.Nil(w.res)
	})

	t.Run("simple name", func(t *testing.T) {
		w, err := string(`{ name: Test }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.NotEmpty(w.res.Name)
		req.True(w.res.Enabled)
	})

	t.Run("disabled", func(t *testing.T) {
		w, err := string(`{ enabled: false }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.False(w.res.Enabled)
	})

	t.Run("compose namespace file 1", func(t *testing.T) {
		doc, err := file(1)
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.namespaces, 1)
		req.NotNil(doc.compose.namespaces[0])
		req.Equal("CRM", doc.compose.namespaces[0].res.Name)
		req.Equal("crm", doc.compose.namespaces[0].res.Slug)
	})

	t.Run("compose namespace file 2", func(t *testing.T) {
		doc, err := file(2)
		req.NoError(err)
		req.NotNil(doc)
		spew.Dump(doc.compose)
		req.NotNil(doc.compose)
		req.Len(doc.compose.namespaces, 1)
		req.NotNil(doc.compose.namespaces[0])

		req.Equal("CRM", doc.compose.namespaces[0].res.Name)
		req.Equal("crm", doc.compose.namespaces[0].res.Slug)
	})

}
