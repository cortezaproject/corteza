package tests

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testAutomationWorkflows(t *testing.T, s store.AutomationWorkflows) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.Workflow {
			return &types.Workflow{
				ID:        id.Next(),
				Handle:    handle,
				CreatedAt: *now(),
				Enabled:   true,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Workflow) {
			req := require.New(t)
			req.NoError(s.TruncateAutomationWorkflows(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateAutomationWorkflow(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		wf := &types.Workflow{
			ID:        id.Next(),
			CreatedAt: *now(),
		}
		req.NoError(s.CreateAutomationWorkflow(ctx, wf))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, wf := truncAndCreate(t)

		fetched, err := s.LookupAutomationWorkflowByID(ctx, wf.ID)
		req.NoError(err)
		req.Equal(wf.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, wf := truncAndCreate(t)
		wf.Enabled = true
		req.NoError(s.UpdateAutomationWorkflow(ctx, wf))
		fetched, err := s.LookupAutomationWorkflowByID(ctx, wf.ID)
		req.NoError(err)
		req.Equal(wf.ID, fetched.ID)
		req.True(fetched.Enabled)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Workflow", func(t *testing.T) {
			req, wf := truncAndCreate(t)
			req.NoError(s.DeleteAutomationWorkflow(ctx, wf))
			_, err := s.LookupAutomationWorkflowByID(ctx, wf.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, wf := truncAndCreate(t)
			req.NoError(s.DeleteAutomationWorkflowByID(ctx, wf.ID))
			_, err := s.LookupAutomationWorkflowByID(ctx, wf.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Workflow{
			makeNew("one-one"),
			makeNew("one-two"),
			makeNew("two-one"),
			makeNew("two-two"),
			makeNew("two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateAutomationWorkflows(ctx))
		req.NoError(s.CreateAutomationWorkflow(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchAutomationWorkflows(ctx, types.WorkflowFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchAutomationWorkflows(ctx, types.WorkflowFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchAutomationWorkflows(ctx, types.WorkflowFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		// find all prefixed
		set, f, err = s.SearchAutomationWorkflows(ctx, types.WorkflowFilter{Query: "two-"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}
