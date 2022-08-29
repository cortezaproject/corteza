package dal

import (
	"context"
	"fmt"
	"io"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/spf13/cast"
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
		scanRow     *Row
		planned     bool
		filtered    bool

		// Index the attributes for easier lookups later on
		outAttrIndex   map[string]int
		leftAttrIndex  map[string]int
		rightAttrIndex map[string]int

		joinRightAttr AttributeMapping
		joinLeftAttr  AttributeMapping

		rowTester tester

		// Index to keep track of related rows
		relIndex *relIndex

		// Output placeholder for sorted rows
		// @todo consider a generic slice for cases when sorting is not needed.
		//       This will probably save up on memory/time since we don't even need
		//       to pull everything.
		outSorted *btree.Generic[*Row]
		i         int
	}
)

func (xs *joinLeft) init(ctx context.Context) (err error) {
	xs.relIndex = newRelIndex()
	xs.indexAttributes()

	xs.rowTester, err = prepareGenericRowTester(xs.def.Filter)
	if err != nil {
		return
	}

	// @note careful here if you throw routines into the mix; see the NoLocks flag.
	//       Enabling locks does have a performance impact so you might be better off by
	//       constructing multiple of these but then you'll also need to complicate
	//       the .Next methods a bit.
	xs.outSorted = btree.NewGenericOptions[*row](makeRowComparator(xs.filter.OrderBy()...), btree.Options{NoLocks: true})

	xs.joinLeftAttr = xs.def.LeftAttributes[xs.leftAttrIndex[xs.def.On.Left]]
	xs.joinRightAttr = xs.def.RightAttributes[xs.rightAttrIndex[xs.def.On.Right]]

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
	xs.def.Filter.cursor, err = filter.PagingCursorFrom(xs.def.Filter.OrderBy(), v, xs.collectPrimaryAttributes()...)
	if err != nil {
		return
	}

	// Redo row tester
	xs.rowTester, err = prepareGenericRowTester(xs.def.Filter)
	if err != nil {
		return
	}

	// Redo the state
	// @todo adjust based on aggregation plan; reuse buffered, etc.
	xs.relIndex = newRelIndex()
	xs.outSorted = btree.NewGenericOptions[*Row](makeRowComparator(xs.filter.OrderBy()...), btree.Options{NoLocks: true})
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
	pc, err = filter.PagingCursorFrom(xs.def.Filter.OrderBy(), v, xs.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	pc.ROrder = true
	pc.LThen = xs.def.Filter.OrderBy().Reversed()

	return
}

func (xs *joinLeft) ForwardCursor(v ValueGetter) (pc *filter.PagingCursor, err error) {
	pc, err = filter.PagingCursorFrom(xs.def.Filter.OrderBy(), v, xs.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utility methods

// next prepares the next scan row based on the defined join plan
func (xs *joinLeft) next(ctx context.Context) (more bool, err error) {
	more, err = xs.pullNext(ctx)
	if !more || err != nil {
		return more, err
	}

	xs.scanRow, _ = xs.outSorted.GetAt(xs.i)
	xs.i++
	return true, nil
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
		for _, r := range b.rows {
			// Merge the two
			err = mergeRows(xs.def.OutAttributes, r, left, r)
			if err != nil {
				return
			}

			// Assert if we want to keep
			if !xs.keep(ctx, r) {
				continue
			}

			xs.outSorted.Set(r)
		}
	}

	return
}

// getRelatedBuffers returns all of the right rows corresponding to the given left row
func (xs *joinLeft) getRelatedBuffers(l *Row) (out []*relIndexBuffer, ok bool, err error) {
	attrIdent := xs.joinLeftAttr.Identifier()
	attrType := xs.joinLeftAttr.Properties().Type
	var aux *relIndexBuffer

	for c := uint(0); c < l.counters[attrIdent]; c++ {
		v, _ := l.GetValue(attrIdent, c)

		switch attrType.(type) {
		case TypeNumber:
			aux, ok = xs.relIndex.GetInt(cast.ToInt64(v))
			if !ok {
				continue
			}
			out = append(out, aux)
			continue

		case TypeText:
			aux, ok = xs.relIndex.GetString(cast.ToString(v))
			if !ok {
				continue
			}
			out = append(out, aux)
			continue

		case TypeID,
			TypeRef:
			aux, ok = xs.relIndex.GetID(cast.ToUint64(v))
			if !ok {
				continue
			}
			out = append(out, aux)
			continue

		default:
			// @note this should be validated way before
			err = fmt.Errorf("cannot use type %s ad join predicate", attrType.Type())
		}

		return
	}

	return
}

// indexRightRow pushes the provided row onto the rel index
// @todo consider moving most of this logic to the relIndex struct.
func (xs *joinLeft) indexRightRow(r *Row) (err error) {
	attrIdent := xs.joinRightAttr.Identifier()
	attrType := xs.joinRightAttr.Properties().Type

	for i := uint(0); i < r.CountValues()[attrIdent]; i++ {
		v, err := r.GetValue(attrIdent, i)
		if err != nil {
			return err
		}

		// @todo not so sure about this switch; see above coment about moving this out
		switch attrType.(type) {
		case TypeNumber:
			xs.relIndex.AddInt(cast.ToInt64(v), r)
			continue

		case TypeText:
			xs.relIndex.AddString(cast.ToString(v), r)
			continue

		case TypeID,
			TypeRef:
			xs.relIndex.AddID(cast.ToUint64(v), r)
			continue

		default:
			// @note this should be validated way before
			return fmt.Errorf("cannot use type %s as join predicate", attrType.Type())
		}
	}

	return
}

// keep checks if the row should be kept or discarded
func (xs *joinLeft) keep(ctx context.Context, r *Row) bool {
	if xs.rowTester == nil {
		return true
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

func (xs *joinLeft) indexAttributes() {
	xs.outAttrIndex = make(map[string]int)
	for i, a := range xs.def.OutAttributes {
		xs.outAttrIndex[a.Identifier()] = i
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
