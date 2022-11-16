package tests

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func testAttachment(t *testing.T, s store.Attachments) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(nn ...string) *types.Attachment {
			name := strings.Join(nn, "")
			return &types.Attachment{
				ID:        id.Next(),
				Name:      "handle_" + name,
				Kind:      "test-kind" + name,
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Attachment) {
			req := require.New(t)
			req.NoError(s.TruncateAttachments(ctx))
			attachment := makeNew()
			req.NoError(s.CreateAttachment(ctx, attachment))
			return req, attachment
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, types.AttachmentSet) {
			req := require.New(t)
			req.NoError(s.TruncateAttachments(ctx))

			set := make([]*types.Attachment, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateAttachment(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateAttachment(ctx, makeNew()))
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

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, att := truncAndCreate(t)
			att.Url = "url"

			req.NoError(s.UpsertAttachment(ctx, att))

			upserted, err := s.LookupAttachmentByID(ctx, att.ID)
			req.NoError(err)
			req.Equal(att.Name, upserted.Name)
		})

		t.Run("new", func(t *testing.T) {
			att := makeNew("upsert me", "upsert-me")

			req.NoError(s.UpsertAttachment(ctx, att))

			upserted, err := s.LookupAttachmentByID(ctx, att.ID)
			req.NoError(err)
			req.Equal(att.Name, upserted.Name)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Attachment", func(t *testing.T) {
			req, att := truncAndCreate(t)
			req.NoError(s.DeleteAttachment(ctx, att))
			_, err := s.LookupAttachmentByID(ctx, att.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, att := truncAndCreate(t)
			req.NoError(s.DeleteAttachmentByID(ctx, att.ID))
			_, err := s.LookupAttachmentByID(ctx, att.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by kind", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, f, err := s.SearchAttachments(ctx, types.AttachmentFilter{Kind: prefill[0].Kind})
			req.NoError(err)
			req.Equal(prefill[0].Kind, f.Kind)
			req.Len(set, 1)
		})

		t.Run("with check", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, _, err := s.SearchAttachments(ctx, types.AttachmentFilter{
				Check: func(attachment *types.Attachment) (bool, error) {
					return attachment.Kind == prefill[0].Kind, nil
				},
			})

			req.NoError(err)
			req.Len(set, 1)
			req.Equal(prefill[0].Kind, set[0].Kind)
		})
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
