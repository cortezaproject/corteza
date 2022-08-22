package dal

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func TestStepLinkleft(t *testing.T) {
	crs1 := &filter.PagingCursor{}
	crs1.Set("l_pk", 1, false)
	crs1.Set("l_val", "l1 v1", false)
	crs1.Set("f_pk", 2, false)
	crs1.Set("f_fk", 1, false)
	crs1.Set("f_val", "f2 v1", false)

	basicLeftAttrs := []simpleAttribute{
		{ident: "l_pk", t: TypeID{}},
		{ident: "l_val", t: TypeText{}},
	}
	basicRightAttrs := []simpleAttribute{
		{ident: "f_pk", t: TypeID{}},
		{ident: "f_fk", t: TypeRef{}},
		{ident: "f_val", t: TypeText{}},
	}
	basicOutLeftAttrs := []simpleAttribute{
		{ident: "l_pk", t: TypeID{}},
		{ident: "l_val", t: TypeText{}},
	}
	basicOutRightAttrs := []simpleAttribute{
		{ident: "f_pk", t: TypeID{}},
		{ident: "f_fk", t: TypeRef{}},
		{ident: "f_val", t: TypeText{}},
	}

	tcc := []struct {
		name string

		leftAttributes     []simpleAttribute
		rightAttributes    []simpleAttribute
		leftOutAttributes  []simpleAttribute
		rightOutAttributes []simpleAttribute
		linkPred           LinkPredicate

		lIn []simpleRow
		fIn []simpleRow
		out []simpleRow

		f internalFilter
	}{
		// Basic behavior
		{
			name:               "basic link",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
			out: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
		},
		{
			name:               "basic link multiple right",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 1, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 1, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
			},
			out: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 1, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 1, "f_val": "f3 v1"},

				{"l_pk": 2, "l_val": "l2 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
			},
		},
		{
			name:               "basic link omit missing rows",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 9999, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},
		},
		{
			name:               "basic link no rows joined",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},
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
			name:               "basic link empty right",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{},
			out: []simpleRow{},
		},
		{
			name:               "basic link empty left",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},
			lIn:                []simpleRow{},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 123, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 9999, "f_val": "f2 v1"},
			},
			out: []simpleRow{},
		},
		{
			name:               "empty input",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},
			lIn:                []simpleRow{},
			fIn:                []simpleRow{},
			out:                []simpleRow{},
		},

		// Filtering
		{
			name:               "filtering constraints single attr",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_const": "c1", "l_val": "l1 v1"},
				{"l_pk": 2, "l_const": "c2", "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_const": "c1", "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				constraints: map[string][]any{
					"l_const": {"c1"},
				},
			},
		},
		{
			name:               "filtering constraints right single attr",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_const": "c1", "l_val": "l1 v1"},
				{"l_pk": 2, "l_const": "c2", "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_const": "c1", "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				constraints: map[string][]any{
					"f_val": {"f1 v1"},
				},
			},
		},
		{
			name:               "filtering constraints both single attr",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_const": "c1", "l_val": "l1 v1"},
				{"l_pk": 2, "l_const": "c2", "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{},

			f: internalFilter{
				constraints: map[string][]any{
					"l_const": {"c2"},
					"f_val":   {"f1 v1"},
				},
			},
		},
		{
			name:               "filtering constraints multiple attrs",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_const_a": "cac1", "l_const_b": "cbc1", "l_val": "l1 v1"},
				{"l_pk": 2, "l_const_a": "cac1", "l_const_b": "cbc2", "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_const_a": "cac1", "l_const_b": "cbc1", "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				constraints: map[string][]any{"l_const_a": {"cac1"}, "l_const_b": {"cbc1"}},
			},
		},
		{
			name:               "filtering constraints single attr multiple options",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_const_a": "cac1", "l_const_b": "cbc1", "l_val": "l1 v1"},
				{"l_pk": 2, "l_const_a": "cac1", "l_const_b": "cbc2", "l_val": "l2 v1"},
				{"l_pk": 3, "l_const_a": "cac2", "l_const_b": "cbc3", "l_val": "l3 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_const_a": "cac1", "l_const_b": "cbc1", "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_const_a": "cac1", "l_const_b": "cbc2", "l_val": "l2 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			f: internalFilter{
				constraints: map[string][]any{"l_const_b": {"cbc1", "cbc2"}},
			},
		},
		{
			name:               "filtering expressions constant true",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			f: internalFilter{
				expression: "true",
			},
		},
		{
			name:               "filtering expressions constant false",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

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
			name:               "filtering expressions simple",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			f: internalFilter{
				expression: "l_val == 'l2 v1'",
			},
		},
		{
			name:               "filtering expressions right simple",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},

			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				expression: "f_val == 'f1 v1'",
			},
		},

		// Paging
		{
			name:               "paging cut off first entry",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			f: internalFilter{
				cursor: crs1,
			},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 1, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 1, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
			},
			out: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
			},
		},
		{
			name:               "paging cut off first entry with constant true",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},

			f: internalFilter{
				cursor:     crs1,
				expression: "true",
			},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 1, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 1, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
			},
			out: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
			},
		},
		{
			name:               "paging cut off first entry with constant false",
			leftAttributes:     basicLeftAttrs,
			rightAttributes:    basicRightAttrs,
			leftOutAttributes:  basicOutLeftAttrs,
			rightOutAttributes: basicOutRightAttrs,
			linkPred:           LinkPredicate{Left: "l_pk", Right: "f_fk"},
			f: internalFilter{
				cursor:     crs1,
				expression: "false",
			},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 1, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 1, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
			},
			out: []simpleRow{},
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

			def := Link{
				Ident:              "foo",
				On:                 tc.linkPred,
				LeftAttributes:     saToMapping(tc.leftAttributes...),
				RightAttributes:    saToMapping(tc.rightAttributes...),
				OutLeftAttributes:  saToMapping(tc.leftOutAttributes...),
				OutRightAttributes: saToMapping(tc.rightOutAttributes...),
				Filter:             tc.f,
			}

			xs, err := def.Initialize(ctx, l, f)
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

