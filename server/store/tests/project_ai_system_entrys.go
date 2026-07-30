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

func testProjectGroupEntrys(t *testing.T, s store.ProjectGroupEntrys) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(groupID uint64, ref string) *types.ProjectGroupEntry {
			return &types.ProjectGroupEntry{
				ID:             id.Next(),
				ProjectGroupID: groupID,
				ResourceRef:    ref,
				CreatedAt:      time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectGroupEntry) {
			req := require.New(t)
			req.NoError(s.TruncateProjectGroupEntrys(ctx))
			res := makeNew(id.Next(), "corteza::system/user:1")
			req.NoError(s.CreateProjectGroupEntry(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectGroupEntrys(ctx))
		entry := makeNew(id.Next(), "corteza::system/user:1")
		req.NoError(s.CreateProjectGroupEntry(ctx, entry))
	})

	t.Run("lookup by group ID and resource ref", func(t *testing.T) {
		req, entry := truncAndCreate(t)
		fetched, err := s.LookupProjectGroupEntryByProjectGroupIDResourceRef(ctx, entry.ProjectGroupID, entry.ResourceRef)
		req.NoError(err)
		req.Equal(entry.ProjectGroupID, fetched.ProjectGroupID)
		req.Equal(entry.ResourceRef, fetched.ResourceRef)
	})

	t.Run("delete by group ID and resource ref", func(t *testing.T) {
		req, entry := truncAndCreate(t)
		req.NoError(s.DeleteProjectGroupEntryByProjectGroupIDResourceRef(ctx, entry.ProjectGroupID, entry.ResourceRef))
	})
}
