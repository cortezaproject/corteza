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

// strategizeLoading initializes and prepares a set of load contexts based on the provided state
func (d *joinedDataset) strategizeLoading(ctx context.Context, inverted bool, local, foreign *FrameDefinition) (localLoader *frameLoadCtx, foreignLoader *frameLoadCtx, err error) {
	if !inverted {
		return d.stratLoadPrimary(ctx, inverted, local, foreign)
	}
	return d.stratLoadInverted(ctx, inverted, local, foreign)
}

// stratLoadPrimary is the generic strategy where the local DS dictates the output
//
// The strategy should be always be used, except when the foreign DS defines the initial sort.
func (d *joinedDataset) stratLoadPrimary(ctx context.Context, inverted bool, local, foreign *FrameDefinition) (localLoader *frameLoadCtx, foreignLoader *frameLoadCtx, err error) {
	// @todo partitioned local when needed

	// - local
	ldr, clsr, err := d.local.Load(ctx, local)
	if err != nil {
		return
	}
	localLoader = &frameLoadCtx{
		loader: ldr,
		closer: clsr,

		sorting:     local.Sort,
		keyCol:      d.def.LocalColumn,
		keyColIndex: -1,
	}

	// - foreign
	//
	// -- only allow partitionable datasources
	pForeignDS, ok := d.foreign.(PartitionableDatasource)
	if !ok {
		// @todo allow alternatives also
		err = fmt.Errorf("foreign datasource is not partitionable")
		return
	}
	//
	// Clone foreign filter so that we don't corrupt the initial one
	ffilter := foreign.Filter.Clone()
	//
	foreignLoader = &frameLoadCtx{
		initLoader: func(cap int, f *Filter) (Loader, Closer, error) {
			foreign.Filter = merger(ffilter.Clone(), f, "and")

			// - partition
			ok, err := pForeignDS.Partition(uint(cap), d.def.ForeignColumn)
			if err != nil {
				return nil, nil, err
			}
			if !ok {
				return nil, nil, fmt.Errorf("foreign datasource is not partitionable: %s", d.foreign.Name())
			}

			return pForeignDS.Load(ctx, foreign)
		},

		sorting:     foreign.Sort,
		keyCol:      d.def.ForeignColumn,
		keyColIndex: -1,
	}

	return
}

// stratLoadInverted makes the foreign DS dictate the output
//
// This strategy should only be used for cases where the foreign DS defines the initial sort.
// When we're using inverted sort, the foreign DS should define the initial sort.
func (d *joinedDataset) stratLoadInverted(ctx context.Context, inverted bool, local, foreign *FrameDefinition) (localLoader *frameLoadCtx, foreignLoader *frameLoadCtx, err error) {
	// - foreign
	//
	// -- only allow partitionable datasources
	pForeignDS, ok := d.foreign.(PartitionableDatasource)
	if !ok {
		// @todo allow alternatives also
		err = fmt.Errorf("foreign datasource is not partitionable")
		return
	}
	//
	// - partition
	ok, err = pForeignDS.Partition(uint(foreign.Paging.Limit), d.def.ForeignColumn)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, fmt.Errorf("foreign datasource is not partitionable: %s", d.foreign.Name())
	}
	//
	// - init laoder
	ldr, clsr, err := d.foreign.Load(ctx, foreign)
	if err != nil {
		return
	}
	foreignLoader = &frameLoadCtx{
		loader: ldr,
		closer: clsr,

		sorting:     foreign.Sort,
		keyCol:      d.def.ForeignColumn,
		keyColIndex: -1,
	}

	// - local
	//
	// Clone foreign filter so that we don't corrupt the initial one
	// @todo partitioned local when needed
	lfilter := local.Filter.Clone()
	localLoader = &frameLoadCtx{
		initLoader: func(cap int, f *Filter) (Loader, Closer, error) {
			local.Filter = merger(lfilter.Clone(), f, "and")

			return d.local.Load(ctx, local)
		},

		sorting:     local.Sort,
		keyCol:      d.def.LocalColumn,
		keyColIndex: -1,
	}

	return
}

// frameLoader is a utility function to load the appropriate data from the given source
//
// Outline:
//  . load a chunk of data
//  . post processing
//  . overfetch rows when required
func (d *joinedDataset) frameLoader(loader Loader, loadCtx *frameLoadCtx, cap int, partitioned bool, post func(ff []*Frame) ([]*Frame, error)) (out []*Frame, err error) {
	var ff []*Frame
	out = make([]*Frame, 0, 32)

	for {
		// . load a chunk of data
		ff, err = loader(cap, false)
		if err != nil || ff == nil {
			return
		}

		// . post processing
		// .. update metadata
		for _, f := range ff {
			f.Source = d.Name()
		}
		//
		// .. additional provided post processing
		if post != nil {
			ff, err = post(ff)
			if err != nil {
				return
			}
		}
		//
		// .. update the output buffer
		if len(out) == 0 {
			out = append(out, ff...)
		} else {
			if !partitioned {
				// When not partitioned, we always pull one frame
				out[0].Rows = append(out[0].Rows, ff[0].Rows...)
			} else {
				// - Append to existing buffer when partitioned.
				//   Pulled chunks can be considered complete so we don't have to worry about that
				out = append(out, ff...)
			}
		}

		// . overfetch rows when required
		if cap == 0 || !d.overfetch(ff, partitioned, ff[0].Columns, loadCtx.sorting) {
			return out, nil
		}
	}
}

