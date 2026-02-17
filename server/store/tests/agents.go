package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func testAgents(t *testing.T, s store.Agents) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.Agent {
			return &types.Agent{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
				Status:    "active",
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Agent) {
			req := require.New(t)
			req.NoError(s.TruncateAgents(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateAgent(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateAgents(ctx))
		agent := makeNew("AgentCRUD")
		req.NoError(s.CreateAgent(ctx, agent))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, agent := truncAndCreate(t)
		fetched, err := s.LookupAgentByID(ctx, agent.ID)
		req.NoError(err)
		req.Equal(agent.Handle, fetched.Handle)
		req.Equal(agent.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, agent := truncAndCreate(t)
		agent.Status = "disabled"
		req.NoError(s.UpdateAgent(ctx, agent))
		fetched, err := s.LookupAgentByID(ctx, agent.ID)
		req.NoError(err)
		req.Equal("disabled", fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, agent := truncAndCreate(t)
		req.NoError(s.DeleteAgentByID(ctx, agent.ID))
	})
}
