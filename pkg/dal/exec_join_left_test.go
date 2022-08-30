package dal

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func TestStepJoinLocal(t *testing.T) {
	crs1 := &filter.PagingCursor{}
	crs1.Set("l_pk", 1, false)
	crs1.Set("l_val", "l1 v1", false)
	crs1.Set("f_pk", 1, false)
	crs1.Set("f_fk", 1, false)
	crs1.Set("f_val", "f1 v1", false)

	basicAttrs := []simpleAttribute{
		{ident: "l_pk", t: TypeID{}},
		{ident: "l_val", t: TypeText{}},
		{ident: "f_pk", t: TypeID{}},
		{ident: "f_fk", t: TypeRef{}},
		{ident: "f_val", t: TypeText{}},
	}
	basicLocalAttrs := []simpleAttribute{
		{ident: "l_pk", t: TypeID{}},
		{ident: "l_val", t: TypeText{}},
	}
	basicForeignAttrs := []simpleAttribute{
		{ident: "f_pk", t: TypeID{}},
		{ident: "f_fk", t: TypeRef{}},
		{ident: "f_val", t: TypeText{}},
	}

	tcc := []struct {
		name string

		outAttributes   []simpleAttribute
		leftAttributes  []simpleAttribute
		rightAttributes []simpleAttribute
		joinPred        JoinPredicate

		lIn []simpleRow
		fIn []simpleRow
		out []simpleRow

		f internalFilter
	}{
		// Basic behavior
		{
			name:            "basic link",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
			out: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
		},
		{
			name:            "basic link omit missing rows",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 9999, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},
		},
		{
			name:            "basic link no rows joined",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 123, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 9999, "f_val": "f2 v1"},
			},
			out: []simpleRow{},
		},
		{
			name:            "basic link empty foreign",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{},
			out: []simpleRow{},
		},
		{
			name:            "basic link empty local",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 123, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 9999, "f_val": "f2 v1"},
			},

			out: []simpleRow{},
		},
		{
			name:            "empty input",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{},
			fIn: []simpleRow{},
			out: []simpleRow{},
		},

		// Filtering
		{
			name:            "filtering constraints single attr",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "l_const"}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "l_const"}),
			rightAttributes: basicForeignAttrs,

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_const": "c1", "l_val": "l1 v1"},
				{"l_pk": 2, "l_const": "c2", "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_const": "c1", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				constraints: map[string][]any{
					"l_const": {"c1"},
				},
			},
		},
		{
			name:            "filtering constraints multiple attrs",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "l_const_a"}, simpleAttribute{ident: "l_const_b"}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "l_const_a"}, simpleAttribute{ident: "l_const_b"}),
			rightAttributes: basicForeignAttrs,

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_const_a": "cac1", "l_const_b": "cbc1", "l_val": "l1 v1"},
				{"l_pk": 2, "l_const_a": "cac1", "l_const_b": "cbc2", "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{{
				"l_pk":      1,
				"l_const_a": "cac1",
				"l_const_b": "cbc1",
				"l_val":     "l1 v1",
				"f_pk":      1,
				"f_fk":      1,
				"f_val":     "f1 v1",
			}},

			f: internalFilter{
				constraints: map[string][]any{"l_const_a": {"cac1"}, "l_const_b": {"cbc1"}},
			},
		},
		{
			name:            "filtering constraints single attr multiple options",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "l_const_a"}, simpleAttribute{ident: "l_const_b"}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "l_const_a"}, simpleAttribute{ident: "l_const_b"}),
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_const_a": "cac1", "l_const_b": "cbc1", "l_val": "l1 v1"},
				{"l_pk": 2, "l_const_a": "cac1", "l_const_b": "cbc2", "l_val": "l2 v1"},
				{"l_pk": 3, "l_const_a": "cac2", "l_const_b": "cbc3", "l_val": "l3 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{{
				"l_pk":      1,
				"l_const_a": "cac1",
				"l_const_b": "cbc1",
				"l_val":     "l1 v1",
				"f_pk":      1,
				"f_fk":      1,
				"f_val":     "f1 v1",
			}, {
				"l_pk":      2,
				"l_const_a": "cac1",
				"l_const_b": "cbc2",
				"l_val":     "l2 v1",
				"f_pk":      2,
				"f_fk":      2,
				"f_val":     "f2 v1",
			}},

			f: internalFilter{
				constraints: map[string][]any{"l_const_b": {"cbc1", "cbc2"}},
			},
		},
		{
			name:            "filtering expressions constant true",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			f: internalFilter{
				expression: "true",
			},
		},
		{
			name:            "filtering expressions constant false",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{},

			f: internalFilter{
				expression: "false",
			},
		},
		{
			name:            "filtering expressions simple",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{{
				"l_pk":  2,
				"l_val": "l2 v1",
				"f_pk":  2,
				"f_fk":  2,
				"f_val": "f2 v1",
			}},

			f: internalFilter{
				expression: "l_val == 'l2 v1'",
			},
		},

		// Paging
		{
			name:            "paging cut off first entry",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{{
				"l_pk":  2,
				"l_val": "l2 v1",
				"f_pk":  2,
				"f_fk":  2,
				"f_val": "f2 v1",
			}},

			f: internalFilter{
				cursor: crs1,
			},
		},
		{
			name:            "paging cut off last entry with constant true",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{{
				"l_pk":  2,
				"l_val": "l2 v1",
				"f_pk":  2,
				"f_fk":  2,
				"f_val": "f2 v1",
			}},

			f: internalFilter{
				expression: "true",
				cursor:     crs1,
			},
		},
		{
			name:            "paging cut off last entry with constant false",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{},

			f: internalFilter{
				expression: "false",
				cursor:     crs1,
			},
		},
	}

	ctx := context.Background()
	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			l := InMemoryBuffer()
			for _, r := range tc.lIn {
				require.NoError(t, l.Add(ctx, r))
			}

			f := InMemoryBuffer()
			for _, r := range tc.fIn {
				require.NoError(t, f.Add(ctx, r))
			}

			tc.f.orderBy = filter.SortExprSet{
				{Column: "l_pk"},
				{Column: "f_pk"},
			}

			def := Join{
				Ident:           "foo",
				On:              tc.joinPred,
				OutAttributes:   saToMapping(tc.outAttributes...),
				LeftAttributes:  saToMapping(tc.leftAttributes...),
				RightAttributes: saToMapping(tc.rightAttributes...),
				filter:          tc.f,

				plan: joinPlan{},
			}

			err := def.init(ctx)
			require.NoError(t, err)
			xs, err := def.exec(ctx, l, f)
			require.NoError(t, err)

			i := 0
			for xs.Next(ctx) {
				require.NoError(t, xs.Err())
				out := simpleRow{}
				require.NoError(t, xs.Err())
				require.NoError(t, xs.Scan(out))

				require.Equal(t, tc.out[i], out)

				i++
			}
			require.Equal(t, len(tc.out), i)
		})
	}
}

