package csv

import (
	"io"
	"strings"
	"testing"

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

		aux := make(map[string]string)
		var more bool

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r1f1", aux["f1"])
		req.Equal("r1f2", aux["f2"])
		req.Equal("r1f3", aux["f3"])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r2f1", aux["f1"])
		req.Equal("r2f2", aux["f2"])
		req.Equal("r2f3", aux["f3"])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r3f1", aux["f1"])
		req.Equal("r3f2", aux["f2"])
		req.Equal("r3f3", aux["f3"])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.False(more)
	})
}

func testReader() io.Reader {
	src := `f1,f2,f3
r1f1,r1f2,r1f3
r2f1,r2f2,r2f3
r3f1,r3f2,r3f3`

	return strings.NewReader(src)
}
