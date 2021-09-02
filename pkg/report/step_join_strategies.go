package report

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/qlng"
	"github.com/spf13/cast"
)

type (
	// partialPagingCnd is a wrapper struck for parts of the processed paging cursor
	partialPagingCnd struct {
		filterCut     *qlng.ASTNode
		filterInclude *qlng.ASTNode
		// when ref is "" that means local source
		ref string
	}
)

// strategizeLoad uses the given context to determine what join strategy we should
// use to achieve the correct result as optimally as possible
func (d *joinedDataset) strategizeLoad(ctx context.Context, inverted bool, local, foreign *FrameDefinition) (ml *frameBuffer, sl *frameBuffer, err error) {
	if !inverted {
		return d.stratLocalMain(ctx, local, foreign)
	}

	return d.stratForeignMain(ctx, local, foreign)
}

// stratLocalMain uses the local datasource for the main loader
//
// The strategy is used always, except for when the initial sort is controlled
// by the foreign datasource.
//
// In this strategy, the local datasource controlls the initial sort and the
// filter that should be applied when pulling data from the foreign source.
func (d *joinedDataset) stratLocalMain(ctx context.Context, local, foreign *FrameDefinition) (ml *frameBuffer, sl *frameBuffer, err error) {
	var ok bool

	// Prepare the main loader from the local source
	// Local source can be partitioned to support nested joining.
	if d.partitioned {
		ok, err = (d.local.(PartitionableDatasource)).Partition(d.partitionSize, d.partitionCol)
		if err != nil {
			return
		}
		if !ok {
			err = fmt.Errorf("local datasource is not partitionable: %s", d.local.Name())
			return
		}
	}

	ldr, clsr, err := d.local.Load(ctx, local)
	if err != nil {
		return
	}

	ml = &frameBuffer{
		sourceName:  d.Name(),
		keyCol:      d.def.LocalColumn,
		keyColIndex: -1,

		loader: func(_ *Filter, cap uint) ([]*Frame, error) {
			return ldr(int(cap), false)
		},
		closer:  clsr,
		sorting: local.Sort,

		// Overfetch frames when the last two entries define the same sort.
		// This is required to support paging.
		more: func(ff []*Frame, sc []int) bool {
			// No sorting, we don't care
			if len(local.Sort) == 0 {
				return false
			}

			// With sorting and using primary/unique columns, we don't care
			for _, s := range local.Sort {
				c := local.Columns[local.Columns.Find(s.Column)]
				if c.Primary || c.Unique {
					return false
				}
			}

			// This is the last frame we can pull out, we don't care
			if uint(ff[len(ff)-1].Size()) < local.Paging.Limit {
				return false
			}

			// With sorting and regular'ol columns, we care only if the over-fetched row
			// is the same as the last requested row.
			//
			// If we have multiple frames that means that it is partitioned; check whole frames
			if len(ff) > 1 {
				tmpl := len(ff)
				return ff[tmpl-1].FirstRow().Compare(ff[tmpl-2].FirstRow(), sc...) == 0
			}

			// Else it is a regular'ol regular'ol frame; check rows
			return ff[0].LastRow().Compare(ff[0].LastLastRow(), sc...) == 0
		},
	}

	// Partitioned sources need to be sliced
	if d.partitioned {
		ml.postFetch = func(f []*Frame) ([]*Frame, error) {
			ff, err := d.sliceFrames(f, d.def.LocalColumn, d.def.ForeignColumn)
			for _, f := range ff {
				f.RefValue = ""
				f.RelColumn = ""
			}
			return ff, err
		}
	}

	//
	//

	// Prepare the sub loader from the foreign source
	subDS, ok := d.foreign.(PartitionableDatasource)
	if !ok {
		// @todo allow alternatives also
		err = fmt.Errorf("foreign datasource is not partitionable")
		return
	}

	ffilter := foreign.Filter.Clone()
	sl = &frameBuffer{
		sourceName: d.Name(),
		// meta...
		keyCol:      d.def.ForeignColumn,
		keyColIndex: -1,

		loader: func(keyFilter *Filter, cap uint) ([]*Frame, error) {
			foreign.Filter = merger(ffilter.Clone(), keyFilter, "and")

			ok, err := subDS.Partition(cap, d.def.ForeignColumn)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, fmt.Errorf("foreign datasource is not partitionable: %s", d.foreign.Name())
			}
			loader, closer, err := subDS.Load(ctx, foreign)
			if closer != nil {
				defer closer()
			}
			if err != nil {
				return nil, err
			}
			return loader(0, false)
		},
		sorting: foreign.Sort,

		more: func(_ []*Frame, _ []int) bool {
			return false
		},

		postFetch: func(ff []*Frame) ([]*Frame, error) {
			ff, err = d.sliceFrames(ff, d.def.ForeignColumn, d.def.LocalColumn)
			if err != nil {
				return nil, err
			}

			for i := range ff {
				if ff[i].Name == "" {
					ff[i].Name = foreign.Name
				}
				if ff[i].Source == "" {
					ff[i].Source = foreign.Source
				}
				if ff[i].Ref == "" {
					ff[i].Ref = foreign.Ref
				}
			}

			return ff, nil
		},
	}

	return
}

