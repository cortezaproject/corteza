package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testDmlConnections(t *testing.T, s store.DmlConnections) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.DmlConnection {
			return &types.DmlConnection{
				ID:        id.Next(),
				Handle:    handle,
				Label:     "test",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.DmlConnection) {
			req := require.New(t)
			req.NoError(s.TruncateDmlConnections(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateDmlConnection(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateDmlConnections(ctx))
		req.NoError(s.CreateDmlConnection(ctx, makeNew("TestDmlConnection")))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, c := truncAndCreate(t)
		fetched, err := s.LookupDmlConnectionByID(ctx, c.ID)
		req.NoError(err)
		req.Equal(c.ID, fetched.ID)
		req.Equal(c.Handle, fetched.Handle)
	})

	t.Run("update", func(t *testing.T) {
		req, c := truncAndCreate(t)
		c.Label = "updated"
		req.NoError(s.UpdateDmlConnection(ctx, c))
		fetched, err := s.LookupDmlConnectionByID(ctx, c.ID)
		req.NoError(err)
		req.Equal("updated", fetched.Label)
	})

	t.Run("delete", func(t *testing.T) {
		req, c := truncAndCreate(t)
		req.NoError(s.DeleteDmlConnectionByID(ctx, c.ID))
	})
}
