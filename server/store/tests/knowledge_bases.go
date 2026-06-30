package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testKnowledgeBases(t *testing.T, s store.KnowledgeBases) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.KnowledgeBase {
			return &types.KnowledgeBase{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
				Title:     "Test KB",
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.KnowledgeBase) {
			req := require.New(t)
			req.NoError(s.TruncateKnowledgeBases(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateKnowledgeBase(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateKnowledgeBases(ctx))
		kb := makeNew("KnowledgeBaseCRUD")
		req.NoError(s.CreateKnowledgeBase(ctx, kb))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, kb := truncAndCreate(t)
		fetched, err := s.LookupKnowledgeBaseByID(ctx, kb.ID)
		req.NoError(err)
		req.Equal(kb.Handle, fetched.Handle)
		req.Equal(kb.ID, fetched.ID)
	})

	t.Run("lookup by handle", func(t *testing.T) {
		req, kb := truncAndCreate(t)
		fetched, err := s.LookupKnowledgeBaseByHandle(ctx, kb.Handle)
		req.NoError(err)
		req.Equal(kb.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, kb := truncAndCreate(t)
		kb.Title = "Updated"
		req.NoError(s.UpdateKnowledgeBase(ctx, kb))
		fetched, err := s.LookupKnowledgeBaseByID(ctx, kb.ID)
		req.NoError(err)
		req.Equal("Updated", fetched.Title)
	})

	t.Run("delete", func(t *testing.T) {
		req, kb := truncAndCreate(t)
		req.NoError(s.DeleteKnowledgeBaseByID(ctx, kb.ID))
	})
}
