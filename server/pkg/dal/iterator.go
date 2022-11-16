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

// PreLoadCursor into iterator and check it exist then return the cursor
// @todo this should be move to the Iterator
func PreLoadCursor(ctx context.Context, iter Iterator, limit uint, reverse bool, r ValueGetter) (out *filter.PagingCursor, err error) {
	if reverse {
		out, err = iter.BackCursor(r)
	} else {
		out, err = iter.ForwardCursor(r)
	}

	if err != nil {
		return
	}

	err = iter.(iterator).Preload(ctx, limit, out)
	if err != nil {
		return nil, nil
	}

	if !iter.Next(ctx) {
		out = nil
	}

	return
}

// IteratorSorting return iterator sorting
// @todo this should be move to the Iterator
func IteratorSorting(iter Iterator) filter.SortExprSet {
	return iter.(iterator).Sorting()
}
