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

func testProjectIncidents(t *testing.T, s store.ProjectIncidents) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(projectID uint64) *types.ProjectIncident {
			return &types.ProjectIncident{
				ID:        id.Next(),
				ProjectID: projectID,
				Title:     "test incident",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectIncident) {
			req := require.New(t)
			req.NoError(s.TruncateProjectIncidents(ctx))
			res := makeNew(id.Next())
			req.NoError(s.CreateProjectIncident(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectIncidents(ctx))
		req.NoError(s.CreateProjectIncident(ctx, makeNew(id.Next())))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, res := truncAndCreate(t)
		fetched, err := s.LookupProjectIncidentByID(ctx, res.ID)
		req.NoError(err)
		req.Equal(res.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, res := truncAndCreate(t)
		res.Title = "updated"
		req.NoError(s.UpdateProjectIncident(ctx, res))
		fetched, err := s.LookupProjectIncidentByID(ctx, res.ID)
		req.NoError(err)
		req.Equal("updated", fetched.Title)
	})

	t.Run("delete", func(t *testing.T) {
		req, res := truncAndCreate(t)
		req.NoError(s.DeleteProjectIncidentByID(ctx, res.ID))
	})
}
