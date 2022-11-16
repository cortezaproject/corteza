package tests

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testAutomationTriggers(t *testing.T, s store.AutomationTriggers) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func() *types.Trigger {
			return &types.Trigger{
				ID:        id.Next(),
				CreatedAt: *now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Trigger) {
			req := require.New(t)
			req.NoError(s.TruncateAutomationTriggers(ctx))
			res := makeNew()
			req.NoError(s.CreateAutomationTrigger(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		tr := &types.Trigger{
			ID:        id.Next(),
			CreatedAt: *now(),
		}
		req.NoError(s.CreateAutomationTrigger(ctx, tr))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, tr := truncAndCreate(t)

		fetched, err := s.LookupAutomationTriggerByID(ctx, tr.ID)
		req.NoError(err)
		req.Equal(tr.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, tr := truncAndCreate(t)
		tr.Enabled = true
		req.NoError(s.UpdateAutomationTrigger(ctx, tr))
		fetched, err := s.LookupAutomationTriggerByID(ctx, tr.ID)
		req.NoError(err)
		req.Equal(tr.ID, fetched.ID)
		req.True(fetched.Enabled)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Trigger", func(t *testing.T) {
			req, tr := truncAndCreate(t)
			req.NoError(s.DeleteAutomationTrigger(ctx, tr))
			_, err := s.LookupAutomationTriggerByID(ctx, tr.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, tr := truncAndCreate(t)
			req.NoError(s.DeleteAutomationTriggerByID(ctx, tr.ID))
			_, err := s.LookupAutomationTriggerByID(ctx, tr.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		auxWf := id.Next()

		// Create test triggers
		{
			req.NoError(s.CreateAutomationTrigger(ctx, &types.Trigger{
				ID:        id.Next(),
				CreatedAt: *now(),
				Enabled:   true,

				WorkflowID: auxWf,
			}))
			req.NoError(s.CreateAutomationTrigger(ctx, &types.Trigger{
				ID:        id.Next(),
				CreatedAt: *now(),
				Enabled:   true,

				EventType: "event.test",
			}))
			req.NoError(s.CreateAutomationTrigger(ctx, &types.Trigger{
				ID:        id.Next(),
				CreatedAt: *now(),
				Enabled:   true,

				ResourceType: "resource:test",
			}))
			req.NoError(s.CreateAutomationTrigger(ctx, &types.Trigger{
				ID:        id.Next(),
				CreatedAt: *now(),
				Enabled:   true,

				DeletedAt: now(),
			}))
			req.NoError(s.CreateAutomationTrigger(ctx, &types.Trigger{
				ID:        id.Next(),
				CreatedAt: *now(),

				Enabled: false,
			}))
		}

		tt, _, err := store.SearchAutomationTriggers(ctx, s, types.TriggerFilter{
			WorkflowID: []uint64{auxWf},
		})
		req.NoError(err)
		req.Len(tt, 1)
		tt, _, err = store.SearchAutomationTriggers(ctx, s, types.TriggerFilter{
			EventType: "event.test",
		})
		req.NoError(err)
		req.Len(tt, 1)
		tt, _, err = store.SearchAutomationTriggers(ctx, s, types.TriggerFilter{
			ResourceType: "resource:test",
		})
		req.NoError(err)
		req.Len(tt, 1)
		tt, _, err = store.SearchAutomationTriggers(ctx, s, types.TriggerFilter{
			Deleted: filter.StateExclusive,
		})
		req.NoError(err)
		req.Len(tt, 1)
		tt, _, err = store.SearchAutomationTriggers(ctx, s, types.TriggerFilter{
			Disabled: filter.StateExclusive,
		})
		req.NoError(err)
		req.Len(tt, 1)
	})
}
