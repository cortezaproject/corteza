package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
)

func testMessagingMessageAttachments(t *testing.T, s store.MessagingMessageAttachments) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func() *types.MessageAttachment {
			// minimum data set for new messageAttachment
			return &types.MessageAttachment{
				MessageID:    id.Next(),
				AttachmentID: id.Next(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.MessageAttachment) {
			req := require.New(t)
			req.NoError(s.TruncateMessagingMessageAttachments(ctx))
			res := makeNew()
			req.NoError(s.CreateMessagingMessageAttachment(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		mma := &types.MessageAttachment{
			MessageID:    id.Next(),
			AttachmentID: id.Next(),
		}
		req.NoError(s.CreateMessagingMessageAttachment(ctx, mma))
	})

	t.Run("update", func(t *testing.T) {
		req, att := truncAndCreate(t)
		att.AttachmentID = id.Next()
		req.NoError(s.UpdateMessagingMessageAttachment(ctx, att))
		fetched, err := s.LookupMessagingMessageAttachmentByMessageID(ctx, att.MessageID)
		req.NoError(err)
		req.Equal(att.AttachmentID, fetched.AttachmentID)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, att := truncAndCreate(t)
			att.AttachmentID = id.Next()
			req.NoError(s.UpsertMessagingMessageAttachment(ctx, att))
			fetched, err := s.LookupMessagingMessageAttachmentByMessageID(ctx, att.MessageID)
			req.NoError(err)
			req.Equal(att.AttachmentID, fetched.AttachmentID)
		})

		t.Run("new", func(t *testing.T) {
			att := makeNew()
			att.AttachmentID = id.Next()
			req.NoError(s.UpsertMessagingMessageAttachment(ctx, att))

			upserted, err := s.LookupMessagingMessageAttachmentByMessageID(ctx, att.MessageID)
			req.NoError(err)
			req.Equal(att.AttachmentID, upserted.AttachmentID)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by MessageAttachment", func(t *testing.T) {
			req, mma := truncAndCreate(t)
			req.NoError(s.DeleteMessagingMessageAttachment(ctx, mma))
			_, err := s.LookupMessagingMessageAttachmentByMessageID(ctx, mma.MessageID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by MessageID", func(t *testing.T) {
			req, mma := truncAndCreate(t)
			req.NoError(s.DeleteMessagingMessageAttachmentByMessageID(ctx, mma.MessageID))
			_, err := s.LookupMessagingMessageAttachmentByMessageID(ctx, mma.MessageID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})
}
