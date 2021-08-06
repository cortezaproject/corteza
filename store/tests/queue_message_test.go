package tests

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testQueueMessage(t *testing.T, s store.QueueMessages) {
	var (
		ctx           = context.Background()
		foobarMessage = &types.QueueMessage{Payload: []byte(`foobar`), Created: now()}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateQueueMessages(ctx))
		req.NoError(s.CreateQueueMessage(ctx, foobarMessage))
	})

	t.Run("update", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateQueueMessages(ctx))
		req.NoError(s.UpdateQueueMessage(ctx, foobarMessage))
	})

	t.Run("search", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateQueueMessage(ctx,
			&types.QueueMessage{ID: 1, Queue: "test", Payload: []byte{}, Created: now()},
			&types.QueueMessage{ID: 2, Queue: "test", Payload: []byte{}, Created: now()},
			&types.QueueMessage{ID: 3, Queue: "test", Payload: []byte{}, Created: now()},
		))

		set, _, err := s.SearchQueueMessages(ctx, types.QueueMessageFilter{})
		req.NoError(err)
		req.Len(set, 3)
	})
}
