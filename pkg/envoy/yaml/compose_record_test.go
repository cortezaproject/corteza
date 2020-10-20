package yaml

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestComposeRecord_UnmarshalYAML(t *testing.T) {
	var (
		req = require.New(t)

		w   *ComposeRecord
		err error

		parse = func(w *ComposeRecord, src string) (*ComposeRecord, error) {
			return w, yaml.Unmarshal([]byte(src), w)
		}
	)

	w, err = parse(&ComposeRecord{}, ``)
	req.NoError(err)
	req.NotNil(w)
	req.Nil(w.res)

	w, err = parse(&ComposeRecord{}, `{ values: { foo: bar }, createdBy: foo, updatedAt: 2020-10-10T10:10:00Z, deletedBy: user }`)
	req.NoError(err)
	req.NotNil(w)
	req.NotNil(w.res)
	req.NotEmpty(w.res.Values)
	req.NotEmpty(w.res.UpdatedAt)
	req.Equal("bar", w.res.Values.Get("foo", 0).Value)
}