func TestStepJoinLocal_cursorCollect_forward(t *testing.T) {
	tcc := []struct {
		name  string
		ss    filter.SortExprSet
		in    simpleRow
		attrs []simpleAttribute
		out   func() *filter.PagingCursor
		err   bool
	}{
		{
			name: "simple",
			in:   simpleRow{"pk1": 1, "f1": "v1"},
			attrs: []simpleAttribute{{
				ident:   "pk1",
				primary: true,
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

			jj := &joinLeft{
				def: Join{
					filter: internalFilter{
						orderBy: c.ss,
					},
					OutAttributes: saToMapping(c.attrs...),
				},
			}

			out, err := jj.ForwardCursor(c.in)
			require.NoError(t, err)

			require.Equal(t, c.out(), out)
		})
	}
}

func TestStepJoinLocal_cursorCollect_back(t *testing.T) {
	tcc := []struct {
		name  string
		ss    filter.SortExprSet
		in    simpleRow
		attrs []simpleAttribute
		out   func() *filter.PagingCursor
		err   bool
	}{
		{
			name: "simple",
			in:   simpleRow{"pk1": 1, "f1": "v1"},
			attrs: []simpleAttribute{{
				ident:   "pk1",
				primary: true,
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

			jj := &joinLeft{
				def: Join{
					filter: internalFilter{
						orderBy: c.ss,
					},
					OutAttributes: saToMapping(c.attrs...),
				},
			}

			out, err := jj.BackCursor(c.in)
			require.NoError(t, err)

			require.Equal(t, c.out(), out)
		})
	}
}

func TestStepJoinLocal_more(t *testing.T) {
	tcc := []struct {
		name string

		attributes        []simpleAttribute
		localAttributes   []simpleAttribute
		foreignAttributes []simpleAttribute
		joinPred          JoinPredicate

		lIn []simpleRow
		fIn []simpleRow

		out1 []simpleRow
		out2 []simpleRow

		f internalFilter
	}{
		{
			name: "one",
			attributes: []simpleAttribute{
				{ident: "l_pk", primary: true},
				{ident: "l_val"},
				{ident: "f_pk", primary: true},
				{ident: "f_fk"},
				{ident: "f_val"},
			},
			localAttributes: []simpleAttribute{
				{ident: "l_pk", t: TypeID{}},
				{ident: "l_val", t: TypeText{}},
			},
			foreignAttributes: []simpleAttribute{
				{ident: "f_pk", t: TypeID{}},
				{ident: "f_fk", t: TypeRef{}},
				{ident: "f_val", t: TypeText{}},
			},

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
				{"l_pk": 3, "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
			},

			out1: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},
			out2: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 3, "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
			},
		},
	}

	ctx := context.Background()
	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			l := InMemoryBuffer()
			for _, r := range tc.lIn {
				require.NoError(t, l.Add(ctx, r))
			}

			f := InMemoryBuffer()
			for _, r := range tc.fIn {
				require.NoError(t, f.Add(ctx, r))
			}

			tc.f.orderBy = filter.SortExprSet{
				{Column: "l_pk"},
				{Column: "f_pk"},
			}

			def := Join{
				Ident:           "foo",
				On:              tc.joinPred,
				OutAttributes:   saToMapping(tc.attributes...),
				LeftAttributes:  saToMapping(tc.localAttributes...),
				RightAttributes: saToMapping(tc.foreignAttributes...),
				filter:          tc.f,
			}

			err := def.init(ctx)
			require.NoError(t, err)
			xs, err := def.exec(ctx, l, f)
			require.NoError(t, err)

			require.True(t, xs.Next(ctx))
			out := simpleRow{}
			require.NoError(t, xs.Err())
			require.NoError(t, xs.Scan(out))
			require.Equal(t, tc.out1[0], out)

			require.NoError(t, xs.More(0, out))

			l.Seek(ctx, 0)
			f.Seek(ctx, 0)

			i := 0
			for xs.Next(ctx) {
				out := simpleRow{}
				require.NoError(t, xs.Err())
				require.NoError(t, xs.Scan(out))

				require.Equal(t, tc.out2[i], out)

				i++
			}
			require.Equal(t, len(tc.out2), i)
		})
	}
}