func TestStepLinkleft_cursorCollect_forward(t *testing.T) {
	tcc := []struct {
		name          string
		ss            filter.SortExprSet
		in            simpleRow
		state         simpleRow
		leftAttrs     []simpleAttribute
		rightAttrs    []simpleAttribute
		outleftAttrs  []simpleAttribute
		outrightAttrs []simpleAttribute
		out           func() *filter.PagingCursor
		err           bool
	}{
		{
			name:  "one",
			in:    simpleRow{"f_pk1": 25, "f1": "v25"},
			state: simpleRow{"l_pk1": 1, "f1": "v1"},
			leftAttrs: []simpleAttribute{{
				ident:   "l_pk1",
				primary: true,
			}},
			rightAttrs: []simpleAttribute{{
				ident:   "f_pk1",
				primary: true,
			}},
			outleftAttrs: []simpleAttribute{{
				ident:   "l_pk1",
				primary: true,
			}},
			outrightAttrs: []simpleAttribute{{
				ident:   "f_pk1",
				primary: true,
			}},
			out: func() *filter.PagingCursor {
				pc := &filter.PagingCursor{}
				pc.Set("l_pk1", 1, false)
				pc.Set("f_pk1", 25, false)
				return pc
			},
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			xs := &linkLeft{
				def: Link{
					Filter: internalFilter{
						orderBy: c.ss,
					},
					LeftAttributes:     saToMapping(c.leftAttrs...),
					RightAttributes:    saToMapping(c.rightAttrs...),
					OutLeftAttributes:  saToMapping(c.outleftAttrs...),
					OutRightAttributes: saToMapping(c.outrightAttrs...),
				},
				leftRow: simpleToRow(c.state),
			}

			out, err := xs.ForwardCursor(c.in)
			require.NoError(t, err)

			require.Equal(t, c.out(), out)
		})
	}
}

