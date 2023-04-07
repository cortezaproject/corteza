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

func testComposePageLayouts(t *testing.T, s store.ComposePageLayouts) {
	var (
		ctx = context.Background()
		req = require.New(t)

		namespaceID = id.Next()

		makeNew = func(title, handle string) *types.PageLayout {
			// minimum data set for new composePageLayout
			return &types.PageLayout{
				ID:          id.Next(),
				NamespaceID: namespaceID,
				CreatedAt:   time.Now(),
				Meta: types.PageLayoutMeta{
					Title: title,
				},
				Handle: handle,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.PageLayout) {
			req := require.New(t)
			req.NoError(s.TruncateComposePageLayouts(ctx))
			res := makeNew(string(rand.Bytes(10)), string(rand.Bytes(10)))
			req.NoError(s.CreateComposePageLayout(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		composePageLayout := makeNew("ComposePageLayoutCRUD", "compose-page-crud")
		req.NoError(s.CreateComposePageLayout(ctx, composePageLayout))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup", func(t *testing.T) {
		t.Run("lookup by ID", func(t *testing.T) {
			req, composePageLayout := truncAndCreate(t)
			fetched, err := s.LookupComposePageLayoutByID(ctx, composePageLayout.ID)
			req.NoError(err)
			req.Equal(composePageLayout.Meta.Title, fetched.Meta.Title)
			req.Equal(composePageLayout.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})

		t.Run("lookup by NamespaceID, PageID, Handle", func(t *testing.T) {
			req, composePageLayout := truncAndCreate(t)
			fetched, err := s.LookupComposePageLayoutByNamespaceIDPageIDHandle(ctx, composePageLayout.NamespaceID, composePageLayout.PageID, composePageLayout.Handle)
			req.NoError(err)
			req.Equal(composePageLayout.Meta.Title, fetched.Meta.Title)
			req.Equal(composePageLayout.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})

		t.Run("lookup by NamespaceID, Handle", func(t *testing.T) {
			req, composePageLayout := truncAndCreate(t)
			fetched, err := s.LookupComposePageLayoutByNamespaceIDHandle(ctx, composePageLayout.NamespaceID, composePageLayout.Handle)
			req.NoError(err)
			req.Equal(composePageLayout.Meta.Title, fetched.Meta.Title)
			req.Equal(composePageLayout.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})
	})

	t.Run("update", func(t *testing.T) {
		composePageLayout := makeNew("update me", "update-me")
		req.NoError(s.CreateComposePageLayout(ctx, composePageLayout))

		composePageLayout = &types.PageLayout{
			ID:        composePageLayout.ID,
			CreatedAt: composePageLayout.CreatedAt,
			Meta: types.PageLayoutMeta{
				Title: "ComposePageLayoutCRUD+2",
			},
		}
		req.NoError(s.UpdateComposePageLayout(ctx, composePageLayout))

		updated, err := s.LookupComposePageLayoutByID(ctx, composePageLayout.ID)
		req.NoError(err)
		req.Equal(composePageLayout.Meta.Title, updated.Meta.Title)
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, composePageLayout := truncAndCreate(t)
			composePageLayout.Meta.Title = "ComposePageLayoutCRUD+2"

			req.NoError(s.UpsertComposePageLayout(ctx, composePageLayout))

			upserted, err := s.LookupComposePageLayoutByID(ctx, composePageLayout.ID)
			req.NoError(err)
			req.Equal(composePageLayout.Meta.Title, upserted.Meta.Title)
		})

		t.Run("new", func(t *testing.T) {
			composePageLayout := makeNew("upsert me", "upsert-me")
			composePageLayout.Meta.Title = "ComposePageLayoutCRUD+2"

			req.NoError(s.UpsertComposePageLayout(ctx, composePageLayout))

			upserted, err := s.LookupComposePageLayoutByID(ctx, composePageLayout.ID)
			req.NoError(err)
			req.Equal(composePageLayout.Meta.Title, upserted.Meta.Title)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by PageLayout", func(t *testing.T) {
			req, composePageLayout := truncAndCreate(t)
			req.NoError(s.DeleteComposePageLayout(ctx, composePageLayout))
			_, err := s.LookupComposePageLayoutByID(ctx, composePageLayout.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, composePageLayout := truncAndCreate(t)
			req.NoError(s.DeleteComposePageLayoutByID(ctx, composePageLayout.ID))
			_, err := s.LookupComposePageLayoutByID(ctx, composePageLayout.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.PageLayout{
			makeNew("/one-one", "page-1-1"),
			makeNew("/one-two", "page-1-2"),
			makeNew("/two-one", "page-2-1"),
			makeNew("/two-two", "page-2-2"),
			makeNew("/two-deleted", "page-2-d"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateComposePageLayouts(ctx))
		req.NoError(s.CreateComposePageLayout(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchComposePageLayouts(ctx, types.PageLayoutFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchComposePageLayouts(ctx, types.PageLayoutFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchComposePageLayouts(ctx, types.PageLayoutFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, f, err = s.SearchComposePageLayouts(ctx, types.PageLayoutFilter{Handle: "page-2-1"})
		req.NoError(err)
		req.Len(set, 1)

		// find all prefixed
		set, f, err = s.SearchComposePageLayouts(ctx, types.PageLayoutFilter{Query: "/two-"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}
