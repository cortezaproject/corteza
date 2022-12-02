package filter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	simpleGetter map[string][]any
)

func (g simpleGetter) GetValue(name string, pos uint) (any, error) {
	if g[name] == nil {
		return nil, fmt.Errorf("not found")
	}

	if int(pos) >= len(g[name]) {
		return nil, fmt.Errorf("out of bounds")
	}

	return g[name][pos], nil
}

func (g simpleGetter) CountValues() map[string]uint {
	out := make(map[string]uint)
	for k := range g {
		out[k]++
	}
	return out
}

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

func TestPagingCursorFromValueGetter(t *testing.T) {
	tcc := []struct {
		name      string
		vals      simpleGetter
		ss        SortExprSet
		primaries []string

		out *PagingCursor
		err bool
	}{
		{
			name:      "no pk",
			primaries: []string{},
			vals:      simpleGetter{"k1": {"a"}},
			err:       true,
		},

		{
			name:      "simple without sorting",
			primaries: []string{"k1"},
			vals:      simpleGetter{"k1": {"a"}},
			out: &PagingCursor{
				keys:     []string{"k1"},
				kk:       [][]string{{"k1"}},
				values:   []any{"a"},
				modifier: []string{""},
				desc:     []bool{false},
			},
		},

		{
			name:      "complex without sorting",
			primaries: []string{"k1", "k2"},
			vals:      simpleGetter{"k1": {"a"}, "k2": {"b"}, "something": {10}},
			out: &PagingCursor{
				keys:     []string{"k1", "k2"},
				kk:       [][]string{{"k1"}, {"k2"}},
				values:   []any{"a", "b"},
				modifier: []string{"", ""},
				desc:     []bool{false, false},
			},
		},

		{
			name:      "simple with sorting",
			primaries: []string{"k1"},
			vals:      simpleGetter{"k1": {"a"}, "something": {42}},
			ss:        SortExprSet{{Column: "something"}},
			out: &PagingCursor{
				keys:     []string{"something", "k1"},
				kk:       [][]string{{"something"}, {"k1"}},
				values:   []any{42, "a"},
				modifier: []string{"", ""},
				desc:     []bool{false, false},
			},
		},

		{
			name:      "complex with sorting",
			primaries: []string{"k1", "k2"},
			vals:      simpleGetter{"k1": {"a"}, "k2": {"b"}, "something": {10}},
			ss:        SortExprSet{{Column: "something"}, {Column: "k2"}},
			out: &PagingCursor{
				keys:     []string{"something", "k2", "k1"},
				kk:       [][]string{{"something"}, {"k2"}, {"k1"}},
				values:   []any{10, "b", "a"},
				modifier: []string{"", "", ""},
				desc:     []bool{false, false, false},
			},
		},

		{
			name:      "complex mix and match",
			primaries: []string{"k1", "k2"},
			vals:      simpleGetter{"k1": {"a"}, "k2": {"b"}, "something": {10}, "another_thing": {"qwerty"}},
			ss:        SortExprSet{{Column: "something", Descending: true}, {Column: "k2", Descending: false}, {Column: "another_thing", Descending: true}},
			out: &PagingCursor{
				keys:     []string{"something", "k2", "another_thing", "k1"},
				kk:       [][]string{{"something"}, {"k2"}, {"another_thing"}, {"k1"}},
				values:   []any{10, "b", "qwerty", "a"},
				modifier: []string{"", "", "", ""},
				desc:     []bool{true, false, true, false},
				LThen:    true,
			},
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			out, err := PagingCursorFrom(c.ss, c.vals, c.primaries...)
			if c.err {
				require.Error(t, err)
			}
			require.Equal(t, c.out, out)
		})
	}
}

func TestPagingCursor_ToAST(t *testing.T) {
	t.Skip("TODO")
}
