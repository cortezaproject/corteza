package dal

import (
	"context"
	"fmt"
	"io"
	"sort"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/spf13/cast"
)

type (
	linkLeft struct {
		def    Link
		filter internalFilter

		leftSource  Iterator
		rightSource Iterator
		err         error
		scanRow     *Row
		planned     bool
		filtered    bool

		// Index the attributes for easier lookups later on
		outLeftAttrIndex  map[string]int
		outRightAttrIndex map[string]int
		leftAttrIndex     map[string]int
		rightAttrIndex    map[string]int

		linkRightAttr  AttributeMapping
		linkLeftAttr   AttributeMapping
		rightSortAttrs []string

		rowTester tester

		// Buffer to keep track of pulled left rows
		leftRows   []*Row
		relIndex   *relIndex
		keepLeft   bool
		leftIndex  int
		rightIndex int

		// Some helper fields for temporary data
		leftRow       *Row
		rightRow      *Row
		relScanBuffer *relIndexBuffer
	}
)

func (xs *linkLeft) init(ctx context.Context) (err error) {
	// @note the index is keeping track of the right source attributes so we can
	//       simplify the sorting logic.
	xs.relIndex = newRelIndex(xs.rightSortAttrs...)
	xs.indexAttributes()

	// basic sort breadown
	for _, o := range xs.filter.OrderBy() {
		if _, ok := xs.rightAttrIndex[o.Column]; ok {
			xs.rightSortAttrs = append(xs.rightSortAttrs, o.Column)
		}
	}

	xs.linkLeftAttr = xs.def.LeftAttributes[xs.leftAttrIndex[xs.def.On.Left]]
	xs.linkRightAttr = xs.def.RightAttributes[xs.rightAttrIndex[xs.def.On.Right]]

	xs.rowTester, err = prepareGenericRowTester(xs.filter)
	if err != nil {
		return
	}

	return xs.applyPlan(ctx)
}

func (xs *linkLeft) Next(ctx context.Context) (more bool) {
	xs.err = xs.applyPlan(ctx)
	if xs.err != nil {
		return false
	}

	more, xs.err = xs.next(ctx)
	return
}

func (xs *linkLeft) More(limit uint, v ValueGetter) (err error) {
	xs.filter.cursor, err = xs.ForwardCursor(v)
	if err != nil {
		return
	}

	// Redo row tester
	xs.rowTester, err = prepareGenericRowTester(xs.filter)
	if err != nil {
		return
	}

	// Redo the state
	// @todo adjust based on aggregation plan; reuse buffered, etc.
	xs.relIndex = newRelIndex(xs.rightSortAttrs...)
	xs.leftRows = make([]*Row, 0, 128)
	xs.scanRow = nil
	xs.planned = false
	xs.keepLeft = false
	xs.leftIndex = 0
	xs.rightIndex = 0

	return
}

func (xs *linkLeft) Err() error { return xs.err }

func (xs *linkLeft) Scan(s ValueSetter) (err error) {
	for k, cc := range xs.scanRow.CountValues() {
		for i := uint(0); i < cc; i++ {
			// @note internal row won't raise errors so we can safely omit them
			v, _ := xs.scanRow.GetValue(k, i)
			err = s.SetValue(k, i, v)
			if err != nil {
				return
			}
		}
	}

	return
}

func (xs *linkLeft) Close() (err error) {
	if xs == nil {
		return
	}

	cc := []io.Closer{
		xs.leftSource,
		xs.rightSource,
	}

	for _, c := range cc {
		if c != nil {
			err = c.Close()
			if err != nil {
				return err
			}
		}
	}

	return
}

