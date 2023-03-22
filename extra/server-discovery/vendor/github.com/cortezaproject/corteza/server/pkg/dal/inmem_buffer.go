package dal

import (
	"context"
	"fmt"
	"sort"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	// inmemBuffer isa simple DAL buffer which holds all of the data in memory
	//
	// This buffer should be used for small datasets
	inmemBuffer struct {
		rows []ValueGetter
		ctrs []map[string]uint

		i      int
		cap    int
		closed bool

		more bool
		err  error

		sort filter.SortExprSet
	}
)

// InMemoryBuffer initializes a new DAL buffer where the data is kept in memory
func InMemoryBuffer() *inmemBuffer {
	return &inmemBuffer{
		// @note we'll buffer the value counters along side the actual value counter
		// because constant hashmap initialization is quite memory intensive and introduces
		// a bit of a bottleneck.
		//
		// @todo investigate if this can be optimized at a lower level
		rows: make([]ValueGetter, 0, 100),
		ctrs: make([]map[string]uint, 0, 100),

		// i starts off as -1 so we don't need an extra state flag when doing the first Next
		i:    -1,
		more: true,
	}
}

// InMemoryBufferWith returns a new buffer with the given value getters
func InMemoryBufferWith(ctx context.Context, vv ...ValueGetter) (Buffer, error) {
	b := InMemoryBuffer()
	for _, v := range vv {
		if err := b.Add(ctx, v); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (b *inmemBuffer) InOrder(ss ...*filter.SortExpr) (err error) {
	b.sort = ss

	return nil
}

// Single makes the buffer only keep one element at the time
// @todo make this option convert the buffer into a circular buffer
func (b *inmemBuffer) Single() {
	b.cap = 1
	b.rows = make([]ValueGetter, 1)
	b.ctrs = make([]map[string]uint, 1)
	b.i = -1
}

func (b *inmemBuffer) Seek(_ context.Context, i int) (err error) {
	// we go one level further and start at -1
	b.i = i - 1

	b.more = i < b.Len()
	return
}

func (b *inmemBuffer) isSingle() bool {
	return b.cap == 1
}

func (b *inmemBuffer) isSorted() bool {
	return len(b.sort) > 0
}

func (b *inmemBuffer) Add(ctx context.Context, v ValueGetter) (err error) {
	if b.i != -1 && b.isSorted() {
		return fmt.Errorf("cannot buffer items after an access occurred and the buffer is ordered")
	}

	if b.isSingle() {
		b.rows[0] = v
		b.ctrs[0] = v.CountValues()
	} else {
		b.rows = append(b.rows, v)
		b.ctrs = append(b.ctrs, v.CountValues())
	}

	b.more = true
	return
}

func (b *inmemBuffer) Next(ctx context.Context) bool {
	if !b.more {
		return false
	}
	if b.closed {
		panic("cannot call Next on closed buffer")
	}

	// @todo use something like a b-tree to sort while inserting
	if b.i == -1 && len(b.sort) > 0 {
		sort.Sort(b)
	}

	b.i++
	if b.isSingle() {
		b.i = b.i % b.cap
	}

	b.more = !b.isSingle() && b.i+1 < b.Len()

	return b.i < b.Len()
}

func (b *inmemBuffer) More(uint, ValueGetter) error {
	return nil
}

func (b *inmemBuffer) Err() error {
	return b.err
}

func (b *inmemBuffer) Scan(s ValueSetter) (err error) {
	if b.closed {
		panic("cannot call Scan on closed buffer")
	}

	r, ctrs, err := b.rowAt(b.i)
	if err != nil {
		return
	}

	var v any
	for name, count := range ctrs {
		for i := uint(0); i < count; i++ {
			v, err = r.GetValue(name, i)
			if err != nil {
				return
			}
			err = s.SetValue(name, i, v)
			if err != nil {
				return
			}
		}
	}
	return
}

// @todo should this get a context for things like IO ops
func (b *inmemBuffer) Close() error {
	b.ctrs = nil
	b.rows = nil
	b.closed = true

	return nil
}

func (b *inmemBuffer) BackCursor(ValueGetter) (*filter.PagingCursor, error) {
	return nil, fmt.Errorf("not supported")
}

func (b *inmemBuffer) ForwardCursor(ValueGetter) (*filter.PagingCursor, error) {
	return nil, fmt.Errorf("not supported")
}

// rowAt returns the row at a given index
func (b *inmemBuffer) rowAt(i int) (ValueGetter, map[string]uint, error) {
	if i >= b.Len() {
		return nil, nil, nil
	}

	return b.rows[i], b.ctrs[i], nil
}

// sort.Interface methods

func (b *inmemBuffer) Len() int {
	return len(b.rows)
}

func (b *inmemBuffer) Less(i, j int) bool {
	ra, _, err := b.rowAt(i)
	if err != nil {
		panic(err)
	}

	rb, _, err := b.rowAt(j)
	if err != nil {
		panic(err)
	}

	return makeRowComparator(b.sort...)(ra, rb)
}

func (b *inmemBuffer) Swap(i, j int) {
	b.rows[i], b.rows[j] = b.rows[j], b.rows[i]
	b.ctrs[i], b.ctrs[j] = b.ctrs[j], b.ctrs[i]
}
