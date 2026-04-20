package json

import (
	"io"
	"strings"
	"testing"

	"github.com/crusttech/human/server/pkg/envoyx/datasource"
	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	req := require.New(t)

	t.Run("init & meta", func(t *testing.T) {
		dc, err := Decoder(testReader(), "test.csv")
		req.NoError(err)

		req.Equal(uint64(3), dc.Count())
	})

	t.Run("fields", func(t *testing.T) {
		dc, err := Decoder(testReader(), "test.csv")
		req.NoError(err)

		hh := dc.Fields()
		req.Contains(hh, "f1")
		req.Contains(hh, "f2")
		req.Contains(hh, "f3")
	})

	t.Run("iterate", func(t *testing.T) {
		dc, err := Decoder(testReader(), "test.csv")
		req.NoError(err)

		aux := make(datasource.RawRecord)
		var more bool

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r1f1", aux["f1"].Values[0])
		req.Equal("r1f2", aux["f2"].Values[0])
		req.Equal("r1f3", aux["f3"].Values[0])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r2f1", aux["f1"].Values[0])
		req.Equal("r2f2", aux["f2"].Values[0])
		req.Equal("r2f3", aux["f3"].Values[0])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r3f1", aux["f1"].Values[0])
		req.Equal("r3f2", aux["f2"].Values[0])
		req.Equal("r3f3", aux["f3"].Values[0])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.False(more)
	})
}

func testReader() io.Reader {
	src := `{"f1": "r1f1", "f2": "r1f2", "f3": "r1f3"}
{"f1": "r2f1", "f2": "r2f2", "f3": "r2f3"}
{"f1": "r3f1", "f2": "r3f2", "f3": "r3f3"}`

	return strings.NewReader(src)
}
