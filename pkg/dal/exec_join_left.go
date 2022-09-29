package dal

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/tidwall/btree"
)

type (
	// Considerations for optimizations
	// - Rework how values are stored internally.
	//   Potentially rework the pipeline input/output interfaces, how ValueGetter/Setter work.
	// - Skip initial left row scan/sort when sorting not requested.
	// - Rework row struct to use slices instead of hashmaps.
	//   With how values are now the performance gain is not that impactful and the
	//   rework complexity is a bit too high for now.
	// - Use `.More` to filter rows by key.
	// - When data is provided in a satisfactory order, use that to pull data in chunks.
	joinLeft struct {
		def    Join
		filter internalFilter

		leftSource  Iterator
		rightSource Iterator
		err         error
		scanRow     ValueGetter
		planned     bool
		filtered    bool

		rowTester tester

		// Index to keep track of related rows
		relIndex *relIndex

		// Output placeholder for sorted rows
		// @todo consider a generic slice for cases when sorting is not needed.
		//       This will probably save up on memory/time since we don't even need
		//       to pull everything.
		outSorted *btree.Generic[ValueGetter]
		i         int
	}
)

func (xs *joinLeft) init(ctx context.Context, joinPredType Type) (err error) {
	xs.relIndex, err = newRelIndex(joinPredType)
	if err != nil {
		return
	}

	xs.rowTester, err = prepareGenericRowTester(xs.filter)
	if err != nil {
		return
	}

	// @note careful here if you throw routines into the mix; see the NoLocks flag.
	//       Enabling locks does have a performance impact so you might be better off by
	//       constructing multiple of these but then you'll also need to complicate
	//       the .Next methods a bit.
	xs.outSorted = btree.NewGenericOptions[ValueGetter](makeRowComparator(xs.filter.OrderBy()...), btree.Options{NoLocks: true})

	return xs.applyPlan(ctx)
}

func (xs *joinLeft) Next(ctx context.Context) (more bool) {
	xs.err = xs.applyPlan(ctx)
	if xs.err != nil {
		return false
	}

	more, xs.err = xs.next(ctx)
	return
}

func (xs *joinLeft) More(limit uint, v ValueGetter) (err error) {
	xs.filter.cursor, err = filter.PagingCursorFrom(xs.filter.OrderBy(), v, xs.collectPrimaryAttributes()...)
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
	xs.relIndex.Clear()
	xs.outSorted = btree.NewGenericOptions[ValueGetter](makeRowComparator(xs.filter.OrderBy()...), btree.Options{NoLocks: true})
	xs.scanRow = nil
	xs.planned = false
	xs.i = 0

	return
}

func (xs *joinLeft) Err() error { return xs.err }

