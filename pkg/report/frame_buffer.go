package report

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/qlng"
	"github.com/spf13/cast"
)

type (
	// frameBuffer is a container for pulled frames that we may not use immediately
	frameBuffer struct {
		sourceName string

		loader func(keyFilter *Filter, cap uint) ([]*Frame, error)
		closer Closer

		localFrames   []*Frame
		foreignFrames []*Frame

		// more returns true if we need to pull more
		more      func([]*Frame, []int) bool
		postFetch func([]*Frame) ([]*Frame, error)

		// General state for easier work
		// All of these provided parameters must be validated beforehand; no validation
		// occurs in this struct
		metaInitialized bool
		sorting         filter.SortExprSet
		sortColumns     []int

		keyCol      string
		keyColIndex int
	}
)

// load loads the next chunk into the buffer based on the provided loader
func (bl *frameBuffer) load(cap uint, keyFilter *Filter) (more bool, err error) {
	var ff []*Frame
	auxLocal := make([]*Frame, 0, 32)
	auxForeign := make([]*Frame, 0, 32)

	for {
		// Load
		ff, err = bl.loader(keyFilter, cap)
		if err != nil {
			return false, err
		}
		if ff == nil {
			return false, nil
		}

		for _, f := range ff {
			f.Source = bl.sourceName
		}

		// Bucket
		for _, f := range ff {
			if f.RefValue != "" {
				auxForeign = append(auxForeign, f)
				continue
			}
			auxLocal = append(auxLocal, f)
		}

		// Post processing
		if bl.postFetch != nil {
			auxLocal, err = bl.postFetch(auxLocal)
			if err != nil {
				return false, err
			}
		}

		// Push into the buffer
		if len(bl.localFrames) == 0 {
			bl.localFrames = auxLocal
		} else {
			// - when the source is partitioned we append frames to the buffer
			if len(auxLocal) > 1 {
				bl.localFrames = append(bl.localFrames, auxLocal...)
			} else {
				// - when the source is not partitioned, we append the fetched rows to the only buffer frame
				bl.localFrames[0].Rows = append(bl.localFrames[0].Rows, auxLocal[0].Rows...)
			}
		}

		// foreign frames should not overwrite eachothers as they are all complete
		bl.foreignFrames = append(bl.foreignFrames, auxForeign...)

		// Initialize meta parameters as they are dependant on the resulting frames
		if !bl.metaInitialized {
			bl.metaInitialized = true
			bl.sortColumns = make([]int, len(bl.sorting))
			for i, s := range bl.sorting {
				bl.sortColumns[i] = bl.localFrames[0].Columns.Find(s.Column)
			}

			bl.keyColIndex = bl.localFrames[0].Columns.Find(bl.keyCol)
		}

		// Do we need to fetch more?
		// When chunk size is 0, we are fetching all
		if cap == 0 || !bl.more(bl.localFrames, bl.sortColumns) {
			return true, nil
		}
	}
}

// keys returns a slice of keys that should be used for filtering over the
// sub loader.
func (bl *frameBuffer) keys() (keys []string, err error) {
	keys = make([]string, 0, defaultPageSize)
	keySet := make(map[string]bool)
	var k string

	for _, mf := range bl.localFrames {
		err = mf.WalkRows(func(i int, r FrameRow) error {
			k, err = cast.ToStringE(r[bl.keyColIndex].Get())
			if ok := keySet[k]; !ok {
				keys = append(keys, k)
				keySet[k] = true
			}
			return err
		})
		if err != nil {
			return
		}
	}

	return
}

// keyFilter prepares the filter that should be used when fetching related rows.
//
// @todo do some compression, ie "id > x && id < y"
//       this will return more stuff but it could be faster then the current thing
func (bl *frameBuffer) keyFilter(keys []string) *Filter {
	aa := make(qlng.ASTNodeSet, len(keys))

	for i, k := range keys {
		aa[i] = &qlng.ASTNode{
			Ref: "eq",
			Args: qlng.ASTNodeSet{
				&qlng.ASTNode{Symbol: bl.keyCol},
				&qlng.ASTNode{Value: qlng.MakeValueOf("String", k)},
			},
		}
	}

	return &Filter{
		ASTNode: &qlng.ASTNode{
			Ref: "group",
			Args: qlng.ASTNodeSet{
				&qlng.ASTNode{
					Ref:  "or",
					Args: aa,
				},
			},
		},
	}
}