func (xs *linkLeft) BackCursor(v ValueGetter) (pc *filter.PagingCursor, err error) {
	g := &rowLink{
		a: xs.leftRow,
		b: v,
	}
	pc, err = filter.PagingCursorFrom(xs.filter.OrderBy(), g, xs.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	pc.ROrder = true
	pc.LThen = xs.filter.OrderBy().Reversed()

	return
}

func (xs *linkLeft) ForwardCursor(v ValueGetter) (pc *filter.PagingCursor, err error) {
	g := &rowLink{
		a: xs.leftRow,
		b: v,
	}
	pc, err = filter.PagingCursorFrom(xs.filter.OrderBy(), g, xs.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utility methods

// next prepares the next scan row based on the defined link plan
func (xs *linkLeft) next(ctx context.Context) (more bool, err error) {
	more, err = xs.pullNext(ctx)
	if !more || err != nil {
		return more, err
	}

	return xs.nextBuffered()
}

// pullNext pulls additional data so we can produce more
//
// This step may be omitted based on the join plan.
func (xs *linkLeft) pullNext(ctx context.Context) (more bool, err error) {
	// Pull next chunk from source if not entirely buffered
	if xs.def.plan.partialScan {
		return false, fmt.Errorf("partialScan join plan strategy not implemented")
	}

	// Check if buffer has more
	// We need to check for gt because the right bits may still be relevant
	if xs.leftIndex > len(xs.leftRows) {
		return false, nil
	}
	return true, nil
}

// nextBuffered prepares the next scan row from the buffers
func (xs *linkLeft) nextBuffered() (more bool, err error) {
	// keepLeft indicates if we should keep the left row and move onto the next
	// in the right buffers.
	//
	// If we're keeping it, take the next row from the other side, else take the
	// next left one and reset right counters.

	var (
		ok bool
	)

	for {
		if !xs.keepLeft {
			// Go for the next left row
			if xs.leftIndex >= len(xs.leftRows) {
				return
			}
			xs.leftRow = xs.leftRows[xs.leftIndex]
			xs.leftIndex++

			// Go for the next right buffer
			xs.relScanBuffer, ok, err = xs.getRelatedBuffer(xs.leftRow)
			if !ok || err != nil {
				return
			}
			// if len(rel) > 1 {
			// 	// @todo implement this; not entirely sure how it should be so I'll block it for now
			// 	// @todo move this check futher up
			// 	return false, fmt.Errorf("multi-value link predicates not supported")
			// }
			// xs.relScanBuffer = rel[0]
			xs.rightIndex = 0

			xs.scanRow = xs.leftRow
			xs.keepLeft = true
			return true, nil
		}

		// Related buffer done
		if xs.rightIndex >= len(xs.relScanBuffer.rows) {
			xs.keepLeft = false
			continue
		}

		xs.rightRow = xs.relScanBuffer.rows[xs.rightIndex]
		xs.rightIndex++
		xs.scanRow = xs.rightRow
		break
	}
	return true, nil
}

// getRelatedBuffer returns all of the right rows corresponding to the given left row
func (xs *linkLeft) getRelatedBuffer(l *Row) (out *relIndexBuffer, ok bool, err error) {
	attrIdent := xs.linkLeftAttr.Identifier()
	attrType := xs.linkLeftAttr.Properties().Type

	// @todo mv link predicate attrs
	v, _ := l.GetValue(attrIdent, 0)

	switch attrType.(type) {
	case TypeNumber:
		out, ok = xs.relIndex.GetInt(cast.ToInt64(v))
		return

	case TypeText:
		out, ok = xs.relIndex.GetString(cast.ToString(v))
		return

	case TypeID,
		TypeRef:
		out, ok = xs.relIndex.GetID(cast.ToUint64(v))
		return

	default:
		// @note this should be validated way before
		err = fmt.Errorf("cannot use type %s ad join predicate", attrType.Type())
	}

	return
}

// applyPlan runs plan specific logic to prepare the state
func (xs *linkLeft) applyPlan(ctx context.Context) (err error) {
	if xs.planned || xs.err != nil {
		return
	}

	xs.planned = true
	switch {
	case !xs.def.plan.partialScan:
		return xs.pullEntireSource(ctx)
	}

	return
}

// pullEntireSource pulls both sources into memory and indexes them for later use
func (xs *linkLeft) pullEntireSource(ctx context.Context) (err error) {
	// This bit does the filtering so just mark it of as such here
	xs.filtered = true

	// First the right source
	err = xs.pullEntireRightSource(ctx)
	if err != nil {
		return
	}

	// Next the left source
	err = xs.pullEntireLeftSource(ctx)
	if err != nil {
		return
	}

	// Sort if needed
	return xs.sortLeftRows()
}

// pullEntireRightSource pulls and indexes all of the right bits
func (xs *linkLeft) pullEntireRightSource(ctx context.Context) (err error) {
	// Drain the source
	for xs.rightSource.Next(ctx) {
		r := &Row{
			counters: make(map[string]uint),
			values:   make(valueSet),
		}

		err = xs.rightSource.Scan(r)
		if err != nil {
			return
		}

		err = xs.indexRightRow(r)
		if err != nil {
			return
		}
	}
	return xs.rightSource.Err()
}

// pullEntireLeftSource pulls and indexes all of the left bits
func (xs *linkLeft) pullEntireLeftSource(ctx context.Context) (err error) {
	var (
		rel *relIndexBuffer
		ok  bool
	)

	// Drain the source
	for xs.leftSource.Next(ctx) {
		l := &Row{
			counters: make(map[string]uint),
			values:   make(valueSet),
		}

		err = xs.leftSource.Scan(l)
		if err != nil {
			return
		}

		rel, ok, err = xs.getRelatedBuffer(l)
		if err != nil {
			return
		}
		if !ok {
			continue
		}

		if !xs.keep(ctx, l, rel) {
			continue
		}

		xs.leftRows = append(xs.leftRows, l)
	}
	return xs.leftSource.Err()
}

// indexRightRow pushes the provided row onto the rel index
// @todo consider moving most of this logic to the relIndex struct.
func (xs *linkLeft) indexRightRow(r *Row) (err error) {
	attrIdent := xs.linkRightAttr.Identifier()
	attrType := xs.linkRightAttr.Properties().Type

	// @todo mv link predicate attrs; should be prevented higher up for now
	v, err := r.GetValue(attrIdent, 0)
	if err != nil {
		return err
	}

	// @todo not so sure about this switch; see above coment about moving this out
	switch attrType.(type) {
	case TypeNumber:
		xs.relIndex.AddInt(cast.ToInt64(v), r)
		return

	case TypeText:
		xs.relIndex.AddString(cast.ToString(v), r)
		return

	case TypeID,
		TypeRef:
		xs.relIndex.AddID(cast.ToUint64(v), r)
		return

	default:
		// @note this should be validated way before
		return fmt.Errorf("cannot use type %s as link predicate", attrType.Type())
	}

	return
}

// sortLeftRows sorts the left rows into the correct order
//
// Algorithm outline:
// Compare the two left rows based on the defined sort order.
// If the two left rows match up, or we're primarily sorting using right rows,
// check the min/max stats values produced by the index buffer.
//
// If some min/max of a chunk is greater/lesser then the other, no row in the other
// can appear before/after the prior.
func (xs *linkLeft) sortLeftRows() (err error) {
	var (
		leftRelBufferA *relIndexBuffer
		leftRelBufferB *relIndexBuffer

		a, b any
	)

	// Use stable sort just so we don't needlesly messup the initial order if we decide
	// to preserve sorts produced by further steps.
	sort.SliceStable(xs.leftRows, func(i, j int) bool {
		if err != nil {
			return false
		}

		// Prepare the data
		leftRowA := xs.leftRows[i]
		leftRowB := xs.leftRows[j]

		leftRelBufferA, _, err = xs.getRelatedBuffer(leftRowA)
		if err != nil {
			return false
		}
		leftRelBufferB, _, err = xs.getRelatedBuffer(leftRowB)
		if err != nil {
			return false
		}

		for _, s := range xs.filter.OrderBy() {
			if _, ok := xs.leftAttrIndex[s.Column]; ok {
				// This bit here orders based on the left attributes
				less, skip := evalCmpResult(compareGetters(leftRowA, leftRowB, leftRowA.counters, leftRowB.counters, s.Column), s)
				if !skip {
					return less
				}
			} else {
				// This bit here orders based on the right attributes
				// Check the stats of the buffer; make sure to adjust based on direction

				// Use chunk's values
				if !s.Descending {
					a = leftRelBufferA.min[s.Column]
					b = leftRelBufferB.min[s.Column]
				} else {
					a = leftRelBufferA.max[s.Column]
					b = leftRelBufferB.max[s.Column]
				}

				less, skip := evalCmpResult(compareValues(a, b), s)
				if !skip {
					return less
				}
			}
		}

		return false
	})

	return
}

// keep checks if the row should be kept or discarded
//
// Link's keep is a bit more complicated and it looks at the related buffer as well.
func (xs *linkLeft) keep(ctx context.Context, left *Row, buffer *relIndexBuffer) (keep bool) {
	// If no buffer, we won't keep -- left inner join like behavior
	if buffer == nil {
		return false
	}
	// No tester include all ok rows
	if xs.rowTester == nil {
		return true
	}

	ch := &rowLink{a: left}
	for _, ch.b = range buffer.rows {
		if !xs.rowTester.Test(ctx, ch) {
			return false
		}
	}

	return true
}

func (xs *linkLeft) collectPrimaryAttributes() (out []string) {
	out = make([]string, 0, 2)
	for _, m := range xs.def.OutLeftAttributes {
		if m.Properties().IsPrimary {
			out = append(out, m.Identifier())
		}
	}

	for _, m := range xs.def.OutRightAttributes {
		if m.Properties().IsPrimary {
			out = append(out, m.Identifier())
		}
	}

	return
}

func (xs *linkLeft) indexAttributes() {
	xs.outLeftAttrIndex = make(map[string]int)
	for i, a := range xs.def.LeftAttributes {
		xs.outLeftAttrIndex[a.Identifier()] = i
	}
	xs.outRightAttrIndex = make(map[string]int)
	for i, a := range xs.def.RightAttributes {
		xs.outRightAttrIndex[a.Identifier()] = i
	}

	xs.leftAttrIndex = make(map[string]int)
	for i, a := range xs.def.LeftAttributes {
		xs.leftAttrIndex[a.Identifier()] = i
	}

	xs.rightAttrIndex = make(map[string]int)
	for i, a := range xs.def.RightAttributes {
		xs.rightAttrIndex[a.Identifier()] = i
	}
}
