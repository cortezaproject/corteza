package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testProjectAiSystemEntrys(t *testing.T, s store.ProjectAiSystemEntrys) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(aiSystemID uint64, ref string) *types.ProjectAiSystemEntry {
			return &types.ProjectAiSystemEntry{
				ID:                id.Next(),
				ProjectAiSystemID: aiSystemID,
				ResourceRef:       ref,
				CreatedAt:         time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectAiSystemEntry) {
			req := require.New(t)
			req.NoError(s.TruncateProjectAiSystemEntrys(ctx))
			res := makeNew(id.Next(), "corteza::system/user:1")
			req.NoError(s.CreateProjectAiSystemEntry(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectAiSystemEntrys(ctx))
		entry := makeNew(id.Next(), "corteza::system/user:1")
		req.NoError(s.CreateProjectAiSystemEntry(ctx, entry))
	})

	t.Run("lookup by AI system ID and resource ref", func(t *testing.T) {
		req, entry := truncAndCreate(t)
		fetched, err := s.LookupProjectAiSystemEntryByProjectAiSystemIDResourceRef(ctx, entry.ProjectAiSystemID, entry.ResourceRef)
		req.NoError(err)
		req.Equal(entry.ProjectAiSystemID, fetched.ProjectAiSystemID)
		req.Equal(entry.ResourceRef, fetched.ResourceRef)
	})

	t.Run("delete by AI system ID and resource ref", func(t *testing.T) {
		req, entry := truncAndCreate(t)
		req.NoError(s.DeleteProjectAiSystemEntryByProjectAiSystemIDResourceRef(ctx, entry.ProjectAiSystemID, entry.ResourceRef))
	})
}