func (bl *frameBuffer) walkRowsLocal(fnc func(i int, r FrameRow) error) error {
	if len(bl.localFrames) == 0 {
		return nil
	}

	if len(bl.localFrames) > 1 {
		for i, f := range bl.localFrames {
			fnc(i, f.FirstRow())
		}
		return nil
	}

	return bl.localFrames[0].WalkRows(fnc)
}

func (bl *frameBuffer) getByRefValue(ref expr.TypedValue) *Frame {
	s, err := expr.CastToString(ref.Get())
	if err != nil {
		panic(err)
	}

	for _, f := range bl.localFrames {
		if f.RefValue == s {
			return f
		}
	}
	return nil
}

// cutLocal returns the cap of items from the buffers' localFrames and
// buffers back the remaining items
func (bl *frameBuffer) removeLocal(count int) {
	_, bl.localFrames, _ = bl.cut(bl.localFrames, count)
	return
}

// cutForeign returns the cap of items from the buffers' foreignFrames and
// buffers back the remaining items
func (bl *frameBuffer) removeForeign(count int) {
	_, bl.foreignFrames, _ = bl.cut(bl.foreignFrames, count)
	return
}

// cutLocal returns the cap of items from the buffers' localFrames and
// buffers back the remaining items
func (bl *frameBuffer) cutLocal(cap int) (out []*Frame, more bool) {
	out, bl.localFrames, more = bl.cut(bl.localFrames, cap)
	return
}

// cutForeign returns the cap of items from the buffers' foreignFrames and
// buffers back the remaining items
func (bl *frameBuffer) cutForeign(cap int) (out []*Frame) {
	out, bl.foreignFrames, _ = bl.cut(bl.foreignFrames, cap)
	return
}

func (bl *frameBuffer) sizeLocal() int {
	if len(bl.localFrames) == 0 {
		return 0
	}

	if len(bl.localFrames) > 1 {
		return len(bl.localFrames)
	}

	return len(bl.localFrames[0].Rows)
}

func (bl *frameBuffer) sizeForeign() int {
	return len(bl.foreignFrames)
}

// cut is a generic sub fnc for cutLocal and cutForeign
func (bl *frameBuffer) cut(ff []*Frame, cap int) (out []*Frame, buffer []*Frame, more bool) {
	lfl := len(ff)

	if lfl == 0 {
		return nil, nil, false
	}

	if lfl == 1 {
		o := ff[0].CloneMeta()
		if ff[0].Size() <= cap {
			o.Rows = ff[0].Rows
			return []*Frame{o}, nil, false
		}

		o.Rows = ff[0].Rows[0:cap]
		ff[0].Rows = ff[0].Rows[cap:]
		return []*Frame{o}, ff, true
	}

	if lfl < cap {
		oo := ff
		return oo, nil, false
	}

	return ff[0:cap], ff[cap:], true
}

func (bl *frameBuffer) calculatePagingCursors(out []*Frame, sorting filter.SortExprSet, cursor *filter.PagingCursor, hasNext bool) []*Frame {
	if !hasNext {
		return out
	}

	// index frames by ref
	index := make(map[string][]*Frame)
	for _, o := range out {
		index[o.Ref] = append(index[o.Ref], o)
	}

	// @todo add support for prev cursors.
	//       For now only next page is supported.
	//       Keep a reference of previously seen cursors if you need this functionality.
	// if hasPrev {
	// 	out[0].Paging = &filter.Paging{}
	// 	out[0].Paging.PrevPage = bl.calculatePagingCursor(out[0].FirstRow(), out[0].Columns, index, true, sorting...)
	// 	out[0].Paging.PrevPage.ROrder = true
	// 	out[0].Paging.PrevPage.LThen = !sorting.Reversed()
	// }

	if out[0].Paging == nil {
		out[0].Paging = &filter.Paging{}
	}
	out[0].Paging.NextPage = bl.calculatePagingCursor(out[0].LastRow(), out[0].Columns, index, false, sorting...)
	out[0].Paging.NextPage.LThen = sorting.Reversed()

	return out
}

