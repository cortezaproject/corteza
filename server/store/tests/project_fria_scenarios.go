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

func testProjectFriaScenarios(t *testing.T, s store.ProjectFriaScenarios) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(projectID uint64) *types.ProjectFriaScenario {
			return &types.ProjectFriaScenario{
				ID:         id.Next(),
				ProjectID:  projectID,
				AiSystemID: id.Next(),
				Title:      "test scenario",
				Severity:   "high",
				CreatedAt:  time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectFriaScenario) {
			req := require.New(t)
			req.NoError(s.TruncateProjectFriaScenarios(ctx))
			res := makeNew(id.Next())
			req.NoError(s.CreateProjectFriaScenario(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectFriaScenarios(ctx))
		req.NoError(s.CreateProjectFriaScenario(ctx, makeNew(id.Next())))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, res := truncAndCreate(t)
		fetched, err := s.LookupProjectFriaScenarioByID(ctx, res.ID)
		req.NoError(err)
		req.Equal(res.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, res := truncAndCreate(t)
		res.Title = "updated"
		req.NoError(s.UpdateProjectFriaScenario(ctx, res))
		fetched, err := s.LookupProjectFriaScenarioByID(ctx, res.ID)
		req.NoError(err)
		req.Equal("updated", fetched.Title)
	})

	t.Run("delete", func(t *testing.T) {
		req, res := truncAndCreate(t)
		req.NoError(s.DeleteProjectFriaScenarioByID(ctx, res.ID))
	})
}
