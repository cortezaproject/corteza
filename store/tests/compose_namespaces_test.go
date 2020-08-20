package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
	)

	t.Run("create", func(t *testing.T) {
		composeNamespace := makeNew("ComposeNamespaceCRUD", "compose-namespace-crud")
		req.NoError(s.CreateComposeNamespace(ctx, composeNamespace))
	})

	t.Run("create with duplicate slug", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		composeNamespace := makeNew("look up by id", "look-up-by-id")
		req.NoError(s.CreateComposeNamespace(ctx, composeNamespace))
		fetched, err := s.LookupComposeNamespaceByID(ctx, composeNamespace.ID)
		req.NoError(err)
		req.Equal(composeNamespace.Name, fetched.Name)
		req.Equal(composeNamespace.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		composeNamespace := makeNew("Delete", "Delete")
		req.NoError(s.CreateComposeNamespace(ctx, composeNamespace))
		req.NoError(s.DeleteComposeNamespace(ctx))
	})

	t.Run("Delete by ID", func(t *testing.T) {
		composeNamespace := makeNew("Delete by id", "Delete-by-id")
		req.NoError(s.CreateComposeNamespace(ctx, composeNamespace))
		req.NoError(s.DeleteComposeNamespace(ctx))
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
		set, f, err = s.SearchComposeNamespaces(ctx, types.NamespaceFilter{Deleted: rh.FilterStateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchComposeNamespaces(ctx, types.NamespaceFilter{Deleted: rh.FilterStateExclusive})
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
