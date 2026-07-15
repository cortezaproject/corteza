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

func testProjectReviews(t *testing.T, s store.ProjectReviews) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(projectID uint64) *types.ProjectReview {
			return &types.ProjectReview{
				ID:        id.Next(),
				ProjectID: projectID,
				Title:     "test review",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectReview) {
			req := require.New(t)
			req.NoError(s.TruncateProjectReviews(ctx))
			res := makeNew(id.Next())
			req.NoError(s.CreateProjectReview(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectReviews(ctx))
		req.NoError(s.CreateProjectReview(ctx, makeNew(id.Next())))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, res := truncAndCreate(t)
		fetched, err := s.LookupProjectReviewByID(ctx, res.ID)
		req.NoError(err)
		req.Equal(res.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, res := truncAndCreate(t)
		res.Title = "updated"
		req.NoError(s.UpdateProjectReview(ctx, res))
		fetched, err := s.LookupProjectReviewByID(ctx, res.ID)
		req.NoError(err)
		req.Equal("updated", fetched.Title)
	})

	t.Run("delete", func(t *testing.T) {
		req, res := truncAndCreate(t)
		req.NoError(s.DeleteProjectReviewByID(ctx, res.ID))
	})
}
