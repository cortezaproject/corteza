package dal

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func TestStepAggregate(t *testing.T) {
	basicAttrs := []simpleAttribute{
		{ident: "k1"},
		{ident: "k2"},
		{ident: "v1"},
		{ident: "txt"},
	}

	tcc := []struct {
		name string

		group            []simpleAttribute
		outAttributes    []simpleAttribute
		sourceAttributes []simpleAttribute

		in  []simpleRow
		out []simpleRow

		f internalFilter
	}{
		// Basic behavior
		{
			name:             "basic one key group",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "g1", "v1": 10, "txt": "foo"},
				{"k1": "g1", "v1": 20, "txt": "fas"},
				{"k1": "g2", "v1": 15, "txt": "bar"},
			},

			out: []simpleRow{
				{"k1": "g1", "v1": float64(30)},
				{"k1": "g2", "v1": float64(15)},
			},

			f: internalFilter{orderBy: filter.SortExprSet{{Column: "k1"}}},
		},
		{
			name:             "basic multi key group",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
				{"k1": "a", "k2": "b", "v1": float64(6)},

				{"k1": "b", "k2": "a", "v1": float64(51)},
			},

			f: internalFilter{orderBy: filter.SortExprSet{{Column: "k1"}, {Column: "k2"}}},
		},
		{
			name:             "basic expr in value aggregation",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(add(v1, 2))",
			}},

			in: []simpleRow{
				{"k1": "g1", "v1": 10, "txt": "foo"},
				{"k1": "g1", "v1": 20, "txt": "fas"},
				{"k1": "g2", "v1": 15, "txt": "bar"},
			},

			out: []simpleRow{
				{"k1": "g1", "v1": float64(34)},
				{"k1": "g2", "v1": float64(17)},
			},

			f: internalFilter{orderBy: filter.SortExprSet{{Column: "k1"}}},
		},

		// Filtering
		{
			name:             "filtering constraints single attr",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			},
			},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
				{"k1": "a", "k2": "b", "v1": float64(6)},
			},

			f: internalFilter{
				constraints: map[string][]any{"k1": {"a"}},
			},
		},
		{
			name:             "filtering constraints multiple attrs",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "b", "v1": float64(6)},
			},

			f: internalFilter{
				constraints: map[string][]any{"k1": {"a"}, "k2": {"b"}},
			},
		},
		{
			name:             "filtering constraints single attr multiple options",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "b", "v1": 2, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 3, "txt": "fas"},
				{"k1": "c", "k2": "a", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(10)},
				{"k1": "a", "k2": "b", "v1": float64(2)},
				{"k1": "b", "k2": "a", "v1": float64(3)},
			},

			f: internalFilter{
				orderBy:     filter.SortExprSet{{Column: "k1"}, {Column: "k2"}},
				constraints: map[string][]any{"k1": {"a", "b"}},
			},
		},
		{
			name:             "filtering expression simple expression",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
			},

			f: internalFilter{
				expression: "v1 > 10 && v1 < 20",
			},
		},
		{
			name:             "filtering expression constant true",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
				{"k1": "a", "k2": "b", "v1": float64(6)},

				{"k1": "b", "k2": "a", "v1": float64(51)},
			},

			f: internalFilter{
				expression: "true",
				orderBy:    filter.SortExprSet{{Column: "k1"}, {Column: "k2"}},
			},
		},
		{
			name:             "filtering expression constant false",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{},

			f: internalFilter{
				expression: "false",
			},
		},

		// Sorting
		{
			name:             "sorting single key full key asc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "v1": 2, "txt": "fas"},
				{"k1": "b", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "v1": float64(12)},
				{"k1": "b", "v1": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: false}},
			},
		},
		{
			name:             "sorting single key full key dsc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "v1": 2, "txt": "fas"},
				{"k1": "b", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "b", "v1": float64(3)},
				{"k1": "a", "v1": float64(12)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: true}},
			},
		},
		{
			name:             "sorting multiple key full key asc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "b", "v1": 2, "txt": "fas"},
				{"k1": "b", "k2": "c", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(10)},
				{"k1": "a", "k2": "b", "v1": float64(2)},
				{"k1": "b", "k2": "c", "v1": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: false}, {Column: "k2", Descending: false}},
			},
		},
		{
			name:             "sorting multiple key full key dsc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "b", "v1": 2, "txt": "fas"},
				{"k1": "b", "k2": "c", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "b", "k2": "c", "v1": float64(3)},
				{"k1": "a", "k2": "b", "v1": float64(2)},
				{"k1": "a", "k2": "a", "v1": float64(10)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: true}, {Column: "k2", Descending: true}},
			},
		},
		{
			name:             "sorting multiple key full key mixed",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "b", "v1": 2, "txt": "fas"},
				{"k1": "b", "k2": "c", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "b", "v1": float64(2)},
				{"k1": "a", "k2": "a", "v1": float64(10)},
				{"k1": "b", "k2": "c", "v1": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: false}, {Column: "k2", Descending: true}},
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			bootstrapAggregate(t, func(ctx context.Context, t *testing.T, sa *Aggregate, b Buffer) {
				for _, r := range tc.in {
					require.NoError(t, b.Add(ctx, r))
				}
				sa.Ident = tc.name
				sa.SourceAttributes = saToMapping(tc.sourceAttributes...)
				sa.Group = saToMapping(tc.group...)
				sa.OutAttributes = saToMapping(tc.outAttributes...)
				sa.filter = tc.f

				err := sa.init(ctx)
				require.NoError(t, err)
				aa, err := sa.exec(ctx, b)
				require.NoError(t, err)

				i := 0
				for aa.Next(ctx) {
					out := simpleRow{}
					require.NoError(t, aa.Scan(out))
					require.Equal(t, tc.out[i], out)
					i++
				}
				require.NoError(t, aa.Err())
				require.Equal(t, len(tc.out), i)
			})
		})
	}
}

