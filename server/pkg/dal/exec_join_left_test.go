package dal

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/filter"
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
		{ident: "l_pk", t: TypeID{}, primary: true},
		{ident: "l_val", t: TypeText{}},
		{ident: "f_pk", t: TypeID{}, primary: true},
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

	type (
		testCase struct {
			name string

			outAttributes   []simpleAttribute
			leftAttributes  []simpleAttribute
			rightAttributes []simpleAttribute
			joinPred        JoinPredicate

			lIn []simpleRow
			fIn []simpleRow
			out []simpleRow

			f internalFilter
		}
	)

	baseBehavior := []testCase{
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
	}
	sorting := []testCase{
		{
			name:            "sorting single key full asc",
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
				orderBy: filter.SortExprSet{{Column: "l_pk", Descending: false}},
			},
		},
		{
			name:            "sorting single key full desc",
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
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "l_pk", Descending: true}},
			},
		},

		{
			name:            "sorting multiple key left first all asc",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}, simpleAttribute{ident: "right_order", t: TypeText{}}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "right_order", t: TypeText{}}),
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "left_order": "a", "l_val": "l1 v1"},
				{"l_pk": 2, "left_order": "b", "l_val": "l2 v1"},
				{"l_pk": 3, "left_order": "b", "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "right_order": "a", "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "right_order": "c", "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "right_order": "b", "f_fk": 3, "f_val": "f3 v1"},
			},
			out: []simpleRow{
				{"l_pk": 1, "left_order": "a", "right_order": "a", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 3, "left_order": "b", "right_order": "b", "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"l_pk": 2, "left_order": "b", "right_order": "c", "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "left_order", Descending: false}, {Column: "right_order", Descending: false}},
			},
		},
		{
			name:            "sorting multiple key left first all desc",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}, simpleAttribute{ident: "right_order", t: TypeText{}}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "right_order", t: TypeText{}}),
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "left_order": "a", "l_val": "l1 v1"},
				{"l_pk": 2, "left_order": "b", "l_val": "l2 v1"},
				{"l_pk": 3, "left_order": "b", "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "right_order": "a", "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "right_order": "c", "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "right_order": "b", "f_fk": 3, "f_val": "f3 v1"},
			},
			out: []simpleRow{
				{"l_pk": 2, "left_order": "b", "right_order": "c", "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 3, "left_order": "b", "right_order": "b", "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"l_pk": 1, "left_order": "a", "right_order": "a", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "left_order", Descending: true}, {Column: "right_order", Descending: true}},
			},
		},
		{
			name:            "sorting multiple key left first asc desc",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}, simpleAttribute{ident: "right_order", t: TypeText{}}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "right_order", t: TypeText{}}),
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "left_order": "a", "l_val": "l1 v1"},
				{"l_pk": 2, "left_order": "b", "l_val": "l2 v1"},
				{"l_pk": 3, "left_order": "b", "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "right_order": "a", "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "right_order": "c", "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "right_order": "b", "f_fk": 3, "f_val": "f3 v1"},
			},
			out: []simpleRow{
				{"l_pk": 1, "left_order": "a", "right_order": "a", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "left_order": "b", "right_order": "c", "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 3, "left_order": "b", "right_order": "b", "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "left_order", Descending: false}, {Column: "right_order", Descending: true}},
			},
		},
		{
			name:            "sorting multiple key left first desc asc",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}, simpleAttribute{ident: "right_order", t: TypeText{}}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "right_order", t: TypeText{}}),
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "left_order": "a", "l_val": "l1 v1"},
				{"l_pk": 2, "left_order": "b", "l_val": "l2 v1"},
				{"l_pk": 3, "left_order": "b", "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "right_order": "a", "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "right_order": "c", "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "right_order": "b", "f_fk": 3, "f_val": "f3 v1"},
			},
			out: []simpleRow{
				{"l_pk": 3, "left_order": "b", "right_order": "b", "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"l_pk": 2, "left_order": "b", "right_order": "c", "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 1, "left_order": "a", "right_order": "a", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "left_order", Descending: true}, {Column: "right_order", Descending: false}},
			},
		},

		{
			name:            "sorting multiple key right first all asc",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}, simpleAttribute{ident: "right_order", t: TypeText{}}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "right_order", t: TypeText{}}),
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "left_order": "a", "l_val": "l1 v1"},
				{"l_pk": 2, "left_order": "c", "l_val": "l2 v1"},
				{"l_pk": 3, "left_order": "b", "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "right_order": "a", "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "right_order": "b", "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "right_order": "b", "f_fk": 3, "f_val": "f3 v1"},
			},
			out: []simpleRow{
				{"l_pk": 1, "left_order": "a", "right_order": "a", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 3, "left_order": "b", "right_order": "b", "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"l_pk": 2, "left_order": "c", "right_order": "b", "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "right_order", Descending: false}, {Column: "left_order", Descending: false}},
			},
		},
		{
			name:            "sorting multiple key right first all desc",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}, simpleAttribute{ident: "right_order", t: TypeText{}}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "right_order", t: TypeText{}}),
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "left_order": "a", "l_val": "l1 v1"},
				{"l_pk": 2, "left_order": "c", "l_val": "l2 v1"},
				{"l_pk": 3, "left_order": "b", "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "right_order": "a", "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "right_order": "b", "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "right_order": "b", "f_fk": 3, "f_val": "f3 v1"},
			},
			out: []simpleRow{
				{"l_pk": 2, "left_order": "c", "right_order": "b", "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 3, "left_order": "b", "right_order": "b", "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"l_pk": 1, "left_order": "a", "right_order": "a", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "right_order", Descending: true}, {Column: "left_order", Descending: true}},
			},
		},
		{
			name:            "sorting multiple key right first asc desc",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}, simpleAttribute{ident: "right_order", t: TypeText{}}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "right_order", t: TypeText{}}),
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "left_order": "a", "l_val": "l1 v1"},
				{"l_pk": 2, "left_order": "c", "l_val": "l2 v1"},
				{"l_pk": 3, "left_order": "b", "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "right_order": "a", "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "right_order": "b", "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "right_order": "b", "f_fk": 3, "f_val": "f3 v1"},
			},
			out: []simpleRow{
				{"l_pk": 1, "left_order": "a", "right_order": "a", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "left_order": "c", "right_order": "b", "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 3, "left_order": "b", "right_order": "b", "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "right_order", Descending: false}, {Column: "left_order", Descending: true}},
			},
		},
		{
			name:            "sorting multiple key right first desc asc",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}, simpleAttribute{ident: "right_order", t: TypeText{}}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "left_order", t: TypeText{}}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "right_order", t: TypeText{}}),
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "left_order": "a", "l_val": "l1 v1"},
				{"l_pk": 2, "left_order": "c", "l_val": "l2 v1"},
				{"l_pk": 3, "left_order": "b", "l_val": "l3 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "right_order": "a", "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "right_order": "b", "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "right_order": "b", "f_fk": 3, "f_val": "f3 v1"},
			},
			out: []simpleRow{
				{"l_pk": 3, "left_order": "b", "right_order": "b", "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"l_pk": 2, "left_order": "c", "right_order": "b", "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 1, "left_order": "a", "right_order": "a", "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "right_order", Descending: true}, {Column: "left_order", Descending: false}},
			},
		},

		{
			name:            "sorting nulls asc",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": nil},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
			out: []simpleRow{
				{"l_pk": 2, "l_val": nil, "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "l_val", Descending: false}},
			},
		},
		{
			name:            "sorting nulls desc",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": nil},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
			out: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 1, "l_val": nil, "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "l_val", Descending: true}},
			},
		},
	}
	filtering := []testCase{
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

		{
			name:            "filtering mv field single edge-case",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "l_const", multivalue: true}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "l_const", multivalue: true}),
			rightAttributes: basicForeignAttrs,

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_const": []string{"a"}, "l_val": "l1 v1"},
				// For this row we'll use a single 'b', combined with a single value in the
				// mv field has an edge-case in the current implementation
				{"l_pk": 2, "l_const": []string{"bbbbbbb"}, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_const": []string{"a"}, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				expression: "'a' IN l_const OR 'b' IN l_const",
			},
		},
		{
			name:            "filtering mv field on left attr",
			outAttributes:   append(basicAttrs, simpleAttribute{ident: "l_const", multivalue: true}),
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "l_const", multivalue: true}),
			rightAttributes: basicForeignAttrs,

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_const": []string{"a", "b"}, "l_val": "l1 v1"},
				{"l_pk": 2, "l_const": []string{"b", "c"}, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "l_const": []string{"a", "b"}, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				expression: "'b' IN l_const AND 'c' NOT IN l_const",
			},
		},
		{
			name:            "filtering mv field on right attr",
			outAttributes:   basicAttrs,
			leftAttributes:  append(basicLocalAttrs, simpleAttribute{ident: "r_const", multivalue: true}),
			rightAttributes: append(basicForeignAttrs, simpleAttribute{ident: "r_const", multivalue: true}),

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "r_const": []string{"a", "b"}, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "r_const": []string{"b", "c"}, "f_val": "f2 v1"},
			},

			out: []simpleRow{
				{"l_pk": 1, "r_const": []string{"a", "b"}, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},

			f: internalFilter{
				expression: "'b' IN r_const AND 'c' NOT IN r_const",
			},
		},
	}
	paging := []testCase{
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
				orderBy: filter.SortExprSet{
					{Column: "l_pk", Descending: false},
					{Column: "l_val", Descending: false},
					{Column: "f_pk", Descending: false},
					{Column: "f_fk", Descending: false},
					{Column: "f_val", Descending: false},
				},
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
				orderBy: filter.SortExprSet{
					{Column: "l_pk", Descending: false},
					{Column: "l_val", Descending: false},
					{Column: "f_pk", Descending: false},
					{Column: "f_fk", Descending: false},
					{Column: "f_val", Descending: false},
				},
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
				orderBy: filter.SortExprSet{
					{Column: "l_pk", Descending: false},
					{Column: "l_val", Descending: false},
					{Column: "f_pk", Descending: false},
					{Column: "f_fk", Descending: false},
					{Column: "f_val", Descending: false},
				},
			},
		},
	}
	nilValues := []testCase{
		{
			name:            "basic nil left join value",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": nil, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
			out: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
		},
		{
			name:            "basic nil right join value",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_pk", Right: "f_fk"},

			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": nil, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
			out: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
		},
	}

	batches := [][]testCase{
		baseBehavior,
		sorting,
		filtering,
		paging,
		nilValues,
	}

	ctx := context.Background()
	for _, batch := range batches {
		for _, tc := range batch {
			t.Run(tc.name, func(t *testing.T) {
				l := InMemoryBuffer()
				for _, r := range tc.lIn {
					require.NoError(t, l.Add(ctx, r))
				}

				f := InMemoryBuffer()
				for _, r := range tc.fIn {
					require.NoError(t, f.Add(ctx, r))
				}

				// @todo please, move this into test cases; this was a cheat
				if len(tc.f.orderBy) == 0 {
					tc.f.orderBy = filter.SortExprSet{
						{Column: "l_pk"},
						{Column: "f_pk"},
					}
				}

				def := Join{
					Ident:           "foo",
					On:              tc.joinPred,
					OutAttributes:   saToMapping(tc.outAttributes...),
					LeftAttributes:  saToMapping(tc.leftAttributes...),
					RightAttributes: saToMapping(tc.rightAttributes...),
					Filter:          tc.f,

					plan: joinPlan{},
				}

				xs, err := def.iterator(ctx, l, f)
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
}

func TestStepJoinValidation(t *testing.T) {
	ctx := context.Background()

	basicLocalAttrs := []simpleAttribute{
		{ident: "l_pk", t: TypeID{}},
		{ident: "l_val", t: TypeText{}},
	}
	basicForeignAttrs := []simpleAttribute{
		{ident: "f_pk", t: TypeID{}},
		{ident: "f_fk", t: TypeRef{}},
		{ident: "f_val", t: TypeText{}},
	}
	basicPred := JoinPredicate{Left: "l_pk", Right: "f_fk"}

	run := func(t *testing.T, attrs []simpleAttribute) (err error) {
		sa := &Join{
			Ident:           "jn",
			LeftAttributes:  saToMapping(basicLocalAttrs...),
			RightAttributes: saToMapping(basicForeignAttrs...),

			On: basicPred,

			OutAttributes: saToMapping(attrs...),
		}

		return sa.dryrun(ctx)
	}

	runP := func(t *testing.T, pred JoinPredicate, attrs []simpleAttribute) (err error) {
		sa := &Join{
			Ident:           "jn",
			LeftAttributes:  saToMapping(basicLocalAttrs...),
			RightAttributes: saToMapping(basicForeignAttrs...),

			On: pred,

			OutAttributes: saToMapping(attrs...),
		}

		return sa.dryrun(ctx)
	}

	runF := func(t *testing.T, f internalFilter, attrs []simpleAttribute) (err error) {
		sa := &Join{
			Ident:           "jn",
			LeftAttributes:  saToMapping(basicLocalAttrs...),
			RightAttributes: saToMapping(basicForeignAttrs...),

			On: basicPred,

			Filter: f,

			OutAttributes: saToMapping(attrs...),
		}

		return sa.dryrun(ctx)
	}

	basicAttrs := []simpleAttribute{
		{ident: "l_pk", t: TypeID{}},
		{ident: "l_val", t: TypeText{}},
		{ident: "f_pk", t: TypeID{}},
		{ident: "f_fk", t: TypeRef{}},
		{ident: "f_val", t: TypeText{}},
	}
	_ = basicAttrs

	t.Run("out ident doesn't exist", func(t *testing.T) {
		basicAttrs := []simpleAttribute{
			{ident: "i_not_real"},
		}

		err := run(t, basicAttrs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "i_not_real")
	})

	t.Run("left predicate doesn't exist", func(t *testing.T) {
		err := runP(t, JoinPredicate{Left: "i_not_exist", Right: "f_fk"}, basicAttrs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "i_not_exist")
	})

	t.Run("right predicate doesn't exist", func(t *testing.T) {
		err := runP(t, JoinPredicate{Left: "l_pk", Right: "i_not_exist"}, basicAttrs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "i_not_exist")
	})

	t.Run("sort ident does not exist", func(t *testing.T) {
		err := runF(t, internalFilter{orderBy: filter.SortExprSet{{Column: "i_not_yes"}}}, basicAttrs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "i_not_yes")
	})
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

			xs, err := def.iterator(ctx, l, f)
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

func TestStepJoin_paging(t *testing.T) {
	basicLeftSrcAttrs := []simpleAttribute{
		{ident: "l_pk", t: TypeID{}},
		{ident: "l_val", t: TypeText{}},
	}
	basicRightSrcAttrs := []simpleAttribute{
		{ident: "f_pk", t: TypeID{}},
		{ident: "f_fk", t: TypeRef{}},
		{ident: "f_val", t: TypeText{}},
	}

	basicAttrs := []simpleAttribute{
		{ident: "l_pk", primary: true},
		{ident: "l_val"},
		{ident: "f_pk", primary: true},
		{ident: "f_fk"},
		{ident: "f_val"},
	}

	tcc := []struct {
		name string

		outAttributes   []simpleAttribute
		leftAttributes  []simpleAttribute
		rightAttributes []simpleAttribute
		joinPred        JoinPredicate

		lIn []simpleRow
		fIn []simpleRow

		f internalFilter

		outF1 []simpleRow
		outF2 []simpleRow
		outB1 []simpleRow
	}{
		{
			name:            "key asc",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLeftSrcAttrs,
			rightAttributes: basicRightSrcAttrs,

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
				{"l_pk": 3, "l_val": "l3 v1"},
				{"l_pk": 4, "l_val": "l4 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
			},

			f: internalFilter{
				limit:   2,
				orderBy: filter.SortExprSet{{Column: "l_pk", Descending: false}, {Column: "f_pk", Descending: false}},
			},

			outF1: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
			outF2: []simpleRow{
				{"l_pk": 3, "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"l_pk": 4, "l_val": "l4 v1", "f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
			},
			outB1: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
		},
		{
			name:            "key desc",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLeftSrcAttrs,
			rightAttributes: basicRightSrcAttrs,

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
				{"l_pk": 3, "l_val": "l3 v1"},
				{"l_pk": 4, "l_val": "l4 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
			},

			f: internalFilter{
				limit:   2,
				orderBy: filter.SortExprSet{{Column: "l_pk", Descending: true}, {Column: "f_pk", Descending: true}},
			},

			outF1: []simpleRow{
				{"l_pk": 4, "l_val": "l4 v1", "f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
				{"l_pk": 3, "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
			},
			outF2: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},
			outB1: []simpleRow{
				{"l_pk": 4, "l_val": "l4 v1", "f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
				{"l_pk": 3, "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
			},
		},

		{
			name:            "val asc",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLeftSrcAttrs,
			rightAttributes: basicRightSrcAttrs,

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
				{"l_pk": 3, "l_val": "l3 v1"},
				{"l_pk": 4, "l_val": "l4 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
			},

			f: internalFilter{
				limit:   2,
				orderBy: filter.SortExprSet{{Column: "f_val", Descending: false}, {Column: "l_pk", Descending: false}, {Column: "f_pk", Descending: false}},
			},

			outF1: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
			outF2: []simpleRow{
				{"l_pk": 3, "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"l_pk": 4, "l_val": "l4 v1", "f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
			},
			outB1: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
			},
		},
		{
			name:            "val desc",
			outAttributes:   basicAttrs,
			leftAttributes:  basicLeftSrcAttrs,
			rightAttributes: basicRightSrcAttrs,

			joinPred: JoinPredicate{Left: "l_pk", Right: "f_fk"},
			lIn: []simpleRow{
				{"l_pk": 1, "l_val": "l1 v1"},
				{"l_pk": 2, "l_val": "l2 v1"},
				{"l_pk": 3, "l_val": "l3 v1"},
				{"l_pk": 4, "l_val": "l4 v1"},
			},
			fIn: []simpleRow{
				{"f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
				{"f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
				{"f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
			},

			f: internalFilter{
				limit:   2,
				orderBy: filter.SortExprSet{{Column: "f_val", Descending: true}, {Column: "l_pk", Descending: true}, {Column: "f_pk", Descending: true}},
			},

			outF1: []simpleRow{
				{"l_pk": 4, "l_val": "l4 v1", "f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
				{"l_pk": 3, "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
			},
			outF2: []simpleRow{
				{"l_pk": 2, "l_val": "l2 v1", "f_pk": 2, "f_fk": 2, "f_val": "f2 v1"},
				{"l_pk": 1, "l_val": "l1 v1", "f_pk": 1, "f_fk": 1, "f_val": "f1 v1"},
			},
			outB1: []simpleRow{
				{"l_pk": 4, "l_val": "l4 v1", "f_pk": 4, "f_fk": 4, "f_val": "f4 v1"},
				{"l_pk": 3, "l_val": "l3 v1", "f_pk": 3, "f_fk": 3, "f_val": "f3 v1"},
			},
		},
	}

	ctx := context.Background()
	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			buffL := InMemoryBuffer()
			for _, r := range tc.lIn {
				require.NoError(t, buffL.Add(ctx, r))
			}
			buffR := InMemoryBuffer()
			for _, r := range tc.fIn {
				require.NoError(t, buffR.Add(ctx, r))
			}

			var d Join

			prep := func(f internalFilter) {
				d = Join{
					Filter:          f,
					On:              tc.joinPred,
					OutAttributes:   saToMapping(tc.outAttributes...),
					LeftAttributes:  saToMapping(tc.leftAttributes...),
					RightAttributes: saToMapping(tc.rightAttributes...),
				}
			}
			check := func(iter Iterator, assert []simpleRow) (first, last simpleRow) {
				i := 0
				for iter.Next(ctx) {
					out := simpleRow{}

					require.NoError(t, iter.Scan(out))

					require.Equal(t, assert[i], out)
					if i == 0 {
						first = out
					}
					last = out
					i++
				}
				require.NoError(t, iter.Err())
				require.Equal(t, len(assert), i)

				return
			}

			f := tc.f
			var (
				first, last simpleRow
			)

			// First page, no cursor
			prep(f)
			aa, err := d.iterator(ctx, buffL, buffR)
			require.NoError(t, err)
			_, last = check(aa, tc.outF1)

			// Second page, cursor
			require.NoError(t, buffL.Seek(ctx, 0))
			require.NoError(t, buffR.Seek(ctx, 0))
			f.cursor, err = aa.ForwardCursor(last)
			require.NoError(t, err)

			prep(f)
			aa, err = d.iterator(ctx, buffL, buffR)
			require.NoError(t, err)
			first, _ = check(aa, tc.outF2)

			// Third page, back, cursor
			require.NoError(t, buffL.Seek(ctx, 0))
			require.NoError(t, buffR.Seek(ctx, 0))
			f.cursor, err = aa.BackCursor(first)
			require.NoError(t, err)

			prep(f)
			aa, err = d.iterator(ctx, buffL, buffR)
			require.NoError(t, err)
			check(aa, tc.outB1)
		})
	}
}

func TestStepJoin_multiValueFields(t *testing.T) {
	basicLocalAttrs := []simpleAttribute{
		{ident: "l_pk", t: TypeID{}},
		{ident: "l_ref", t: TypeID{}},
		{ident: "l_val", t: TypeText{}},
	}
	basicForeignAttrs := []simpleAttribute{
		{ident: "f_pk", t: TypeID{}},
		{ident: "f_fk", t: TypeRef{}},
		{ident: "f_val", t: TypeText{}},
	}

	basicAttrs := append(basicLocalAttrs, basicForeignAttrs...)

	tcc := []struct {
		name string

		outAttributes   []simpleAttribute
		leftAttributes  []simpleAttribute
		rightAttributes []simpleAttribute
		joinPred        JoinPredicate

		lIn []*Row
		fIn []*Row
		out []*Row
	}{
		{
			name: "multiple left keys",

			outAttributes:   basicAttrs,
			leftAttributes:  basicLocalAttrs,
			rightAttributes: basicForeignAttrs,
			joinPred:        JoinPredicate{Left: "l_ref", Right: "f_fk"},

			lIn: []*Row{
				(&Row{}).
					WithValue("l_pk", 0, 1).
					WithValue("l_ref", 0, 1).
					WithValue("l_ref", 1, 2).
					WithValue("l_val", 0, "l1 v1"),
			},
			fIn: []*Row{
				(&Row{}).
					WithValue("f_pk", 0, 1).
					WithValue("f_fk", 0, 1).
					WithValue("f_val", 0, "f1 v1"),

				(&Row{}).
					WithValue("f_pk", 0, 2).
					WithValue("f_fk", 0, 2).
					WithValue("f_val", 0, "f1 v2"),
			},

			out: []*Row{
				(&Row{}).
					WithValue("l_pk", 0, 1).
					WithValue("l_ref", 0, 1).
					WithValue("l_ref", 1, 2).
					WithValue("l_val", 0, "l1 v1").
					// ...
					WithValue("f_pk", 0, 1).
					WithValue("f_fk", 0, 1).
					WithValue("f_val", 0, "f1 v1"),

				(&Row{}).
					WithValue("l_pk", 0, 1).
					WithValue("l_ref", 0, 1).
					WithValue("l_ref", 1, 2).
					WithValue("l_val", 0, "l1 v1").
					// ...
					WithValue("f_pk", 0, 2).
					WithValue("f_fk", 0, 2).
					WithValue("f_val", 0, "f1 v2"),
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

			def := Join{
				Ident:           "foo",
				On:              tc.joinPred,
				OutAttributes:   saToMapping(tc.outAttributes...),
				LeftAttributes:  saToMapping(tc.leftAttributes...),
				RightAttributes: saToMapping(tc.rightAttributes...),
				Filter:          filter.Generic(filter.WithOrderBy(filter.SortExprSet{{Column: "l_pk"}, {Column: "f_pk"}})),

				plan: joinPlan{},
			}

			xs, err := def.iterator(ctx, l, f)
			require.NoError(t, err)

			i := 0
			for xs.Next(ctx) {
				require.NoError(t, xs.Err())
				out := &Row{}
				require.NoError(t, xs.Err())
				require.NoError(t, xs.Scan(out))

				require.Equal(t, tc.out[i], out)

				i++
			}
			require.Equal(t, len(tc.out), i)
		})
	}
}
