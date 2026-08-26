package filter

import (
	"fmt"
	"testing"
	"time"

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

// A cursor value must come back as the instant it went in. RFC3339 writes a
// zone offset as hours and minutes, so a zone whose offset carries seconds
// loses them; the cursor's equality test then never matches the row it was
// built from, and paging over a block of rows sharing that value skips the
// block ascending and repeats it forever descending.
func TestPagingCursor_timeSurvivesEncoding(t *testing.T) {
	// +00:58:04 — Ljubljana's LMT, the zone Go's zero time lands in on a
	// machine set to Europe/Ljubljana.
	lmt := time.FixedZone("LMT", 58*60+4)

	tcc := []struct {
		name string
		in   time.Time
	}{
		{"go zero time in a sub-minute-offset zone", time.Time{}.In(lmt)},
		{"ordinary time in a sub-minute-offset zone", time.Date(2026, 8, 26, 13, 4, 5, 0, lmt)},
		{"time already in UTC", time.Date(2026, 8, 26, 13, 4, 5, 0, time.UTC)},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			enc := &PagingCursor{}
			enc.Set("created_at", c.in, true)

			raw, err := enc.MarshalJSON()
			require.NoError(t, err)

			dec := &PagingCursor{}
			require.NoError(t, dec.UnmarshalJSON(raw))

			// Timestamps ride the cursor as RFC3339 strings and are handed to
			// the database as such, so what matters is the instant they denote.
			enc0, ok := dec.Values()[0].(string)
			require.True(t, ok, "cursor value decoded as %T, not a string", dec.Values()[0])

			got, err := time.Parse(time.RFC3339Nano, enc0)
			require.NoError(t, err, "cursor value %q is not RFC3339", enc0)

			require.True(t, c.in.Equal(got),
				"cursor shifted the instant: put in %s, encoded as %q, got back %s (%s off)",
				c.in.UTC(), enc0, got.UTC(), got.Sub(c.in))
		})
	}
}
