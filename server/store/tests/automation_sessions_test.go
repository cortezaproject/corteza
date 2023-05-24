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

func testAutomationSessions(t *testing.T, s store.AutomationSessions) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(wfID uint64, completed bool) *types.Session {
			ses := &types.Session{
				ID:         id.Next(),
				WorkflowID: wfID,
				CreatedAt:  *now(),
			}
			if completed {
				ses.CompletedAt = now()
			}
			return ses
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Session) {
			req := require.New(t)
			req.NoError(s.TruncateAutomationSessions(ctx))
			res := makeNew(0, false)
			req.NoError(s.CreateAutomationSession(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		wf := &types.Session{
			ID:        id.Next(),
			CreatedAt: *now(),
		}
		req.NoError(s.CreateAutomationSession(ctx, wf))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, wf := truncAndCreate(t)

		fetched, err := s.LookupAutomationSessionByID(ctx, wf.ID)
		req.NoError(err)
		req.Equal(wf.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, wf := truncAndCreate(t)
		wf.Status = types.SessionCompleted
		req.NoError(s.UpdateAutomationSession(ctx, wf))
		fetched, err := s.LookupAutomationSessionByID(ctx, wf.ID)
		req.NoError(err)
		req.Equal(wf.ID, fetched.ID)
		req.Equal(types.SessionCompleted, fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Session", func(t *testing.T) {
			req, wf := truncAndCreate(t)
			req.NoError(s.DeleteAutomationSession(ctx, wf))
			_, err := s.LookupAutomationSessionByID(ctx, wf.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, wf := truncAndCreate(t)
			req.NoError(s.DeleteAutomationSessionByID(ctx, wf.ID))
			_, err := s.LookupAutomationSessionByID(ctx, wf.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Session{
			makeNew(1001, false),
			makeNew(1001, false),
			makeNew(1001, false),
			makeNew(1001, false),
			makeNew(1001, true),
		}

		count := len(prefill)

		prefill[4].CompletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateAutomationSessions(ctx))
		req.NoError(s.CreateAutomationSession(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchAutomationSessions(ctx, types.SessionFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchAutomationSessions(ctx, types.SessionFilter{Completed: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchAutomationSessions(ctx, types.SessionFilter{Completed: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		// find all prefixed
		set, f, err = s.SearchAutomationSessions(ctx, types.SessionFilter{WorkflowID: id.Strings(1001)})
		req.NoError(err)
		req.Len(set, 4)

		_ = f // dummy
	})
}
