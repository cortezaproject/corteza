package tests

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func testComposeAttachments(t *testing.T, s store.ComposeAttachments) {
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
			req.NoError(s.TruncateComposeAttachments(ctx))
			res := makeNew()
			req.NoError(s.CreateComposeAttachment(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		attachment := &types.Attachment{
			ID:        id.Next(),
			CreatedAt: *now(),
		}
		req.NoError(s.CreateComposeAttachment(ctx, attachment))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, att := truncAndCreate(t)

		fetched, err := s.LookupComposeAttachmentByID(ctx, att.ID)
		req.NoError(err)
		req.Equal(att.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, att := truncAndCreate(t)
		att.Url = "url"
		req.NoError(s.UpdateComposeAttachment(ctx, att))
		fetched, err := s.LookupComposeAttachmentByID(ctx, att.ID)
		req.NoError(err)
		req.Equal(att.ID, fetched.ID)
		req.Equal("url", fetched.Url)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Attachment", func(t *testing.T) {
			req, att := truncAndCreate(t)
			req.NoError(s.DeleteComposeAttachment(ctx, att))
			_, err := s.LookupComposeAttachmentByID(ctx, att.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, att := truncAndCreate(t)
			req.NoError(s.DeleteComposeAttachmentByID(ctx, att.ID))
			_, err := s.LookupComposeAttachmentByID(ctx, att.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
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
