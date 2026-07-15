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

func testProjectPrivacys(t *testing.T, s store.ProjectPrivacys) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(projectID uint64) *types.ProjectPrivacy {
			return &types.ProjectPrivacy{
				ID:        id.Next(),
				ProjectID: projectID,
				Title:     "test privacy",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectPrivacy) {
			req := require.New(t)
			req.NoError(s.TruncateProjectPrivacys(ctx))
			res := makeNew(id.Next())
			req.NoError(s.CreateProjectPrivacy(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectPrivacys(ctx))
		req.NoError(s.CreateProjectPrivacy(ctx, makeNew(id.Next())))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, res := truncAndCreate(t)
		fetched, err := s.LookupProjectPrivacyByID(ctx, res.ID)
		req.NoError(err)
		req.Equal(res.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, res := truncAndCreate(t)
		res.Title = "updated"
		req.NoError(s.UpdateProjectPrivacy(ctx, res))
		fetched, err := s.LookupProjectPrivacyByID(ctx, res.ID)
		req.NoError(err)
		req.Equal("updated", fetched.Title)
	})

	t.Run("delete", func(t *testing.T) {
		req, res := truncAndCreate(t)
		req.NoError(s.DeleteProjectPrivacyByID(ctx, res.ID))
	})
}
