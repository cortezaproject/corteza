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
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		mma := &types.MessageAttachment{
			MessageID:    id.Next(),
			AttachmentID: id.Next(),
		}
		req.NoError(s.CreateMessagingMessageAttachment(ctx, mma))
	})

	t.Run("delete", func(t *testing.T) {
		req := require.New(t)
		mma := &types.MessageAttachment{
			MessageID:    id.Next(),
			AttachmentID: id.Next(),
		}
		req.NoError(s.DeleteMessagingMessageAttachment(ctx, mma))
	})

}
