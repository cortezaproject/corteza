package dal

import (
	"context"
	"fmt"
	"sort"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/ql"
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

		groupDefs     []AggregateAttr
		aggregateDefs []AggregateAttr

		rowTester tester
		keyWalker keyWalker

		// Keep track of the registered groupIndex so we can easier implement the iterator.
		// Simple index tracking won't be enough, because btree is self-balancing -- things change
		//
		// @todo worth considering using the inmembuffer and fork it around...
		groupIndex *btree.Generic[*aggregateGroup]
		groups     []*aggregateGroup
		i          int

		ctr int
	}

	keyWalker func(context.Context, ValueGetter, func(context.Context, groupKey, ValueGetter) error) error
)

// init initializes the execution step's state
func (xs *aggregate) init(ctx context.Context) (err error) {
	// Initialize state variables
	xs.scanRow = xs.initScanRow()
	xs.groupIndex = btree.NewGeneric[*aggregateGroup](xs.compareGroupKeys)
	xs.groups = make([]*aggregateGroup, 0, 128)
	xs.rowTester, err = prepareGenericRowTester(xs.filter)
	if err != nil {
		return
	}

	// Initialize the key maker
	kk := make([]*ql.ASTNode, 0, len(xs.groupDefs))
	for _, a := range xs.groupDefs {
		kk = append(kk, a.Expression)
	}

	xs.keyWalker, err = aggregateGroupKeyWalker(kk...)
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
		if xs.limitExceeded() {
			return false, nil
		}

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
		k, err := xs.keep(ctx, xs.scanRow)
		if err != nil {
			return false, err
		}
		if !k {
			continue
		}

		xs.ctr++
		break
	}

	return true, nil
}

func (xs *aggregate) limitExceeded() bool {
	return xs.filter.limit > 0 && xs.ctr >= int(xs.filter.limit)
}

func (xs *aggregate) More(limit uint, v ValueGetter) (err error) {
	// Redo the cursor
	xs.filter.cursor, err = filter.PagingCursorFrom(xs.filter.OrderBy(), v, xs.collectPrimaryAttributes()...)
	if err != nil {
		return
	}

	// Redo the row tester
	xs.rowTester, err = prepareGenericRowTester(xs.filter)
	if err != nil {
		return
	}

	// Redo the state
	// @todo adjust based on aggregation plan; reuse buffered, etc.
	xs.scanRow.Reset()
	xs.groupIndex = btree.NewGeneric[*aggregateGroup](xs.compareGroupKeys)
	xs.groups = make([]*aggregateGroup, 0, 128)
	xs.planned = false
	xs.i = 0

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
	c, err := filter.PagingCursorFrom(s.filter.OrderBy(), v, s.collectPrimaryAttributes()...)
	if err != nil {
		return nil, err
	}

	c.ROrder = true
	c.LThen = s.filter.OrderBy().Reversed()

	return c, nil
}

