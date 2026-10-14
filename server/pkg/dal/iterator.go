package dal

import (
	"context"
	"encoding/json"
	"io"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	// Iterator provides an interface for loading data from the underlying store
	Iterator interface {
		Next(ctx context.Context) bool
		More(uint, ValueGetter) error
		Err() error
		Scan(ValueSetter) error
		Close() error

		BackCursor(ValueGetter) (*filter.PagingCursor, error)
		ForwardCursor(ValueGetter) (*filter.PagingCursor, error)

		// // -1 means unknown
		// Total() int
		// Cursor() any
		// // ... do we need anything else here?
	}

	iterator interface {
		Preload(context.Context, uint, *filter.PagingCursor) error
		Sorting() filter.SortExprSet
	}
)

// IteratorEncodeJSON helper function that encodes each item from the iterator as JSON
// and writes it to th given io.Writer.
//
// target initialization function is intentionally used to avoid use of reflection
func IteratorEncodeJSON(ctx context.Context, w io.Writer, iter Iterator, initTarget func() ValueSetter) (err error) {
	var (
		target   ValueSetter
		firstOut = false
	)

	for iter.Next(ctx) {
		if err = iter.Err(); err != nil {
			return
		}

		if firstOut {
			if _, err = w.Write([]byte(`,`)); err != nil {
				return
			}
		}

		firstOut = true

		target = initTarget()

		if err = iter.Scan(target); err != nil {
			return
		}

		err = json.NewEncoder(w).Encode(target)
		if err != nil {
			return
		}
	}

	return
}

// PreLoadCursor returns the cursor past r when at least one row beyond it
// passes fx, reading on in windows of limit rows until one does or rows run out
// @todo this should be move to the Iterator
func PreLoadCursor(ctx context.Context, iter Iterator, limit uint, reverse bool, r ValueGetter, fx func(Iterator) (ValueGetter, bool, error)) (out *filter.PagingCursor, err error) {
	makeCursor := func() (*filter.PagingCursor, error) {
		if reverse {
			return iter.BackCursor(r)
		} else {
			return iter.ForwardCursor(r)
		}
	}

	out, err = makeCursor()
	if err != nil {
		return
	}

	err = iter.(iterator).Preload(ctx, limit, out)
	if err != nil {
		return nil, nil
	}

	for {
		var (
			scanned uint

			// row the next window continues from; windows come in sort order,
			// so a reverse one continues before its first row
			edge ValueGetter
		)

		for iter.Next(ctx) {
			row, ok, err := fx(iter)
			if err != nil {
				return nil, err
			}

			if ok {
				return out, nil
			}

			scanned++
			if edge == nil || !reverse {
				edge = row
			}
		}

		if err = iter.Err(); err != nil {
			return nil, err
		}

		if scanned < limit {
			return nil, nil
		}

		if err = iter.More(limit, edge); err != nil {
			return nil, err
		}
	}
}

// IteratorSorting return iterator sorting
// @todo this should be move to the Iterator
func IteratorSorting(iter Iterator) filter.SortExprSet {
	return iter.(iterator).Sorting()
}
