package ctxcache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	var (
		req   = require.New(t)
		loads = 0
		load  = func() (int, error) { loads++; return loads, nil }
	)

	t.Run("without a cache every call loads", func(t *testing.T) {
		loads = 0
		ctx := context.Background()
		_, _ = Get(ctx, "k", load)
		v, err := Get(ctx, "k", load)
		req.NoError(err)
		req.Equal(2, v)
	})

	t.Run("with a cache a key loads once", func(t *testing.T) {
		loads = 0
		ctx := With(context.Background(), New(time.Minute))
		_, _ = Get(ctx, "k", load)
		v, err := Get(ctx, "k", load)
		req.NoError(err)
		req.Equal(1, v)

		v, err = Get(ctx, "other", load)
		req.NoError(err)
		req.Equal(2, v)
	})

	t.Run("expired entries load again", func(t *testing.T) {
		defer func() { now = time.Now }()

		loads = 0
		at := time.Now()
		now = func() time.Time { return at }

		ctx := With(context.Background(), New(time.Second))
		_, _ = Get(ctx, "k", load)

		at = at.Add(2 * time.Second)
		v, err := Get(ctx, "k", load)
		req.NoError(err)
		req.Equal(2, v)
	})

	t.Run("errors are not cached", func(t *testing.T) {
		var (
			calls = 0
			fail  = func() (int, error) { calls++; return 0, errors.New("failed") }
			ctx   = With(context.Background(), New(time.Minute))
		)

		_, err := Get(ctx, "k", fail)
		req.Error(err)
		_, err = Get(ctx, "k", fail)
		req.Error(err)
		req.Equal(2, calls)
	})
}
