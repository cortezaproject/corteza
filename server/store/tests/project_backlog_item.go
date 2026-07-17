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

func testProjectBacklogItems(t *testing.T, s store.ProjectBacklogItems) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(projectID uint64) *types.ProjectBacklogItem {
			return &types.ProjectBacklogItem{
				ID:        id.Next(),
				ProjectID: projectID,
				EventID:   id.Next(),
				Category:  "task",
				Title:     "test backlog item",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectBacklogItem) {
			req := require.New(t)
			req.NoError(s.TruncateProjectBacklogItems(ctx))
			res := makeNew(id.Next())
			req.NoError(s.CreateProjectBacklogItem(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectBacklogItems(ctx))
		req.NoError(s.CreateProjectBacklogItem(ctx, makeNew(id.Next())))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, res := truncAndCreate(t)
		fetched, err := s.LookupProjectBacklogItemByID(ctx, res.ID)
		req.NoError(err)
		req.Equal(res.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, res := truncAndCreate(t)
		res.Title = "updated"
		req.NoError(s.UpdateProjectBacklogItem(ctx, res))
		fetched, err := s.LookupProjectBacklogItemByID(ctx, res.ID)
		req.NoError(err)
		req.Equal("updated", fetched.Title)
	})

	t.Run("delete", func(t *testing.T) {
		req, res := truncAndCreate(t)
		req.NoError(s.DeleteProjectBacklogItemByID(ctx, res.ID))
	})
}
