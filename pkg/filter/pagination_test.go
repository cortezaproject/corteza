package filter

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_cursorEncDec(t *testing.T) {
	var (
		req = require.New(t)

		id = uint64(201244712307261628)

		enc string
		cur = &PagingCursor{}
		dec = &PagingCursor{}
	)

	{
		cur.Set("uint64", id, true)
		cur.Set("string", "foo", false)
		req.Len(cur.values, 2)
		req.Equal(id, cur.values[0])
		req.Equal(fmt.Sprintf("<uint64: %d DESC, string: foo, forward>", id), cur.String())
	}

	{
		enc = cur.Encode()
		req.NotEmpty(enc)
	}
	{
		req.NoError(dec.Decode(enc[1 : len(enc)-1]))
		req.Len(dec.values, 2)
		req.Equal(id, dec.values[0])
		req.Equal("foo", dec.values[1])
		req.Equal(fmt.Sprintf("<uint64: %d DESC, string: foo, forward>", id), cur.String())
	}
}

func Test_cursorValueUnmarshal(t *testing.T) {
	var (
		req = require.New(t)
		pcv = &pagingCursorValue{}
	)

	req.NoError(pcv.UnmarshalJSON([]byte("201244712307261628")))
	req.Equal(pcv.v, uint64(201244712307261628))

	req.NoError(pcv.UnmarshalJSON([]byte("42")))
	req.Equal(pcv.v, uint64(42))

	req.NoError(pcv.UnmarshalJSON([]byte("-42")))
	req.Equal(pcv.v, int64(-42))

	req.NoError(pcv.UnmarshalJSON([]byte("true")))
	req.Equal(pcv.v, true)

	req.NoError(pcv.UnmarshalJSON([]byte("42.42")))
	req.Equal(pcv.v, 42.42)

	req.NoError(pcv.UnmarshalJSON([]byte(`"foo"`)))
	req.Equal(pcv.v, "foo")
}