func TestStepLinkleft_cursorCollect_back(t *testing.T) {
	tcc := []struct {
		name          string
		ss            filter.SortExprSet
		in            simpleRow
		state         simpleRow
		leftAttrs     []simpleAttribute
		rightAttrs    []simpleAttribute
		outleftAttrs  []simpleAttribute
		outrightAttrs []simpleAttribute
		out           func() *filter.PagingCursor
		err           bool
	}{
		{
			name:  "one",
			in:    simpleRow{"f_pk1": 25, "f1": "v25"},
			state: simpleRow{"l_pk1": 1, "f1": "v1"},
			leftAttrs: []simpleAttribute{{
				ident:   "l_pk1",
				primary: true,
			}},
			rightAttrs: []simpleAttribute{{
				ident:   "f_pk1",
				primary: true,
			}},
			outleftAttrs: []simpleAttribute{{
				ident:   "l_pk1",
				primary: true,
			}},
			outrightAttrs: []simpleAttribute{{
				ident:   "f_pk1",
				primary: true,
			}},
			out: func() *filter.PagingCursor {
				pc := &filter.PagingCursor{}
				pc.Set("l_pk1", 1, false)
				pc.Set("f_pk1", 25, false)
				pc.ROrder = true
				return pc
			},
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			jj := &linkLeft{
				def: Link{
					Filter: internalFilter{
						orderBy: c.ss,
					},
					LeftAttributes:     saToMapping(c.leftAttrs...),
					RightAttributes:    saToMapping(c.rightAttrs...),
					OutLeftAttributes:  saToMapping(c.outleftAttrs...),
					OutRightAttributes: saToMapping(c.outrightAttrs...),
				},
				leftRow: simpleToRow(c.state),
			}

			out, err := jj.BackCursor(c.in)
			require.NoError(t, err)

			require.Equal(t, c.out(), out)
		})
	}
}

func TestStepLinkleft_more(t *testing.T) {
	tcc := []struct {
		name          string
		linkPred      LinkPredicate
		leftAttrs     []simpleAttribute
		rightAttrs    []simpleAttribute
		outleftAttrs  []simpleAttribute
		outrightAttrs []simpleAttribute

		lIn  []simpleRow
		fIn  []simpleRow
		out1 []simpleRow
		out2 []simpleRow

		f internalFilter
	}{
		{
			name: "one",
			leftAttrs: []simpleAttribute{{
				ident:   "l_pk",
				primary: true,
				t:       TypeID{},
			}},
			rightAttrs: []simpleAttribute{{
				ident:   "f_pk",
				primary: true,
				t:       TypeID{},
			}, {
				ident:   "f_fk",
				primary: false,
				t:       TypeID{},
			}},
			outleftAttrs: []simpleAttribute{{
				ident:   "l_pk",
				primary: true,
				t:       TypeID{},
			}},
			outrightAttrs: []simpleAttribute{{
				ident:   "f_pk",
				primary: true,
				t:       TypeID{},
			}, {
				ident:   "f_fk",
				primary: false,
				t:       TypeID{},
			}},
			linkPred: LinkPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 1, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 1, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
			},

			out1: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},
			out2: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1"},
				{"f_pk": 4, "f_fk": 2, "f_val": "f4 v1"},
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

			def := Link{
				Ident:              "foo",
				On:                 tc.linkPred,
				LeftAttributes:     saToMapping(tc.leftAttrs...),
				RightAttributes:    saToMapping(tc.rightAttrs...),
				OutLeftAttributes:  saToMapping(tc.outleftAttrs...),
				OutRightAttributes: saToMapping(tc.outrightAttrs...),
				Filter:             tc.f,
			}

			xs, err := def.Initialize(ctx, l, f)
			require.NoError(t, err)

			require.True(t, xs.Next(ctx))
			out := simpleRow{}
			require.NoError(t, xs.Err())
			require.NoError(t, xs.Scan(out))
			require.Equal(t, tc.out1[0], out)

			require.True(t, xs.Next(ctx))
			out = simpleRow{}
			require.NoError(t, xs.Err())
			require.NoError(t, xs.Scan(out))
			require.Equal(t, tc.out1[1], out)

			l.Seek(ctx, 0)
			f.Seek(ctx, 0)
			require.NoError(t, xs.More(0, out))

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

func simpleToRow(in simpleRow) (out *row) {
	out = &row{}
	for k, v := range in {
		out.SetValue(k, 0, v)
	}
	return
}
