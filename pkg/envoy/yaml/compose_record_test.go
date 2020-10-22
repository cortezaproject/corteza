package yaml

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestComposeRecord_UnmarshalYAML(t *testing.T) {
	var (
		req = require.New(t)

		parseString = func(src string) (*ComposeRecord, error) {
			w := &ComposeRecord{}
			return w, yaml.Unmarshal([]byte(src), w)
		}

		parseDocument = func(i int) (*Document, error) {
			doc := &Document{}
			f, err := os.Open(fmt.Sprintf("testdata/compose_record_%d.yaml", i))
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

	t.Run("empty", func(t *testing.T) {
		w, err := parseString(`{ values: { foo: bar }, createdBy: foo, updatedAt: 2020-10-10T10:10:00Z, deletedBy: user }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.NotEmpty(w.res.Values)
		req.NotEmpty(w.res.UpdatedAt)
		req.Equal("bar", w.res.Values.Get("foo", 0).Value)
	})

	t.Run("compose record file 1", func(t *testing.T) {
		doc, err := parseDocument(1)
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.records, 3)

		req.NotNil(doc.compose.records[0].res)
		req.Equal("Department", doc.compose.records[0].refModule)
		rec := doc.compose.records[0].res
		req.Equal("Service", rec.Values.Get("Name", 0).Value)
		req.Equal("50", rec.Values.Get("HourCost", 0).Value)

		req.NotNil(doc.compose.records[1].res)
		req.Equal("EmailTemplate", doc.compose.records[1].refModule)

		req.NotNil(doc.compose.records[2].res)
		req.Equal("Settings", doc.compose.records[2].refModule)
	})
}
