package dal

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// OrderedBuffer provides the buffered data in the specified order
	OrderedBuffer interface {
		Buffer

		// InOrder makes the buffer provide the stored data in the specified order
		InOrder(ss ...*filter.SortExpr) (err error)
	}

	// Buffer provides a place where you can buffer the data provided by DAL
	// @note Buffers are currently primarily used for testing
	Buffer interface {
		// Seek moves the index pointer to the specified location
		// After the Seek call, a Next() call is required
		Seek(context.Context, int) error

		// Len returns the number of elements in the buffer
		Len() int

		Iterator
		Adder
	}

	Adder interface {
		// Add adds a new ValueGetter to the buffer
		Add(context.Context, ValueGetter) (err error)
	}
)
