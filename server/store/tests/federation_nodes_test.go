package tests

import (
	"context"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testFederationNodes(t *testing.T, s store.FederationNodes) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(name string) *types.Node {
			// minimum data set for new FederationNode
			return &types.Node{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Name:      name,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Node) {
			req := require.New(t)
			req.NoError(s.TruncateFederationNodes(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateFederationNode(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		FederationNode := makeNew("FederationNodeCRUD")
		req.NoError(s.CreateFederationNode(ctx, FederationNode))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, FederationNode := truncAndCreate(t)
		fetched, err := s.LookupFederationNodeByID(ctx, FederationNode.ID)
		req.NoError(err)
		req.Equal(FederationNode.Name, fetched.Name)
		req.Equal(FederationNode.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, FederationNode := truncAndCreate(t)
		FederationNode.Name = "FederationNodeCRUD+2"

		req.NoError(s.UpdateFederationNode(ctx, FederationNode))

		updated, err := s.LookupFederationNodeByID(ctx, FederationNode.ID)
		req.NoError(err)
		req.Equal(FederationNode.Name, updated.Name)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, FederationNode := truncAndCreate(t)
			FederationNode.Name = "FederationNodeCRUD+2"

			req.NoError(s.UpsertFederationNode(ctx, FederationNode))

			updated, err := s.LookupFederationNodeByID(ctx, FederationNode.ID)
			req.NoError(err)
			req.Equal(FederationNode.Name, updated.Name)
		})

		t.Run("new", func(t *testing.T) {
			FederationNode := makeNew("upsert me")
			FederationNode.Name = "ComposeChartCRUD+2"

			req.NoError(s.UpsertFederationNode(ctx, FederationNode))

			upserted, err := s.LookupFederationNodeByID(ctx, FederationNode.ID)
			req.NoError(err)
			req.Equal(FederationNode.Name, upserted.Name)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by FederationNode", func(t *testing.T) {
			req, FederationNode := truncAndCreate(t)
			req.NoError(s.DeleteFederationNode(ctx, FederationNode))
			_, err := s.LookupFederationNodeByID(ctx, FederationNode.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, FederationNode := truncAndCreate(t)
			req.NoError(s.DeleteFederationNodeByID(ctx, FederationNode.ID))
			_, err := s.LookupFederationNodeByID(ctx, FederationNode.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Node{
			makeNew("/one-one"),
			makeNew("/one-two"),
			makeNew("/two-one"),
			makeNew("/two-two"),
			makeNew("/two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateFederationNodes(ctx))
		req.NoError(s.CreateFederationNode(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchFederationNodes(ctx, types.NodeFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchFederationNodes(ctx, types.NodeFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchFederationNodes(ctx, types.NodeFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		// find all prefixed
		set, f, err = s.SearchFederationNodes(ctx, types.NodeFilter{Query: "/two"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}