func (bl *frameBuffer) calculatePagingCursor(r FrameRow, cols FrameColumnSet, index map[string][]*Frame, first bool, cc ...*filter.SortExpr) *filter.PagingCursor {
	// The check for existence should be performed way in advanced so we won't bother here.
	// A unique value is also assured at way before.
	cursor := &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}
	var foreignFrames []*Frame
	var v interface{}

	for _, c := range cc {
		foreignFrames = nil

		if strings.Contains(c.Column, ".") {
			pts := strings.Split(c.Column, ".")
			foreignFrames = index[pts[0]]
			var r FrameRow
			var f *Frame
			if first {
				f = foreignFrames[0]
				r = f.FirstRow()
			} else {
				f = foreignFrames[len(foreignFrames)-1]
				r = f.FirstRow()
			}

			if r[f.Columns.Find(pts[1])] != nil {
				v = r[f.Columns.Find(pts[1])].Get()
			}
			cursor.Set(c.Column, v, c.Descending)
		} else {
			if r[cols.Find(c.Column)] != nil {
				v = r[cols.Find(c.Column)].Get()
			}
			cursor.Set(c.Column, v, c.Descending)
		}
	}

	return cursor
}

// prepareResponse takes the provided buffers, metadata and prepares the result of the step
func prepareResponse(main, sub *frameBuffer, inverted, processed bool, lfd *FrameDefinition, keyColumn string, dscr FrameDescriptionSet) (oo []*Frame, err error) {
	var local, foreign *frameBuffer

	// Determine which one was local/foreign
	if inverted {
		local = sub
		foreign = main
	} else {
		local = main
		foreign = sub
	}

	more := false
	if processed {
		// cut
		// - things from local take priority
		oo, more = local.cutLocal(int(lfd.Paging.Limit))
		oo = append(oo, local.cutForeign(int(lfd.Paging.Limit))...)

		// - followed by things in foreign
		aux, _ := foreign.cutLocal(int(lfd.Paging.Limit))
		for _, a := range aux {
			a.RelSource = oo[0].Ref
		}
		oo = append(oo, aux...)
		oo = append(oo, foreign.cutForeign(int(lfd.Paging.Limit))...)

		// paging
		oo = local.calculatePagingCursors(oo, lfd.Sort, lfd.Paging.PageCursor, more)
	} else {
		// whole
		// - things from local take priority
		oo = local.localFrames
		local.localFrames = []*Frame{}
		oo = append(oo, local.foreignFrames...)
		local.foreignFrames = []*Frame{}

		// - followed by things in foreign
		aux := foreign.localFrames
		for _, a := range aux {
			a.RelSource = oo[0].Ref
		}
		oo = append(oo, aux...)
		foreign.localFrames = []*Frame{}
		oo = append(oo, foreign.foreignFrames...)
		foreign.foreignFrames = []*Frame{}
	}

	if len(oo) == 0 {
		return prepareResponseEmpty(lfd, dscr), nil
	}

	return oo, nil
}

func prepareResponseEmpty(lfd *FrameDefinition, dd FrameDescriptionSet) (oo []*Frame) {
	dscr := dd.FilterByRef(lfd.Ref)[0]

	f := &Frame{
		Name:    lfd.Name,
		Ref:     lfd.Ref,
		Source:  lfd.Source,
		Columns: make(FrameColumnSet, 0, len(lfd.Columns)),
	}

	if len(lfd.Columns) == 0 {
		f.Columns = dscr.Columns
	} else {
		for _, c := range lfd.Columns {
			f.Columns = append(f.Columns, dscr.Columns[dscr.Columns.Find(c.Name)])
		}
	}
	return []*Frame{f}
}
