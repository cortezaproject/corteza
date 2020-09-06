package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testComposePages(t *testing.T, s store.ComposePages) {
	var (
		ctx = context.Background()
		req = require.New(t)

		namespaceID = id.Next()

		makeNew = func(title, handle string) *types.Page {
			// minimum data set for new composePage
			return &types.Page{
				ID:          id.Next(),
				NamespaceID: namespaceID,
				CreatedAt:   time.Now(),
				Title:       title,
				Handle:      handle,
			}
		}
	)

	t.Run("create", func(t *testing.T) {
		composePage := makeNew("ComposePageCRUD", "compose-page-crud")
		req.NoError(s.CreateComposePage(ctx, composePage))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		composePage := makeNew("look up by id", "look-up-by-id")
		req.NoError(s.CreateComposePage(ctx, composePage))
		fetched, err := s.LookupComposePageByID(ctx, composePage.ID)
		req.NoError(err)
		req.Equal(composePage.Title, fetched.Title)
		req.Equal(composePage.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		composePage := makeNew("Delete", "Delete")
		req.NoError(s.CreateComposePage(ctx, composePage))
		req.NoError(s.DeleteComposePage(ctx))
	})

	t.Run("Delete by ID", func(t *testing.T) {
		composePage := makeNew("Delete by id", "Delete-by-id")
		req.NoError(s.CreateComposePage(ctx, composePage))
		req.NoError(s.DeleteComposePage(ctx))
	})

	t.Run("update", func(t *testing.T) {
		composePage := makeNew("update me", "update-me")
		req.NoError(s.CreateComposePage(ctx, composePage))

		composePage = &types.Page{
			ID:        composePage.ID,
			CreatedAt: composePage.CreatedAt,
			Title:     "ComposePageCRUD+2",
		}
		req.NoError(s.UpdateComposePage(ctx, composePage))

		updated, err := s.LookupComposePageByID(ctx, composePage.ID)
		req.NoError(err)
		req.Equal(composePage.Title, updated.Title)
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Page{
			makeNew("/one-one", "page-1-1"),
			makeNew("/one-two", "page-1-2"),
			makeNew("/two-one", "page-2-1"),
			makeNew("/two-two", "page-2-2"),
			makeNew("/two-deleted", "page-2-d"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateComposePages(ctx))
		req.NoError(s.CreateComposePage(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchComposePages(ctx, types.PageFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchComposePages(ctx, types.PageFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchComposePages(ctx, types.PageFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, f, err = s.SearchComposePages(ctx, types.PageFilter{Handle: "page-2-1"})
		req.NoError(err)
		req.Len(set, 1)

		// find all prefixed
		set, f, err = s.SearchComposePages(ctx, types.PageFilter{Query: "/two-"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}
