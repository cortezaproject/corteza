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
		req.Equal(fmt.Sprintf("<uint64: %d DESC, string: foo, [FWD,>]>", id), cur.String())
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
		req.Equal(fmt.Sprintf("<uint64: %d DESC, string: foo, [FWD,>]>", id), cur.String())
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

	req.NoError(pcv.UnmarshalJSON([]byte(`null`)))
	req.Nil(pcv.v)
}

func Test_cursorUnmarshal(t *testing.T) {
	var (
		tt = []struct {
			name   string
			json   string
			cursor *PagingCursor
		}{
			{
				"null",
				`{"K":["StageName","id"],"V":[null,210277506916442116],"D":[false,false],"R":false,"LT":false}`,
				&PagingCursor{
					keys:   []string{"StageName", "id"},
					values: []interface{}{nil, uint64(210277506916442116)},
					desc:   []bool{false, false},
					ROrder: false,
					LThen:  false,
				},
			},
		}
	)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
				cur = &PagingCursor{}
			)

			req.NoError(cur.UnmarshalJSON([]byte(tc.json)))
			req.Equal(cur, tc.cursor)
		})
	}

}