func (d *joinedDataset) stratForeignMain(ctx context.Context, local, foreign *FrameDefinition) (ml *frameBuffer, sl *frameBuffer, err error) {
	// @todo this will be added at the very and as it's an inverse of the above strategy
	return nil, nil, fmt.Errorf("unable to sort by a joined column")
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Sorting

// shouldSort determines if we need to perform additional sorting based on the
// provided frame definition sorting.
//
// * If no sort was requested; consider it sorted.
// * If all sort expressions use only the local source; consider it sorted.
// * If at least one sort expression references a foreign source; consider it NOT sorted.
func (d *joinedDataset) shouldSort(ss filter.SortExprSet) bool {
	if len(ss) == 0 {
		return false
	}

	for _, s := range ss {
		if strings.Contains(s.Column, ".") {
			return true
		}
	}

	return false
}

func (d *joinedDataset) strategizeSort(main, sub *frameBuffer, inverted bool, lfd *FrameDefinition, keyColumn string) (err error) {
	var local, foreign *frameBuffer

	// Determine which one was local/foreign
	if inverted {
		local = sub
		foreign = main
	} else {
		local = main
		foreign = sub
	}

	// - index foreign frames into buckets; here the foreign sort must be respected
	buckets := make(map[string]int)
	for i, mf := range foreign.localFrames {
		buckets[mf.RefValue] = i
	}

	localSort := make(filter.SortExprSet, 0, len(lfd.Sort))
	for _, s := range lfd.Sort {
		if !strings.Contains(s.Column, ".") {
			localSort = append(localSort, s)
		}
	}

	if len(local.localFrames) == 1 {
		err = d.sortFrameRows(local, localSort, buckets, inverted)
	} else {
		err = d.sortBufferFrames(local, localSort, buckets, inverted)
	}

	if err != nil {
		return
	}

	// - assure bucket's order
	// @todo can we do this in one go alongside the original sort?
	aux := make([]*Frame, len(foreign.localFrames))
	for _, f := range foreign.localFrames {
		aux[buckets[f.RefValue]] = f
	}
	foreign.localFrames = aux

	return
}

func (d *joinedDataset) sortFrameRows(local *frameBuffer, ss filter.SortExprSet, buckets map[string]int, inverted bool) error {
	f := local.localFrames[0]

	// - sort based on bucket index
	var ri FrameRow
	var rj FrameRow

	var bi int
	var bj int

	var ci expr.Comparable
	var ok bool

	sort.SliceStable(f.Rows, func(i, j int) bool {
		ri = f.Rows[i]
		rj = f.Rows[j]
		bi = buckets[cast.ToString(ri[local.keyColIndex].Get())]
		bj = buckets[cast.ToString(rj[local.keyColIndex].Get())]

		if inverted {
			// if bi is before bj, the row bellongs to a complitely different bucket
			// so we shouldn't do anything extra
			if bi < bj {
				return !ss.Reversed()
			}
		}

		// go through the sort definitions and sort based on that
		for si, s := range ss {
			ci, ok = ri[local.sortColumns[si]].(expr.Comparable)
			if !ok {
				return !s.Descending
			}

			r, err := ci.Compare(rj[local.sortColumns[si]])
			if err != nil {
				return s.Descending
			}

			if r != 0 {
				if s.Descending {
					return r > 0
				}
				return r < 0
			}
		}

		return bi < bj
	})

	f.WalkRows(func(i int, r FrameRow) error {
		buckets[cast.ToString(r[local.keyColIndex].Get())] = i
		return nil
	})

	return nil
}

func (d *joinedDataset) sortBufferFrames(local *frameBuffer, ss filter.SortExprSet, buckets map[string]int, inverted bool) error {
	f := local.localFrames[0]

	var ri *Frame
	var rj *Frame

	var bi int
	var bj int

	var ci expr.Comparable
	var ok bool

	sort.SliceStable(local.localFrames, func(i, j int) bool {
		ri = local.localFrames[i]
		rj = local.localFrames[j]
		bi = buckets[cast.ToString(ri.Rows[0][local.keyColIndex].Get())]
		bj = buckets[cast.ToString(rj.Rows[0][local.keyColIndex].Get())]

		if inverted {
			// if bi is before bj, the row bellongs to a complitely different bucket
			// so we shouldn't do anything extra
			if bi < bj {
				return !ss.Reversed()
			}
		}

		// go through the sort definitions and sort based on that
		for si, s := range ss {
			ci, ok = ri.Rows[0][local.sortColumns[si]].(expr.Comparable)
			if !ok {
				return !s.Descending
			}

			r, err := ci.Compare(rj.Rows[0][local.sortColumns[si]])
			if err != nil {
				return s.Descending
			}

			if r != 0 {
				if s.Descending {
					return r > 0
				}
				return r < 0
			}
		}

		return bi < bj
	})

	f.WalkRows(func(i int, r FrameRow) error {
		buckets[cast.ToString(r[local.keyColIndex].Get())] = i
		return nil
	})

	return nil
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Paging

func (d *joinedDataset) strategizePaging(local, foreign *FrameDefinition, inverted bool) (pp []partialPagingCnd, err error) {
	pp, err = d.calculatePagingFilters(local, inverted)
	if err != nil {
		return
	}

	local.Filter = merger(&Filter{pp[0].filterCut}, local.Filter, "and")
	local.Paging.PageCursor = nil

	return
}

// calculatePagingFilters produces additional filtering that should be done
// on the datasource level and/or in the join logic.
//
// The core logic is extracted from store/rdbms/builders/cursor.go
func (d *joinedDataset) calculatePagingFilters(local *FrameDefinition, inverted bool) (partials []partialPagingCnd, err error) {
	if len(local.Paging.PageCursor.Keys()) == 0 {
		return
	}

	var localAppendix *qlng.ASTNode

	var (
		cur = local.Paging.PageCursor

		// baseCndAppx is the initial AST for finding rows that match the sort column
		// It's basically the second part of the wrap condition (if the value equals)
		//
		// The correlated string version is: (%s OR ((%s IS NULL AND %s) OR %s = %s))
		baseCndAppx = func(field string, checkNull bool, value interface{}) *qlng.ASTNode {
			return &qlng.ASTNode{
				Ref: "or",
				Args: qlng.ASTNodeSet{
					&qlng.ASTNode{
						Ref: "and",
						Args: qlng.ASTNodeSet{
							&qlng.ASTNode{
								Ref: "is",
								Args: qlng.ASTNodeSet{
									&qlng.ASTNode{
										Symbol: field,
									},
									&qlng.ASTNode{
										Ref: "null",
									},
								},
							}, &qlng.ASTNode{
								Value: qlng.MakeValueOf("Boolean", checkNull),
							},
						},
					},
					&qlng.ASTNode{
						Ref: "eq",
						Args: qlng.ASTNodeSet{{
							Symbol: field,
						}, {
							// @todo type
							Value: qlng.MakeValueOf("String", value),
						}},
					},
				},
			}
		}

		// baseCnd is the initial AST for filtering over the given sort column
		//
		// The correlated string version is: ((%s IS %s AND %s) OR (%s %s %s))
		baseCnd = func(field string, nullVal *qlng.ASTNode, checkNull bool, compOp string, value interface{}, appendix bool) *qlng.ASTNode {
			pp := strings.Split(field, ".")
			field = pp[len(pp)-1]

			out := &qlng.ASTNode{
				Ref: "or",
				Args: qlng.ASTNodeSet{&qlng.ASTNode{
					Ref: "and",
					Args: qlng.ASTNodeSet{&qlng.ASTNode{
						Ref: "is",
						Args: qlng.ASTNodeSet{{
							Symbol: field,
						}, nullVal},
					}, &qlng.ASTNode{
						Value: qlng.MakeValueOf("Boolean", checkNull),
					}},
				}, &qlng.ASTNode{
					Ref: compOp,
					Args: qlng.ASTNodeSet{{
						Symbol: field,
					}, {
						// @todo type
						Value: qlng.MakeValueOf("String", value),
					}},
				},
				},
			}

			if appendix {
				localAppendix = baseCndAppx(field, checkNull, value)

				return &qlng.ASTNode{
					Ref: "or",
					Args: qlng.ASTNodeSet{
						out,
						localAppendix,
					},
				}
			}

			return out
		}

		// wrapCnd is the conjunction between two paging cursor columns
		//
		// The correlated string version is: (%s OR (((%s IS NULL AND %s) OR %s = %s) AND %s))
		wrapCnd = func(base *qlng.ASTNode, field string, value interface{}, checkNull bool, condition *qlng.ASTNode) *qlng.ASTNode {
			pp := strings.Split(field, ".")
			field = pp[len(pp)-1]

			return &qlng.ASTNode{
				Ref: "or",
				Args: qlng.ASTNodeSet{
					base,
					&qlng.ASTNode{
						Ref: "and",
						Args: qlng.ASTNodeSet{
							&qlng.ASTNode{
								Ref: "or",
								Args: qlng.ASTNodeSet{
									&qlng.ASTNode{
										Ref: "and",
										Args: qlng.ASTNodeSet{
											&qlng.ASTNode{
												Ref: "is",
												Args: qlng.ASTNodeSet{
													&qlng.ASTNode{
														Symbol: field,
													},
													&qlng.ASTNode{
														Ref: "null",
													},
												},
											}, &qlng.ASTNode{
												Value: qlng.MakeValueOf("Boolean", checkNull),
											},
										},
									},
									&qlng.ASTNode{
										Ref: "eq",
										Args: qlng.ASTNodeSet{
											&qlng.ASTNode{
												Symbol: field,
											},
											&qlng.ASTNode{
												// @todo type
												Value: qlng.MakeValueOf("String", value),
											},
										},
									},
								},
							},
							condition,
						},
					},
				},
			}
		}
	)

	var (
		cc = cur.Keys()
		vv = cur.Values()

		ltOp = map[bool]string{
			true:  "lt",
			false: "gt",
		}

		notOp = map[bool]*qlng.ASTNode{
			true:  {Ref: "nnull"},
			false: {Ref: "null"},
		}

		isNull = func(i int, neg bool) bool {
			if (isNil(vv[i]) && !neg) || (!isNil(vv[i]) && neg) {
				return true
			}

			return false
		}
	)

	// Some temporary variables to avoid initialization
	var tmp []string
	var field string

	calculateAST := func(cc []string, vv []interface{}, dsc []bool, cut bool) (cnd *qlng.ASTNode) {
		// going from the last key/column to the 1st one
		for i := len(cc) - 1; i >= 0; i-- {
			// We need to cut off the values that are before the cursor (when ascending)
			// and vice-versa for descending.
			lt := dsc[i]
			if cut && cur.IsROrder() {
				lt = !lt
			}
			op := ltOp[lt]

			tmp = strings.Split(cc[i], ".")
			field = tmp[len(tmp)-1]

			base := baseCnd(field, notOp[!lt], isNull(i, lt), op, vv[i], cut && i == len(cc)-1)

			if cnd == nil {
				cnd = base
			} else {
				cnd = wrapCnd(base, field, vv[i], isNull(i, false), cnd)
			}
		}

		return
	}

	// Edge case where only 1 source is used
	ref := ""
	for j := range cc {
		if j > 0 {
			aa := strings.Split(cc[j-1], ".")
			bb := strings.Split(cc[j], ".")
			if len(aa) > 1 {
				ref = aa[0]
			}
			if len(aa) != len(bb) || aa[0] != bb[0] {
				goto out
			}
		}
	}
	partials = append(partials, partialPagingCnd{
		filterCut: calculateAST(cc, vv, cur.Desc(), false),
		ref:       ref,
	})
	return

out:
	// Create a partial filter for each bit of the cursor
	startIndex := 0
	for j := range cc {
		if j == startIndex {
			continue
		}

		aa := strings.Split(cc[startIndex], ".")
		bb := strings.Split(cc[j], ".")
		if len(aa) != len(bb) || aa[0] != bb[0] {
			aux := partialPagingCnd{
				filterCut: calculateAST(cc[startIndex:j], vv[startIndex:j], cur.Desc()[startIndex:j], startIndex == 0),
			}
			if len(aa) > 1 {
				aux.ref = aa[0]
			}
			if startIndex == 0 {
				aux.filterInclude = localAppendix
			}

			partials = append(partials, aux)

			startIndex = j
		}

		if j == len(cc)-1 {
			aux := partialPagingCnd{
				filterCut: calculateAST(cc[startIndex:], vv[startIndex:], cur.Desc()[startIndex:], startIndex == 0),
			}
			if len(aa) > 1 {
				aux.ref = aa[0]
			}
			if startIndex == 0 {
				aux.filterInclude = localAppendix
			}

			partials = append(partials, aux)
		}
	}
	return
}

// pagingFilter applies additional filtering based on the given page cursor
func (d *joinedDataset) pagingFilter(main, sub *frameBuffer, pp []partialPagingCnd) (modified, satisfied bool) {
	cutSize := 0
	done := false

	if len(pp) <= 1 {
		// Already satisfied by the DS
		return false, true
	}

	inclCondition := pp[0]
	fCondition := pp[1]

	main.walkRowsLocal(func(i int, r FrameRow) error {
		if done {
			return nil
		}

		// Firstly we evaluate if the local row falls in the "danger zone"
		// (if the row was right on the edge of where the paging cursor filter applied)
		if d.eval(inclCondition.filterInclude, r, main.localFrames[0].Columns) {

			// If we are in the "danger zone", we check what foreign frames don't pass
			// the cursor filter.
			//
			// If the foreign frame does not pass it, we should remove it along with the local row.
			if fCondition.filterCut != nil && !d.eval(fCondition.filterCut, sub.getByRefValue(r[main.keyColIndex]).FirstRow(), sub.localFrames[0].Columns) {
				cutSize++
			} else {
				done = true
				return nil
			}
		} else {
			done = true
			return nil
		}

		return nil
	})

	if cutSize > 0 {
		main.removeLocal(cutSize)
		main.removeForeign(cutSize)
		sub.removeLocal(cutSize)
		sub.removeForeign(cutSize)

		if main.sizeLocal() <= cutSize {
			// We removed all of the local buffer so the paging is not yet satisfied
			return true, false
		}

		// We removed the portion of the local buffer, so the paging is satisfied
		return true, true
	}

	return false, true
}