func (s *aggregate) ForwardCursor(v ValueGetter) (*filter.PagingCursor, error) {
	c, err := filter.PagingCursorFrom(s.filter.OrderBy(), v, s.collectPrimaryAttributes()...)
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
		// @todo we probably can reuse the key or at least cache keys and avoid re-computation.
		//       My fairly hacky attempt boosted performance by ~20%
		err = xs.keyWalker(ctx, r, xs.addToGroup)
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

func (s *aggregate) wrapGroup(ctx context.Context, key groupKey) (g *aggregateGroup, err error) {
	agg := Aggregator()

	for _, a := range s.aggregateDefs {
		err = agg.AddAggregate(a.Identifier, a.Expression)
		if err != nil {
			return
		}
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

func (xs *aggregate) scanKey(g *aggregateGroup, dst *Row) (err error) {
	for i, attr := range xs.groupDefs {
		// @todo multi value support?
		// omitting err; internal row won't raise them
		dst.SetValue(attr.Identifier, 0, g.key[i])
	}

	return nil
}

func (s *aggregate) keep(ctx context.Context, r *Row) (bool, error) {
	if s.rowTester == nil {
		return true, nil
	}
	return s.rowTester.Test(ctx, r)
}

// Each group key is a PK candidate; all together form a composite key
func (s *aggregate) collectPrimaryAttributes() (out []string) {
	out = make([]string, 0, 2)
	for _, m := range s.def.Group {
		out = append(out, m.Identifier)
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

		for _, o := range s.filter.OrderBy() {
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

func inKeys(kk []AggregateAttr, ident string) int {
	for i, k := range kk {
		if k.Identifier == ident {
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
		out.values[attr.Identifier] = make([]any, 0, 2)
		out.counters[attr.Identifier] = 0
	}

	return
}

// aggregateGroupKeyWalker prepares a function which runs the provided function over all group keys
//
// This is required to handle multi-value attributes when used in group keys.
//
// Algorithm TL;DR
//
// We're going backwards in the slice of idents we need to handle (going the other way around)
// should also work but I chose to do it like so.
//
// For every ident, keep track of the number of elements and the index of the last unhandled ident.
//
// For every iteration, collect all of the values pointed to by the unused index.
// After the key is processed, increment the counter for the current ident.
// If the counter exceeds the number of elements, reset the counter for the current ident and move
// the ident pointer backwards (repeat if that ident's counter also exceeds the limit).
//
// When we find an indent which still has some items to process (counter doesn't exceed limit),
// reset the ident pointer to the end of the slice and repeat the whole thing.
func aggregateGroupKeyWalker(kk ...*ql.ASTNode) (out keyWalker, err error) {
	// @todo option to copy constants and idents
	runners, err := makeExprRunners(kk...)
	if err != nil {
		return
	}

	// We'll sort the idents to keep the output consistent; the order doesn't matter
	// but it will simplify testing.
	idents, hasConstants := keysFromExpr(kk...)
	sort.Strings(idents)

	out = func(ctx context.Context, vg ValueGetter, run func(context.Context, groupKey, ValueGetter) error) error {
		// Edgecase for when all of the expressions return constant values
		if len(idents) == 0 {
			// This should be impossible but better safe then sorry
			if !hasConstants {
				return nil
			}

			k, err := makeGroupKey(ctx, runners, vg)
			if err != nil {
				return err
			}

			return run(ctx, k, vg)
		}

		// For every value combination of idents used in agg. key expressions
		// construct a key and run the runner.

		limits := vg.CountValues()
		ptr := len(idents) - 1
		counts := make(map[string]uint, len(limits))
		for k := range limits {
			counts[k] = 0
		}

		handle := func() (err error) {
			aux := make(map[string]any, len(limits))

			for k, i := range counts {
				aux[k], err = vg.GetValue(k, i)
				if err != nil {
					return
				}
			}

			kk, err := makeGroupKey(ctx, runners, aux)
			if err != nil {
				return err
			}

			return run(ctx, kk, vg)
		}

	outer:
		for ptr >= 0 {
			handle()

			counts[idents[ptr]]++
			for {
				// There are still values to process so we can skip the rest
				if counts[idents[ptr]] < limits[idents[ptr]] {
					continue outer
				}

				// Reset the counter for the current ptr since we'll move back.
				// The outer loop will have this reset all of the counters after the current ptr.
				counts[idents[ptr]] = 0

				ptr--
				if ptr < 0 {
					break outer
				}

				counts[idents[ptr]]++

				if counts[idents[ptr]] >= limits[idents[ptr]] {
					continue
				}

				// We need to reset to the end so the next ident gets all of the values
				// that appear after it.
				ptr = len(idents) - 1
				break
			}
		}

		return nil
	}

	return
}

func makeExprRunners(kk ...*ql.ASTNode) (out []*runnerGval, err error) {
	out = make([]*runnerGval, len(kk))

	for i, k := range kk {
		out[i], err = newRunnerGvalParsed(k)
		if err != nil {
			return
		}
	}

	return
}

func makeGroupKey(ctx context.Context, runners []*runnerGval, vals any) (gk groupKey, err error) {
	gk = make(groupKey, len(runners))
	for i, r := range runners {
		v, err := r.Eval(ctx, vals)
		if err != nil {
			return nil, err
		}

		gk[i] = v
	}
	return gk, nil
}
