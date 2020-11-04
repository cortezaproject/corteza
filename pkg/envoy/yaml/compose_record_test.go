package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestComposeRecord_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*composeRecord, error) {
			w := &composeRecord{}
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

	t.Run("empty", func(t *testing.T) {
		req := require.New(t)

		w, err := parseString(`{ values: { foo: bar }, createdBy: foo, updatedAt: 2020-10-10T10:10:00Z, deletedBy: user }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.NotEmpty(w.res.Values)
		req.NotEmpty(w.res.UpdatedAt)
		req.Equal("bar", w.res.Values.Get("foo", 0).Value)
	})

	t.Run("compose record file 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("compose_record_1")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.compose)
		req.Len(doc.compose.Records, 3)

		req.NotNil(doc.compose.Records[0].res)
		req.Equal("Department", doc.compose.Records[0].refModule)
		rec := doc.compose.Records[0].res
		req.Equal("Service", rec.Values.Get("Name", 0).Value)
		req.Equal("50", rec.Values.Get("HourCost", 0).Value)

		req.NotNil(doc.compose.Records[1].res)
		req.Equal("EmailTemplate", doc.compose.Records[1].refModule)

		req.NotNil(doc.compose.Records[2].res)
		req.Equal("Settings", doc.compose.Records[2].refModule)

		//req.NotNil(doc.compose.records[0].rbac)
		//req.NotEmpty(doc.compose.records[0].rbac.rules)
	})
}
