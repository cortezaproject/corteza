package tests

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testMessagebusQueueMessage(t *testing.T, s store.MessagebusQueueMessages) {
	var (
		ctx           = context.Background()
		foobarMessage = &messagebus.QueueMessage{Payload: []byte(`foobar`), Created: now()}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateMessagebusQueueMessages(ctx))
		req.NoError(s.CreateMessagebusQueueMessage(ctx, foobarMessage))
	})

	t.Run("update", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateMessagebusQueueMessages(ctx))
		req.NoError(s.UpdateMessagebusQueueMessage(ctx, foobarMessage))
	})

	t.Run("search", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateMessagebusQueueMessage(ctx,
			&messagebus.QueueMessage{ID: 1, Queue: "test", Payload: []byte{}, Created: now()},
			&messagebus.QueueMessage{ID: 2, Queue: "test", Payload: []byte{}, Created: now()},
			&messagebus.QueueMessage{ID: 3, Queue: "test", Payload: []byte{}, Created: now()},
		))

		set, _, err := s.SearchMessagebusQueueMessages(ctx, messagebus.QueueMessageFilter{})
		req.NoError(err)
		req.Len(set, 3)
	})
}
