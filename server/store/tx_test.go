package store

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	// only Tx is exercised; every other Storer method panics if reached
	txStore struct {
		Storer
		tx func(context.Context, func(context.Context, Storer) error) error
	}
)

func (s txStore) Tx(ctx context.Context, fn func(context.Context, Storer) error) error {
	return s.tx(ctx, fn)
}

// runs the body once, the way a driver with transactions disabled does
func passthroughTx() txStore {
	var s txStore
	s.tx = func(ctx context.Context, fn func(context.Context, Storer) error) error {
		return fn(ctx, s)
	}
	return s
}

func TestAfterCommit_OutsideTransaction(t *testing.T) {
	ran := false
	AfterCommit(context.Background(), func() { ran = true })
	require.True(t, ran, "nothing to wait for, so it should have run right away")
}

func TestAfterCommit_RunsOnlyOnceTheBodyIsDone(t *testing.T) {
	var order []string

	err := Tx(context.Background(), passthroughTx(), func(ctx context.Context, _ Storer) error {
		AfterCommit(ctx, func() { order = append(order, "after") })
		order = append(order, "write")
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, []string{"write", "after"}, order)
}

func TestAfterCommit_SkippedWhenTheTransactionFails(t *testing.T) {
	boom := errors.New("boom")
	ran := false

	err := Tx(context.Background(), passthroughTx(), func(ctx context.Context, _ Storer) error {
		AfterCommit(ctx, func() { ran = true })
		return boom
	})

	require.ErrorIs(t, err, boom)
	require.False(t, ran, "a rolled back change must not announce itself")
}

func TestAfterCommit_RetriedAttemptRunsItOnce(t *testing.T) {
	var (
		s        txStore
		attempts int
		runs     int
	)

	s.tx = func(ctx context.Context, fn func(context.Context, Storer) error) error {
		if err := fn(ctx, s); err != nil {
			return fn(ctx, s)
		}
		return nil
	}

	err := Tx(context.Background(), s, func(ctx context.Context, _ Storer) error {
		AfterCommit(ctx, func() { runs++ })
		attempts++
		if attempts == 1 {
			return errors.New("retry me")
		}
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, attempts)
	require.Equal(t, 1, runs, "the abandoned attempt's work must not be queued")
}

func TestAfterCommit_NestedTransactionDefersToTheOutermost(t *testing.T) {
	var order []string
	s := passthroughTx()

	err := Tx(context.Background(), s, func(ctx context.Context, _ Storer) error {
		if err := Tx(ctx, s, func(ctx context.Context, _ Storer) error {
			AfterCommit(ctx, func() { order = append(order, "after") })
			return nil
		}); err != nil {
			return err
		}

		order = append(order, "outer still open")
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, []string{"outer still open", "after"}, order)
}
