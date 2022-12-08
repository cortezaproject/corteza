package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
)

func testComposeNamespaces(t *testing.T, s store.ComposeNamespaces) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(name, slug string) *types.Namespace {
			// minimum data set for new composeNamespace
			return &types.Namespace{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Name:      name,
				Slug:      slug,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Namespace) {
			req := require.New(t)
			req.NoError(s.TruncateComposeNamespaces(ctx))
			res := makeNew(string(rand.Bytes(10)), string(rand.Bytes(10)))
			req.NoError(s.CreateComposeNamespace(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateComposeNamespaces(ctx))
		composeNamespace := makeNew("ComposeNamespaceCRUD", "compose-namespace-crud")
		req.NoError(s.CreateComposeNamespace(ctx, composeNamespace))
	})

	t.Run("create with duplicate slug", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup", func(t *testing.T) {
		t.Run("by ID", func(t *testing.T) {
			req, composeNamespace := truncAndCreate(t)
			fetched, err := s.LookupComposeNamespaceByID(ctx, composeNamespace.ID)
			req.NoError(err)
			req.Equal(composeNamespace.Name, fetched.Name)
			req.Equal(composeNamespace.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})

		t.Run("by Slug", func(t *testing.T) {
			req, composeNamespace := truncAndCreate(t)
			fetched, err := s.LookupComposeNamespaceBySlug(ctx, composeNamespace.Slug)
			req.NoError(err)
			req.Equal(composeNamespace.Name, fetched.Name)
			req.Equal(composeNamespace.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})
	})

	t.Run("update", func(t *testing.T) {
		composeNamespace := makeNew("update me", "update-me")
		req.NoError(s.CreateComposeNamespace(ctx, composeNamespace))

		composeNamespace = &types.Namespace{
			ID:        composeNamespace.ID,
			CreatedAt: composeNamespace.CreatedAt,
			Name:      "ComposeNamespaceCRUD+2",
		}
		req.NoError(s.UpdateComposeNamespace(ctx, composeNamespace))

		updated, err := s.LookupComposeNamespaceByID(ctx, composeNamespace.ID)
		req.NoError(err)
		req.Equal(composeNamespace.Name, updated.Name)
	})

	t.Run("update with duplicate slug", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, composeNamespace := truncAndCreate(t)
			composeNamespace.Name = "ComposeNamespaceCRUD+2"

			req.NoError(s.UpsertComposeNamespace(ctx, composeNamespace))

			upserted, err := s.LookupComposeNamespaceByID(ctx, composeNamespace.ID)
			req.NoError(err)
			req.Equal(composeNamespace.Name, upserted.Name)
		})

		t.Run("new", func(t *testing.T) {
			composeNamespace := makeNew("upsert me", "upsert-me")
			composeNamespace.Name = "ComposeNamespaceCRUD+3"

			req.NoError(s.UpsertComposeNamespace(ctx, composeNamespace))

			upserted, err := s.LookupComposeNamespaceByID(ctx, composeNamespace.ID)
			req.NoError(err)
			req.Equal(composeNamespace.Name, upserted.Name)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Namespace", func(t *testing.T) {
			req, composeNamespace := truncAndCreate(t)
			req.NoError(s.DeleteComposeNamespace(ctx, composeNamespace))
			_, err := s.LookupComposeNamespaceByID(ctx, composeNamespace.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, composeNamespace := truncAndCreate(t)
			req.NoError(s.DeleteComposeNamespaceByID(ctx, composeNamespace.ID))
			_, err := s.LookupComposeNamespaceByID(ctx, composeNamespace.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Namespace{
			makeNew("/one-one", "namespace-1-1"),
			makeNew("/one-two", "namespace-1-2"),
			makeNew("/two-one", "namespace-2-1"),
			makeNew("/two-two", "namespace-2-2"),
			makeNew("/two-deleted", "namespace-2-d"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateComposeNamespaces(ctx))
		req.NoError(s.CreateComposeNamespace(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchComposeNamespaces(ctx, types.NamespaceFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchComposeNamespaces(ctx, types.NamespaceFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchComposeNamespaces(ctx, types.NamespaceFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, f, err = s.SearchComposeNamespaces(ctx, types.NamespaceFilter{Slug: "namespace-2-1"})
		req.NoError(err)
		req.Len(set, 1)

		// find all prefixed
		set, f, err = s.SearchComposeNamespaces(ctx, types.NamespaceFilter{Query: "/two-"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}
