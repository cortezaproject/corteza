package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func testMessagingUnread(t *testing.T, s store.MessagingUnreads) {
	var (
		ctx = context.Background()
		req = require.New(t)

		channelID = id.Next()
		userID    = id.Next()
		threadID  = id.Next()

		makeNew = func(channelID, threadID, userID uint64) *types.Unread {
			// minimum data set for new messagingUnread
			return &types.Unread{
				ChannelID: channelID,
				ReplyTo:   threadID,
				UserID:    userID,
			}
		}
	)

	t.Run("create", func(t *testing.T) {
		messagingUnread := makeNew(channelID, threadID, userID)
		req.NoError(s.CreateMessagingUnread(ctx, messagingUnread))
	})
}
