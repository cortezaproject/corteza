package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func testAttachment(t *testing.T, s store.Attachments) {
	var (
		ctx = context.Background()

		makeNew = func(nn ...string) *types.Attachment {
			// minimum data set for new attachment
			name := strings.Join(nn, "")
			return &types.Attachment{
				ID:        id.Next(),
				CreatedAt: *now(),
				Name:      "handle_" + name,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Attachment) {
			req := require.New(t)
			req.NoError(s.TruncateAttachments(ctx))
			res := makeNew()
			req.NoError(s.CreateAttachment(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		attachment := &types.Attachment{
			ID:        id.Next(),
			CreatedAt: *now(),
		}
		req.NoError(s.CreateAttachment(ctx, attachment))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, att := truncAndCreate(t)

		fetched, err := s.LookupAttachmentByID(ctx, att.ID)
		req.NoError(err)
		req.Equal(att.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, att := truncAndCreate(t)
		att.Url = "url"
		req.NoError(s.UpdateAttachment(ctx, att))
		fetched, err := s.LookupAttachmentByID(ctx, att.ID)
		req.NoError(err)
		req.Equal(att.ID, fetched.ID)
		req.Equal("url", fetched.Url)

	})

	t.Run("search", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("search by *", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
