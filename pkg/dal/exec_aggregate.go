package dal

import (
	"context"
	"fmt"
	"sort"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/tidwall/btree"
)

type (
	aggregate struct {
		def    Aggregate
		filter internalFilter

		source  Iterator
		err     error
		scanRow *Row
		planned bool

		rowTester tester

		// Keep track of the registered groupIndex so we can easier implement the iterator.
		// Simple index tracking won't be enough, because btree is self-balancing -- things change
		//
		// @todo worth considering using the inmembuffer and fork it around...
		groupIndex *btree.Generic[*aggregateGroup]
		groups     []*aggregateGroup
		i          int
	}
)

// init initializes the execution step's state
func (xs *aggregate) init(ctx context.Context) (err error) {
	// Initialize state variables
	xs.scanRow = xs.initScanRow()
	xs.groupIndex = btree.NewGeneric[*aggregateGroup](xs.compareGroupKeys)
	xs.groups = make([]*aggregateGroup, 0, 128)
	xs.rowTester, err = prepareGenericRowTester(xs.def.Filter)
	if err != nil {
		return
	}

	// Apply the decided aggregation plan
	return xs.applyPlan(ctx)
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Iterator methods

func (xs *aggregate) Next(ctx context.Context) (more bool) {
	// Assure aggregation plan was taken into account
	xs.err = xs.applyPlan(ctx)
	if xs.err != nil {
		return false
	}

	// Go next...
	more, xs.err = xs.next(ctx)
	return
}

// next prepares the next scannable row
//
// The method is adjusted based on the defined aggregation plan.
// The method performs appropriate filtering/sorting (if needed).
func (xs *aggregate) next(ctx context.Context) (more bool, err error) {
	var g *aggregateGroup
	for {
		// Make sure it's cleared out and ready for fresh data
		xs.scanRow.Reset()

		// Make sure we have data
		err = xs.prepareNext(ctx)
		if err != nil {
			return false, err
		}

		// Next group
		g, err = xs.nextGroup(ctx)
		if err != nil || g == nil {
			return false, err
		}

		// Make sure the key is in there
		// @todo potential optimization where we don't scan entire composited keys
		err = xs.scanKey(g, xs.scanRow)
		if err != nil {
			return
		}
		// Scan the data
		err = g.agg.Scan(xs.scanRow)
		if err != nil {
			return
		}

		// Check if we want to keep it
		if !xs.keep(ctx, xs.scanRow) {
			continue
		}

		break
	}

	return true, nil
}

func (xs *aggregate) More(limit uint, v ValueGetter) (err error) {
	// Redo the cursor
	xs.def.Filter.cursor, err = filter.PagingCursorFrom(xs.def.Filter.OrderBy(), v, xs.collectPrimaryAttributes()...)
	if err != nil {
		return
	}

	// Redo the row tester
	xs.rowTester, err = prepareGenericRowTester(xs.def.Filter)
	if err != nil {
		return
	}

	// Redo the state
	// @todo adjust based on aggregation plan; reuse buffered, etc.
	xs.scanRow.Reset()
	xs.groupIndex = btree.NewGeneric[*aggregateGroup](xs.compareGroupKeys)
	xs.groups = make([]*aggregateGroup, 0, 128)
	xs.planned = false

	return
}

func (s *aggregate) Err() error { return s.err }

func (s *aggregate) Scan(dst ValueSetter) (err error) {
	if s.i < 0 {
		return fmt.Errorf("@todo err not initialized; next first")
	}

	var v any
	for name, cc := range s.scanRow.CountValues() {
		for i := uint(0); i < cc; i++ {
			// omitting err here since it won't happen
			v, _ = s.scanRow.GetValue(name, i)
			err = dst.SetValue(name, i, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *aggregate) Close() error {
	if s.source != nil {
		return s.source.Close()
	}
	return nil
}

func (s *aggregate) BackCursor(v ValueGetter) (*filter.PagingCursor, error) {
	c, err := filter.PagingCursorFrom(s.def.Filter.OrderBy(), v, s.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	c.ROrder = true
	c.LThen = s.def.Filter.OrderBy().Reversed()

	return c, nil
}

func (s *aggregate) ForwardCursor(v ValueGetter) (*filter.PagingCursor, error) {
	c, err := filter.PagingCursorFrom(s.def.Filter.OrderBy(), v, s.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities

// applyPlan runs additional initialization code based on the defined aggregation plan
func (xs *aggregate) applyPlan(ctx context.Context) (err error) {
	if xs.planned || xs.err != nil {
		return xs.err
	}

	xs.planned = true
	switch {
	case !xs.def.plan.partialScan:
		return xs.pullEntireSource(ctx)
	}

	return
}

func (xs *aggregate) prepareNext(ctx context.Context) (err error) {
	// Pull next chunk from source if not entirely buffered
	if xs.def.plan.partialScan {
		err = xs.pullNextChunk(ctx)
		if err != nil {
			return
		}
	}
	return
}

// pullEntireSource pulls the entire source into aggregators
//
// The entire source should be pulled when underlaying datasources can't provide
// more appropriately ordered data due to costs.
func (xs *aggregate) pullEntireSource(ctx context.Context) (err error) {
	if xs.source == nil {
		return fmt.Errorf("unable to pull data: no source defined")
	}

	// @todo consider pre-populating the hashmaps
	r := &Row{
		counters: make(map[string]uint),
		values:   make(valueSet),
	}

	// Drain the source
	for xs.source.Next(ctx) {
		err = xs.source.Scan(r)
		if err != nil {
			return
		}

		// Get the key for this row
		// @todo try to reuse key; probably a much simpler thing could work
		k := make(groupKey, len(xs.def.Group))
		err = xs.getGroupKey(ctx, r, k)
		if err != nil {
			return
		}

		// Add the row to the group
		err = xs.addToGroup(ctx, k, r)
		if err != nil {
			return
		}

		r.Reset()
	}
	xs.err = xs.source.Err()
	if xs.err != nil {
		return xs.err
	}

	xs.sortGroups()

	return
}

// pullNextChunk pulls the next chunk into the aggregators
//
// Source should be pulled in chunks when the underlaying datasource can provide
// more appropriately ordered data.
func (xs *aggregate) pullNextChunk(ctx context.Context) (err error) {
	return fmt.Errorf("not implemented")
}

func (s *aggregate) addToGroup(ctx context.Context, key groupKey, r ValueGetter) (err error) {
	if s.groupIndex == nil {
		s.groupIndex = btree.NewGeneric[*aggregateGroup](s.compareGroupKeys)
	}

	// Try to get the existing one; if it doesn't exist, push a new one
	// @todo this causes a bit of a memory pressure; investigate
	g, ok := s.groupIndex.Get(&aggregateGroup{key: key})
	if !ok {
		g, err = s.wrapGroup(ctx, key)
		if err != nil {
			return err
		}
		s.groupIndex.Set(g)

		s.groups = append(s.groups, g)
	}

	// if it exists/was added, add it to the aggregator
	err = g.agg.Aggregate(ctx, r)
	if err != nil {
		return
	}

	return
}

// Group key comparator
// @todo can this be moved to the struct's Less method?
func (t *aggregate) compareGroupKeys(a, b *aggregateGroup) (out bool) {
	for i := range a.key {
		va := a.key[i]
		vb := b.key[i]
		out = out || compareValues(va, vb) < 0
	}

	return
}

func (s *aggregate) getGroupKey(ctx context.Context, r ValueGetter, key groupKey) (err error) {
	for i, attr := range s.def.Group {
		// @todo support expressions?
		v, err := r.GetValue(attr.Expression(), 0)
		if err != nil {
			return err
		}

		// @todo multi-value support?
		key[i] = v
	}

	return nil
}

func (s *aggregate) wrapGroup(ctx context.Context, key groupKey) (g *aggregateGroup, err error) {
	agg, err := Aggregator(s.def.OutAttributes...)
	if err != nil {
		return
	}

	g = &aggregateGroup{
		key: key,
		agg: agg,
	}

	return
}

func (xs *aggregate) nextGroup(ctx context.Context) (_ *aggregateGroup, err error) {
	if xs.i >= len(xs.groups) {
		return
	}

	xs.i++
	return xs.groups[xs.i-1], nil
}

func (s *aggregate) scanKey(g *aggregateGroup, dst *Row) (err error) {
	for i, attr := range s.def.Group {
		// @todo multi value support?
		dst.SetValue(attr.Identifier(), 0, g.key[i])
		// omitting err; internal row won't raise them
	}

	return nil
}

func (s *aggregate) keep(ctx context.Context, r *Row) bool {
	if s.rowTester == nil {
		return true
	}
	return s.rowTester.Test(ctx, r)
}

// Each group key is a PK candidate; all together form a composite key
func (s *aggregate) collectPrimaryAttributes() (out []string) {
	out = make([]string, 0, 2)
	for _, m := range s.def.Group {
		out = append(out, m.Identifier())
	}

	return
}

func (s *aggregate) sortGroups() {
	sort.SliceStable(s.groups, func(i, j int) bool {
		ga := s.groups[i]
		gb := s.groups[j]

		var (
			va any
			vb any
		)

		for _, o := range s.def.Filter.OrderBy() {
			x := inKeys(s.def.Group, o.Column)
			if x > -1 {
				va = ga.key[x]
			} else {
				x := inKeys(s.def.OutAttributes, o.Column)
				va = ga.agg.aggregates[x]
			}

			x = inKeys(s.def.Group, o.Column)
			if x > -1 {
				vb = gb.key[x]
			} else {
				x := inKeys(s.def.OutAttributes, o.Column)
				vb = gb.agg.aggregates[x]
			}

			cmp := compareValues(va, vb)
			if cmp != 0 {
				if o.Descending {
					return cmp > 0
				}
				return cmp < 0
			}
		}
		return false
	})
}

func inKeys(kk []AttributeMapping, ident string) int {
	for i, k := range kk {
		if k.Identifier() == ident {
			return i
		}
	}
	return -1
}

func (xs *aggregate) initScanRow() (out *Row) {
	// base
	out = &Row{
		counters: make(map[string]uint),
		values:   make(valueSet),
	}

	// pre-populate with known attrs
	for _, attr := range append(xs.def.Group, xs.def.OutAttributes...) {
		out.values[attr.Identifier()] = make([]any, 0, 2)
		out.counters[attr.Identifier()] = 0
	}

	return
}
