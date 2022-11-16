package tests

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testComposePages(t *testing.T, s store.ComposePages) {
	var (
		ctx = context.Background()
		req = require.New(t)

		namespaceID = id.Next()
		moduleID    = id.Next()

		makeNew = func(title, handle string) *types.Page {
			// minimum data set for new composePage
			return &types.Page{
				ID:          id.Next(),
				NamespaceID: namespaceID,
				ModuleID:    moduleID,
				CreatedAt:   time.Now(),
				Title:       title,
				Handle:      handle,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Page) {
			req := require.New(t)
			req.NoError(s.TruncateComposePages(ctx))
			res := makeNew(string(rand.Bytes(10)), string(rand.Bytes(10)))
			req.NoError(s.CreateComposePage(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		composePage := makeNew("ComposePageCRUD", "compose-page-crud")
		req.NoError(s.CreateComposePage(ctx, composePage))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup", func(t *testing.T) {
		t.Run("lookup by ID", func(t *testing.T) {
			req, composePage := truncAndCreate(t)
			fetched, err := s.LookupComposePageByID(ctx, composePage.ID)
			req.NoError(err)
			req.Equal(composePage.Title, fetched.Title)
			req.Equal(composePage.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})

		t.Run("lookup by NamespaceID, ModuleID", func(t *testing.T) {
			req, composePage := truncAndCreate(t)
			fetched, err := s.LookupComposePageByNamespaceIDModuleID(ctx, composePage.NamespaceID, composePage.ModuleID)
			req.NoError(err)
			req.Equal(composePage.Title, fetched.Title)
			req.Equal(composePage.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})

		t.Run("lookup by NamespaceID, Handle", func(t *testing.T) {
			req, composePage := truncAndCreate(t)
			fetched, err := s.LookupComposePageByNamespaceIDHandle(ctx, composePage.NamespaceID, composePage.Handle)
			req.NoError(err)
			req.Equal(composePage.Title, fetched.Title)
			req.Equal(composePage.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})
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

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, composePage := truncAndCreate(t)
			composePage.Title = "ComposePageCRUD+2"

			req.NoError(s.UpsertComposePage(ctx, composePage))

			upserted, err := s.LookupComposePageByID(ctx, composePage.ID)
			req.NoError(err)
			req.Equal(composePage.Title, upserted.Title)
		})

		t.Run("new", func(t *testing.T) {
			composePage := makeNew("upsert me", "upsert-me")
			composePage.Title = "ComposePageCRUD+2"

			req.NoError(s.UpsertComposePage(ctx, composePage))

			upserted, err := s.LookupComposePageByID(ctx, composePage.ID)
			req.NoError(err)
			req.Equal(composePage.Title, upserted.Title)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Page", func(t *testing.T) {
			req, composePage := truncAndCreate(t)
			req.NoError(s.DeleteComposePage(ctx, composePage))
			_, err := s.LookupComposePageByID(ctx, composePage.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, composePage := truncAndCreate(t)
			req.NoError(s.DeleteComposePageByID(ctx, composePage.ID))
			_, err := s.LookupComposePageByID(ctx, composePage.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
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
