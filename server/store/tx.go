package store

import (
	"context"
)

type (
	postCommitKey struct{}

	// postCommit holds work that must not run until the transaction it was
	// registered inside has committed.
	postCommit struct {
		fns []func()
	}
)

// AfterCommit schedules fn to run once the transaction ctx belongs to commits.
//
// Outside a transaction there is nothing to wait for, so fn runs right away.
//
// Anything that reads the change back — a cache refresh, an event a handler
// re-reads the store for — belongs here: inside the transaction the row is not
// yet visible to any other connection.
func AfterCommit(ctx context.Context, fn func()) {
	if q, ok := ctx.Value(postCommitKey{}).(*postCommit); ok {
		q.fns = append(q.fns, fn)
		return
	}

	fn()
}

func Tx(ctx context.Context, s Storer, fn func(context.Context, Storer) error) error {
	if _, nested := ctx.Value(postCommitKey{}).(*postCommit); nested {
		// the outermost transaction owns the queue and the commit that drains it
		return s.Tx(ctx, fn)
	}

	q := &postCommit{}
	ctx = context.WithValue(ctx, postCommitKey{}, q)

	err := s.Tx(ctx, func(ctx context.Context, s Storer) error {
		// a retried attempt starts over with an empty queue
		q.fns = nil
		return fn(ctx, s)
	})

	if err != nil {
		return err
	}

	for _, fn := range q.fns {
		fn()
	}

	return nil
}
