package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testQueues(t *testing.T, s store.Queues) {
	var (
		ctx = context.Background()
		new = &types.Queue{
			ID:        42,
			Consumer:  "testConsumer",
			Queue:     "testQueue",
			Meta:      types.QueueMeta{},
			CreatedAt: time.Now(),
			CreatedBy: 1}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateQueues(ctx))
		req.NoError(s.CreateQueue(ctx, new))
	})

	t.Run("update", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateQueues(ctx))
		req.NoError(s.UpdateQueue(ctx, new))
	})

	t.Run("upsert", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateQueues(ctx))
		req.NoError(s.UpsertQueue(ctx, &types.Queue{ID: 42, Queue: "test", CreatedAt: *now()}))
		set, _, err := s.SearchQueues(ctx, types.QueueFilter{})
		req.NoError(err)
		req.Len(set, 1)
		req.True(set[0].Queue == "test")

		req.NoError(s.UpsertQueue(ctx, &types.Queue{ID: 42, Queue: "foobar", CreatedAt: *now()}))
		set, _, err = s.SearchQueues(ctx, types.QueueFilter{})
		req.NoError(err)
		req.Len(set, 1)
		req.True(set[0].Queue == "foobar")
	})

	t.Run("search", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateQueues(ctx))
		req.NoError(s.CreateQueue(ctx,
			new,
		))

		set, _, err := s.SearchQueues(ctx, types.QueueFilter{Query: "EST"})
		req.NoError(err)
		req.Len(set, 1)

		set, _, err = s.SearchQueues(ctx, types.QueueFilter{Query: "foo"})
		req.NoError(err)
		req.Len(set, 0)
	})
}
