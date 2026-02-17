package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func testLlmProviders(t *testing.T, s store.LlmProviders) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.LlmProvider {
			return &types.LlmProvider{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
				Status:    "active",
				Provider:  "openai",
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.LlmProvider) {
			req := require.New(t)
			req.NoError(s.TruncateLlmProviders(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateLlmProvider(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateLlmProviders(ctx))
		llmProvider := makeNew("LlmProviderCRUD")
		req.NoError(s.CreateLlmProvider(ctx, llmProvider))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, llmProvider := truncAndCreate(t)
		fetched, err := s.LookupLlmProviderByID(ctx, llmProvider.ID)
		req.NoError(err)
		req.Equal(llmProvider.Handle, fetched.Handle)
		req.Equal(llmProvider.ID, fetched.ID)
	})

	t.Run("lookup by handle", func(t *testing.T) {
		req, llmProvider := truncAndCreate(t)
		fetched, err := s.LookupLlmProviderByHandle(ctx, llmProvider.Handle)
		req.NoError(err)
		req.Equal(llmProvider.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, llmProvider := truncAndCreate(t)
		llmProvider.Status = "disabled"
		req.NoError(s.UpdateLlmProvider(ctx, llmProvider))
		fetched, err := s.LookupLlmProviderByID(ctx, llmProvider.ID)
		req.NoError(err)
		req.Equal("disabled", fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, llmProvider := truncAndCreate(t)
		req.NoError(s.DeleteLlmProviderByID(ctx, llmProvider.ID))
	})
}