func TestStepAggregate_cursorCollect_forward(t *testing.T) {
	tcc := []struct {
		name          string
		ss            filter.SortExprSet
		in            simpleRow
		group         []simpleAttribute
		outAttributes []simpleAttribute

		out func() *filter.PagingCursor
		err bool
	}{
		{
			name: "simple",
			in:   simpleRow{"pk1": 1, "f1": "v1"},
			group: []simpleAttribute{{
				ident: "pk1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "f1",
			}},
			out: func() *filter.PagingCursor {
				pc := &filter.PagingCursor{}
				pc.Set("pk1", 1, false)
				return pc
			},
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {

			def := Aggregate{
				filter: internalFilter{
					orderBy: c.ss,
				},
				Group:         saToMapping(c.group...),
				OutAttributes: saToMapping(c.outAttributes...),
			}

			out, err := (&aggregate{def: def}).ForwardCursor(c.in)
			require.NoError(t, err)

			require.Equal(t, c.out(), out)
		})
	}
}

func TestStepAggregate_cursorCollect_back(t *testing.T) {
	tcc := []struct {
		name          string
		ss            filter.SortExprSet
		in            simpleRow
		group         []simpleAttribute
		outAttributes []simpleAttribute

		out func() *filter.PagingCursor
		err bool
	}{
		{
			name: "simple",
			in:   simpleRow{"pk1": 1, "f1": "v1"},
			group: []simpleAttribute{{
				ident: "pk1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "f1",
			}},
			out: func() *filter.PagingCursor {
				pc := &filter.PagingCursor{}
				pc.Set("pk1", 1, false)
				pc.ROrder = true
				return pc
			},
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {

			def := Aggregate{
				filter: internalFilter{
					orderBy: c.ss,
				},
				Group:         saToMapping(c.group...),
				OutAttributes: saToMapping(c.outAttributes...),
			}

			out, err := (&aggregate{def: def}).BackCursor(c.in)
			require.NoError(t, err)

			require.Equal(t, c.out(), out)
		})
	}
}

func TestStepAggregate_more(t *testing.T) {
	basicAttrs := []simpleAttribute{
		{ident: "k1"},
		{ident: "k2"},
		{ident: "v1"},
		{ident: "txt"},
	}

	tcc := []struct {
		name string
		in   []simpleRow

		group            []simpleAttribute
		outAttributes    []simpleAttribute
		sourceAttributes []simpleAttribute

		def *Aggregate

		out1 []simpleRow
		out2 []simpleRow
	}{
		{
			name:             "multiple keys",
			sourceAttributes: basicAttrs,
			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out1: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
			},
			out2: []simpleRow{
				{"k1": "a", "k2": "b", "v1": float64(6)},
				{"k1": "b", "k2": "a", "v1": float64(51)},
			},

			def: &Aggregate{},

			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			},
			},
		},
	}

	ctx := context.Background()
	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			buff := InMemoryBuffer()
			for _, r := range tc.in {
				require.NoError(t, buff.Add(ctx, r))
			}

			d := tc.def
			d.Group = saToMapping(tc.group...)
			d.OutAttributes = saToMapping(tc.outAttributes...)
			d.SourceAttributes = saToMapping(tc.sourceAttributes...)
			for _, k := range tc.group {
				d.filter.orderBy = append(d.filter.orderBy, &filter.SortExpr{Column: k.ident})
			}

			err := d.init(ctx)
			require.NoError(t, err)
			aa, err := d.exec(ctx, buff)
			require.NoError(t, err)

			require.True(t, aa.Next(ctx))
			out := simpleRow{}
			require.NoError(t, aa.Err())
			require.NoError(t, aa.Scan(out))
			require.Equal(t, tc.out1[0], out)

			buff.Seek(ctx, 0)
			require.NoError(t, aa.More(0, out))

			i := 0
			for aa.Next(ctx) {
				out := simpleRow{}
				require.NoError(t, aa.Err())
				require.NoError(t, aa.Scan(out))

				require.Equal(t, tc.out2[i], out)

				i++
			}
			require.Equal(t, len(tc.out2), i)
		})
	}
}
