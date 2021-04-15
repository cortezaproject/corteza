package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testMessagebusQueueSettings(t *testing.T, s store.MessagebusQueueSettings) {
	var (
		ctx = context.Background()
		new = &messagebus.QueueSettings{
			ID:        42,
			Consumer:  "testConsumer",
			Queue:     "testQueue",
			Meta:      messagebus.QueueSettingsMeta{},
			CreatedAt: time.Now(),
			CreatedBy: 1}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateMessagebusQueueSettings(ctx))
		req.NoError(s.CreateMessagebusQueueSetting(ctx, new))
	})

	t.Run("update", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateMessagebusQueueSettings(ctx))
		req.NoError(s.UpdateMessagebusQueueSetting(ctx, new))
	})

	t.Run("upsert", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateMessagebusQueueSettings(ctx))
		req.NoError(s.UpsertMessagebusQueueSetting(ctx, &messagebus.QueueSettings{ID: 42, Queue: "test"}))
		set, _, err := s.SearchMessagebusQueueSettings(ctx, messagebus.QueueSettingsFilter{})
		req.NoError(err)
		req.Len(set, 1)
		req.True(set[0].Queue == "test")

		req.NoError(s.UpsertMessagebusQueueSetting(ctx, &messagebus.QueueSettings{ID: 42, Queue: "foobar"}))
		set, _, err = s.SearchMessagebusQueueSettings(ctx, messagebus.QueueSettingsFilter{})
		req.NoError(err)
		req.Len(set, 1)
		req.True(set[0].Queue == "foobar")
	})

	t.Run("search", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateMessagebusQueueSettings(ctx))
		req.NoError(s.CreateMessagebusQueueSetting(ctx,
			new,
		))

		set, _, err := s.SearchMessagebusQueueSettings(ctx, messagebus.QueueSettingsFilter{})
		req.NoError(err)
		req.Len(set, 1)
	})
}
