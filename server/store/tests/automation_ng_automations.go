package tests

import (
	"context"
	"testing"
	"time"

	automationType "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/store"
	"github.com/stretchr/testify/require"
)

func testAutomationNgAutomations(t *testing.T, s store.AutomationNgAutomations) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *automationType.NgAutomation {
			return &automationType.NgAutomation{
				ID:        id.Next(),
				Handle:    handle,
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *automationType.NgAutomation) {
			req := require.New(t)
			req.NoError(s.TruncateAutomationNgAutomations(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateAutomationNgAutomation(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateAutomationNgAutomations(ctx))
		a := makeNew("NgAutomationCRUD")
		req.NoError(s.CreateAutomationNgAutomation(ctx, a))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, a := truncAndCreate(t)
		fetched, err := s.LookupAutomationNgAutomationByID(ctx, a.ID)
		req.NoError(err)
		req.Equal(a.Handle, fetched.Handle)
		req.Equal(a.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, a := truncAndCreate(t)
		a.Handle = string(rand.Bytes(10))
		req.NoError(s.UpdateAutomationNgAutomation(ctx, a))
		fetched, err := s.LookupAutomationNgAutomationByID(ctx, a.ID)
		req.NoError(err)
		req.Equal(a.Handle, fetched.Handle)
	})

	t.Run("delete", func(t *testing.T) {
		req, a := truncAndCreate(t)
		req.NoError(s.DeleteAutomationNgAutomationByID(ctx, a.ID))
	})
}
