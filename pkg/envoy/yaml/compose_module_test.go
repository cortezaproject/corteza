package yaml

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestComposeModule_UnmarshalYAML(t *testing.T) {
	var (
		req = require.New(t)

		parseString = func(src string) (*ComposeNamespace, error) {
			w := &ComposeNamespace{}
			return w, yaml.Unmarshal([]byte(src), w)
		}

		parseDocument = func(i int) (*Document, error) {
			doc := &Document{}
			f, err := os.Open(fmt.Sprintf("testdata/compose_module_%d.yaml", i))
			if err != nil {
				return nil, err
			}

			return doc, yaml.NewDecoder(f).Decode(doc)
		}
	)

	t.Run("empty", func(t *testing.T) {
		w, err := parseString(``)
		req.NoError(err)
		req.NotNil(w)
		req.Nil(w.res)
	})

	t.Run("simple name", func(t *testing.T) {
		w, err := parseString(`{ name: Test }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.NotEmpty(w.res.Name)
	})

	t.Run("compose module file 1", func(t *testing.T) {
		doc, err := parseDocument(1)
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.modules, 3)
		req.Equal(30, len(doc.compose.modules[0].res.Fields))
		req.Equal(21, len(doc.compose.modules[1].res.Fields))
		req.Equal(23, len(doc.compose.modules[2].res.Fields))
	})
}