func (xs *joinLeft) Scan(s ValueSetter) (err error) {
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

func (xs *joinLeft) Close() (err error) {
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

func (xs *joinLeft) BackCursor(v ValueGetter) (pc *filter.PagingCursor, err error) {
	pc, err = filter.PagingCursorFrom(xs.filter.OrderBy(), v, xs.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	pc.ROrder = true
	pc.LThen = xs.filter.OrderBy().Reversed()

	return
}

func (xs *joinLeft) ForwardCursor(v ValueGetter) (pc *filter.PagingCursor, err error) {
	pc, err = filter.PagingCursorFrom(xs.filter.OrderBy(), v, xs.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utility methods

// next prepares the next scan row based on the defined join plan
func (xs *joinLeft) next(ctx context.Context) (more bool, err error) {
	if xs.limitExceeded() {
		return false, nil
	}

	more, err = xs.pullNext(ctx)
	if !more || err != nil {
		return more, err
	}

	xs.scanRow, _ = xs.outSorted.GetAt(xs.i)
	xs.i++
	return true, nil
}

func (xs *joinLeft) limitExceeded() bool {
	return xs.filter.limit > 0 && xs.i >= int(xs.filter.limit)
}

// pullNext pulls additional data so we can produce more
//
// This step may be omitted based on the join plan.
func (xs *joinLeft) pullNext(ctx context.Context) (more bool, err error) {
	// Pull next chunk from source if not entirely buffered
	if xs.def.plan.partialScan {
		// @todo this case is currently not implemented so we're erroring it out
		return false, fmt.Errorf("partialScan join plan strategy not implemented")
	}

	// Check if buffer has more
	if xs.i >= xs.outSorted.Len() {
		return false, nil
	}
	return true, nil
}

// applyPlan runs plan specific logic to prepare the state
func (xs *joinLeft) applyPlan(ctx context.Context) (err error) {
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
func (xs *joinLeft) pullEntireSource(ctx context.Context) (err error) {
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

	return
}

// pullEntireRightSource pulls and indexes all of the right bits
func (xs *joinLeft) pullEntireRightSource(ctx context.Context) (err error) {
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

// pullEntireLeftSource pulls left bits and attempts to do as much of the joining
// work on this stage
func (xs *joinLeft) pullEntireLeftSource(ctx context.Context) (err error) {
	for xs.leftSource.Next(ctx) {
		l := &Row{
			counters: make(map[string]uint),
			values:   make(valueSet),
		}

		err = xs.leftSource.Scan(l)
		if err != nil {
			return
		}

		err = xs.joinRight(ctx, l)
		if err != nil {
			return
		}
	}
	return xs.leftSource.Err()
}

// joinRight finds related right rows for the given left row and matches them up
//
// @note for sorting, we use a b-tree as it's self sorting.
//			 Benchmarking shows that using a slice is negligibly faster if faster at all.
func (xs *joinLeft) joinRight(ctx context.Context, left *Row) (err error) {
	bb, ok, err := xs.getRelatedBuffers(left)
	if err != nil || !ok {
		return
	}

	for _, b := range bb {
		for _, right := range b.rows {
			// Merge the two
			xs.mergeRows(xs.def.OutAttributes, right, left, right)

			// Assert if we want to keep
			k, err := xs.keep(ctx, right)
			if err != nil {
				return err
			}
			if !k {
				continue
			}

			xs.outSorted.Set(right)
		}
	}

	return
}

// getRelatedBuffers returns all of the right rows corresponding to the given left row
func (xs *joinLeft) getRelatedBuffers(l *Row) (out []*relIndexBuffer, ok bool, err error) {
	var aux *relIndexBuffer
	for c := uint(0); c < l.CountValues()[xs.def.On.Left]; c++ {
		// @note internal Row struct never errors
		v, _ := l.GetValue(xs.def.On.Left, c)
		aux, ok = xs.relIndex.Get(v)
		if !ok {
			continue
		}
		out = append(out, aux)
	}

	return
}

// indexRightRow pushes the provided row onto the rel index
// @todo consider moving most of this logic to the relIndex struct.
func (xs *joinLeft) indexRightRow(r *Row) (err error) {
	for i := uint(0); i < r.CountValues()[xs.def.On.Right]; i++ {
		v, err := r.GetValue(xs.def.On.Right, i)
		if err != nil {
			return err
		}

		xs.relIndex.Add(v, r)
	}

	return
}

// keep checks if the row should be kept or discarded
func (xs *joinLeft) keep(ctx context.Context, r *Row) (bool, error) {
	if xs.rowTester == nil {
		return true, nil
	}

	return xs.rowTester.Test(ctx, r)
}

// collectPrimaryAttributes returns all of the attributes of the composited key
//
// For joins, all primary attributes from both of the sides should be in here
// since that is what always uniquely identifies a joined row.
//
// @todo consider applying PK candidates and filter out some of these. I don't
//       think it'll provide much of a performance boost but worth a shot later on.
func (xs *joinLeft) collectPrimaryAttributes() (out []string) {
	out = make([]string, 0, 2)
	for _, m := range xs.def.OutAttributes {
		if m.Properties().IsPrimary {
			out = append(out, m.Identifier())
		}
	}

	return
}

// mergeRows merges the left and right rows based on the provided attrs
func (xs *joinLeft) mergeRows(attrs []AttributeMapping, out, left, right *Row) {
	var (
		identSide int
		srcIdent  string
	)

	for _, attr := range attrs {
		srcIdent = attr.Identifier()
		identSide = xs.identSide(srcIdent)

		if identSide == -1 {
			xs.mergeValuesFrom(srcIdent, out, left)
		} else if identSide == 1 {
			xs.mergeValuesFrom(srcIdent, out, right)
		} else {
			xs.mergeValuesFrom(srcIdent, out, left, right)
		}
	}
}

func (xs *joinLeft) mergeValuesFrom(ident string, out *Row, sources ...*Row) {
	var (
		aux any
	)
	for _, src := range sources {
		for c := uint(0); c < src.CountValues()[ident]; c++ {
			aux, _ = src.GetValue(ident, c)
			out.SetValue(ident, c, aux)
		}
	}
}

// identSide returns -1 if the ident belongs to the left source, 1 if it belongs
// to the right side, and 0 if it's either.
//
// @todo consider adding an additional flag to identify what side it is on.
//       For now, this should be fine.
func (xs *joinLeft) identSide(ident string) int {
	pp := strings.Split(ident, attributeNestingSeparator)
	if len(pp) > 1 {
		if pp[0] == xs.def.RelLeft {
			return -1
		} else if pp[0] == xs.def.RelRight {
			return 1
		}
	}

	return 0
}