// overfetch checks if we need to overfetch data based on the provided sort
//
// We're overfetching when the last rows of the given buffer define the same values
// for the sort columns.
// If we don't overfetch, we may experience issues when applying paging.
func (d *joinedDataset) overfetch(ff []*Frame, partitioned bool, cc FrameColumnSet, ss filter.SortExprSet) bool {
	// No sorting, we don't care
	if len(ss) == 0 {
		return false
	}

	// When sorting over primary/unique columns, we don't care
	for _, s := range ss {
		c := cc[cc.Find(s.Column)]
		if c.Primary || c.Unique {
			return false
		}
	}

	// When the last two things have the same sort column values, we need to overfetch
	var (
		a, b FrameRow
	)

	if !partitioned {
		if ff[0].Size() < 2 {
			return true
		}

		a = ff[0].LastLastRow()
		b = ff[0].LastRow()
	} else {
		if len(ff) < 2 {
			return true
		}

		a = ff[len(ff)-2].FirstRow()
		b = ff[len(ff)-1].FirstRow()
	}

	sci := make([]int, 0, len(ss))
	for _, s := range ss {
		sci = append(sci, cc.Find(s.Column))
	}

	return a.Compare(b, sci...) == 0
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Sorting

// sortFrameBuffers sorts the provided local, foreign frame buffers
//
// Outline:
//  . bucket foreign frames based on their ref value
//  . sort entire local frames when more then one is proviced (covers partitioned local ds)
//  . sort rows for each local frame
//  . assure correct foreign frame order to match local frame rows
func (d *joinedDataset) sortFrameBuffers(localLoader, foreignLoader *frameLoadCtx, localBuffer, foreignBuffer []*Frame, ss filter.SortExprSet, inverted bool) ([]*Frame, []*Frame) {
	// - extract sort expressions for local
	//
	// Foreign frames are chunked and those are already sorted as they should be.
	localSS := make(filter.SortExprSet, 0, len(ss))
	for _, s := range ss {
		if !strings.Contains(s.Column, ".") {
			localSS = append(localSS, s)
		}
	}

	// . bucket foreign frames based on their ref value
	//
	// This will speed up later processing
	buckets := make(map[string]int)
	for i, mf := range foreignBuffer {
		buckets[mf.RefValue] = i
	}

	// . sort entire local frames when more then one is proviced (covers partitioned local ds)
	sort.SliceStable(localBuffer, func(i, j int) bool {
		frameI := localBuffer[i]
		frameDelimiterI := frameI.FirstRow()
		frameJ := localBuffer[j]
		frameDelimiterJ := frameJ.FirstRow()

		// what bucket the frame corresponds to
		bucketI := buckets[cast.ToString(frameDelimiterI[localLoader.keyColIndex].Get())]
		bucketJ := buckets[cast.ToString(frameDelimiterJ[localLoader.keyColIndex].Get())]

		// when inverted, use foreign frames to determine initial sort
		if inverted {
			if bucketI < bucketJ {
				return !ss.Reversed()
			}
		}

		// go through the sort definitions and sort based on that
		for si, s := range localSS {
			ci, ok := frameDelimiterI[localLoader.sortColumns[si]].(expr.Comparable)
			if !ok {
				return !s.Descending
			}

			r, err := ci.Compare(frameDelimiterJ[localLoader.sortColumns[si]])
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

		return bucketI < bucketJ
	})

	// . sort rows for each local frame
	for _, l := range localBuffer {
		sort.SliceStable(l.Rows, func(i, j int) bool {
			rowI := l.Rows[i]
			rowJ := l.Rows[j]

			// what bucket the frame corresponds to
			bucketI := buckets[cast.ToString(rowI[localLoader.keyColIndex].Get())]
			bucketJ := buckets[cast.ToString(rowJ[localLoader.keyColIndex].Get())]

			// when inverted, use foreign frames to determine initial sort
			if inverted {
				if bucketI < bucketJ {
					return !ss.Reversed()
				}
			}

			// go through the sort definitions and sort based on that
			for si, s := range localSS {
				ci, ok := rowI[localLoader.sortColumns[si]].(expr.Comparable)
				if !ok {
					return !s.Descending
				}

				r, err := ci.Compare(rowJ[localLoader.sortColumns[si]])
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

			return bucketI < bucketJ
		})
	}

	// . assure correct foreign frame order to match local frame rows
	//
	// Each foreign frame is on the same index as the corresponding local row.
	// This simplifies later algorithms and removes the need for additional
	// mapping structures.
	for _, l := range localBuffer {
		l.WalkRows(func(i int, r FrameRow) error {
			buckets[cast.ToString(r[localLoader.keyColIndex].Get())] = i
			return nil
		})
	}
	//
	// Update foreign frame order based on reordered buckets
	aux := make([]*Frame, len(foreignBuffer))
	for _, f := range foreignBuffer {
		aux[buckets[f.RefValue]] = f
	}

	return localBuffer, aux
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Paging

// strategizePaging breaks down the given paging cursor into partial conditions and
// applies it to the frame definitions.
func (d *joinedDataset) strategizePaging(local, foreign *FrameDefinition, inverted bool) (pp []partialPagingCnd, err error) {
	// break down the cursor
	pp, err = d.calcualteCursorConditions(local, inverted)
	if err != nil {
		return
	}

	local.Paging.PageCursor = nil

	// apply the "cut" filter to the appropriate DS; foreign when inverted,
	// local otherwise.
	switch pp[0].ref {
	case local.Ref, "":
		local.Filter = merger(&Filter{pp[0].filterCut}, local.Filter, "and")
	case foreign.Ref:
		foreign.Filter = merger(&Filter{pp[0].filterCut}, foreign.Filter, "and")
	}

	return
}

// calcualteCursorConditions processes the given paging cursor and returns a set
// of partial conditions
//
// The core logic is extracted from store/rdbms/builders/cursor.go
//
// Partial cursor conditions are to be used for additional DS filtering and
// additional post filtering when processing the data.
func (d *joinedDataset) calcualteCursorConditions(local *FrameDefinition, inverted bool) (partials []partialPagingCnd, err error) {
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

// applyPagingFilter performs additional filtering based on the provided paging partials
func (d *joinedDataset) applyPagingFilter(buffLocal, buffForeign []*Frame, pp []partialPagingCnd) ([]*Frame, []*Frame, bool, bool, error) {
	var (
		err                 error
		modified, satisfied bool
		cutIndex            int
	)

	// When there is only one partial, it was already handled by the DS.
	// Sorting assures a unique column when working with cursors.
	if len(pp) <= 1 {
		return buffLocal, buffForeign, false, true, nil
	}

	colIndex := make(map[string]FrameColumnSet)
	colIndex[""] = buffLocal[0].Columns
	colIndex[buffLocal[0].Ref] = buffLocal[0].Columns
	colIndex[buffForeign[0].Ref] = buffForeign[0].Columns

	rowIndex := 0
	nextRows := func() (map[string]FrameRow, bool) {
		out := make(map[string]FrameRow)
		more := true

		// local
		if d.partitioned {
			// local can just be ""
			out[""] = buffLocal[rowIndex].FirstRow()
			out[buffLocal[0].Ref] = buffLocal[rowIndex].FirstRow()
		} else {
			// local can just be ""
			out[""] = buffLocal[0].PeekRow(rowIndex)
			out[buffLocal[0].Ref] = buffLocal[0].PeekRow(rowIndex)
		}

		// foreign
		out[buffForeign[0].Ref] = buffForeign[rowIndex].FirstRow()

		rowIndex++
		more = rowIndex < len(buffForeign)

		return out, more
	}

outer:
	for {
		// peek next rows and break if no more
		rr, more := nextRows()

		// check paging partials; break if cursor passes, continue if it does not
		pass := true
		for i, p := range pp {
			n := p.filterInclude
			if i > 0 {
				n = p.filterCut
			}
			tmp := d.eval(n, rr[p.ref], colIndex[p.ref])
			// the first partial is less strict so it only matches the rows that
			// can potentially be included.
			if !tmp && i == 0 {
				break outer
			}

			pass = pass && tmp
			if !pass {
				break
			}
		}

		// if all conditions passed, this is the delimiter for what is included
		if pass {
			break
		}

		cutIndex++
		if !more {
			break
		}
	}

	// cut based on cut index; satisfied when we don't remove everything due to paging
	modified = cutIndex > 0
	satisfied = cutIndex < len(buffForeign)-1
	//
	// - local
	if d.partitioned {
		buffLocal = buffLocal[cutIndex:]
	} else {
		buffLocal[0].Rows = buffLocal[0].Rows[cutIndex:]
	}
	//
	// - foreign
	buffForeign = buffForeign[cutIndex:]

	return buffLocal, buffForeign, modified, satisfied, err
}
