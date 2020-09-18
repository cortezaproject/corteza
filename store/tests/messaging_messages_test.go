package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func testMessagingMessages(t *testing.T, s store.MessagingMessages) {
	var (
		ctx = context.Background()
		req = require.New(t)

		channelID = id.Next()

		makeNew = func(msg string) *types.Message {
			// minimum data set for new messagingMessage
			return &types.Message{
				ID:        id.Next(),
				ChannelID: channelID,
				CreatedAt: *now(),
				Message:   msg,
				Type:      types.MessageTypeSimpleMessage,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Message) {
			req := require.New(t)
			req.NoError(s.TruncateMessagingMessages(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateMessagingMessage(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		messagingMessage := makeNew("MessagingMessageCRUD")
		req.NoError(s.CreateMessagingMessage(ctx, messagingMessage))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		messagingMessage := makeNew("look up by id")
		req.NoError(s.CreateMessagingMessage(ctx, messagingMessage))
		fetched, err := s.LookupMessagingMessageByID(ctx, messagingMessage.ID)
		req.NoError(err)
		req.Equal(messagingMessage.Message, fetched.Message)
		req.Equal(messagingMessage.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		messagingMessage := makeNew("update me")
		req.NoError(s.CreateMessagingMessage(ctx, messagingMessage))

		messagingMessage = &types.Message{
			ID:        messagingMessage.ID,
			CreatedAt: messagingMessage.CreatedAt,
			Message:   "MessagingMessageCRUD+2",
			Type:      types.MessageTypeSimpleMessage,
		}
		req.NoError(s.UpdateMessagingMessage(ctx, messagingMessage))

		updated, err := s.LookupMessagingMessageByID(ctx, messagingMessage.ID)
		req.NoError(err)
		req.Equal(messagingMessage.Message, updated.Message)
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Message", func(t *testing.T) {
			req, messagingMessage := truncAndCreate(t)
			req.NoError(s.DeleteMessagingMessage(ctx, messagingMessage))
			_, err := s.LookupMessagingMessageByID(ctx, messagingMessage.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, messagingMessage := truncAndCreate(t)
			req.NoError(s.DeleteMessagingMessageByID(ctx, messagingMessage.ID))
			_, err := s.LookupMessagingMessageByID(ctx, messagingMessage.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Message{
			makeNew("/one-one"),
			makeNew("/one-two"),
			makeNew("/two-one"),
			makeNew("/two-two"),
			makeNew("/two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateMessagingMessages(ctx))
		req.NoError(s.CreateMessagingMessage(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchMessagingMessages(ctx, types.MessageFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// find all prefixed
		set, f, err = s.SearchMessagingMessages(ctx, types.MessageFilter{Query: "/two-"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}
