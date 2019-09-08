// +build integration

package service

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestPage(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing", true)

	// Set fake Identity (required for permission checks).
	ctx = auth.SetIdentityToContext(ctx, auth.NewIdentity(1337))

	ns1, _ := createTestNamespaces(ctx, t)

	svc := Page().With(ctx)

	// the page object we're working with
	var err error
	page := &types.Page{
		NamespaceID: ns1.ID,
		Title:       "Test",
		ModuleID:    123,
	}
	(&page.Blocks).Scan([]byte("[]"))

	prevPageCount := 0

	{
		{
			m, err := svc.Update(page)
			test.Assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			test.Assert(t, err != nil, "Expected error when updating invalid content")
		}

		// create page
		page, err = svc.Create(page)
		test.Assert(t, err == nil, "Error when creating page: %+v", err)
		test.Assert(t, page.ID > 0, "Expected auto generated ID")

		var firstPageID = page.ID

		page.SelfID = page.ID

		{
			_, err = svc.Create(page)
			test.Assert(t, err != nil, "%+v", errors.Errorf("Expected error when creating duplicate moduleID page"))
		}

		page.ModuleID = 0

		{
			_, err = svc.Create(page)
			test.Assert(t, err == nil, "Unexpected error when creating page, %+v", err)
		}
		{
			_, err = svc.Create(page)
			test.Assert(t, err == nil, "Unexpected error when creating page, %+v", err)
		}

		// fetch created page
		{
			p, err := svc.FindByID(page.NamespaceID, page.ID)
			test.Assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			test.Assert(t, p.ID == page.ID, "Expected ID from database to match, %+v", errors.Errorf("%d != %d", page.ID, p.ID))
			test.Assert(t, p.Title == page.Title, "Expected Title from database to match, %+v", errors.Errorf("%s != %s", page.Title, p.Title))
		}

		// update created page
		{
			page.Title = "Updated test"
			page.UpdatedAt = nil
			_, err := svc.Update(page)
			test.Assert(t, err == nil, "Error when updating page, %+v", err)
		}

		// re-fetch page
		{
			p, err := svc.FindByID(page.NamespaceID, page.ID)
			test.Assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			test.Assert(t, p.ID == page.ID, "re-fetch: Expected ID from database to match, %d != %d", page.ID, p.ID)
			test.Assert(t, p.Title == page.Title, "Expected Title from database to match, %s != %s", page.Title, p.Title)
		}

		// fetch all pages
		{
			p, _, err := svc.FindBySelfID(page.NamespaceID, 0)
			test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
			test.Assert(t, len(p) >= 1, "Expected at least one page, got %d", len(p))
			prevPageCount = len(p)
		}

		// fetch all pages
		{
			p, _, err := svc.FindBySelfID(page.NamespaceID, firstPageID)
			test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
			test.Assert(t, len(p) == 2, "Expected 2 pages with selfID=%d, got %v", page.ID, spew.Sdump(p))
			prevPageCount = len(p)

			ids := []uint64{p[0].ID, p[1].ID}

			{
				err := svc.Reorder(page.NamespaceID, firstPageID, ids)
				test.Assert(t, err == nil, "Error when reordering pages: %+v", err)

				p, _, err = svc.FindBySelfID(page.NamespaceID, firstPageID)
				test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
				test.Assert(t, len(p) == 2, "Expected 2 pages with selfID=%d, got %v", page.ID, spew.Sdump(p))
				test.Assert(t, p[0].Weight < p[1].Weight, "Expected ascending order, %+v", errors.Errorf("%d < %d", p[0].Weight, p[1].Weight))
			}
		}

		// re-fetch page
		{
			err := svc.DeleteByID(page.NamespaceID, page.ID)
			test.Assert(t, err == nil, "Error when deleting page by id: %+v", err)
		}

		// fetch all pages
		{
			p, _, err := svc.FindBySelfID(page.NamespaceID, 0)
			test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
			test.Assert(t, len(p) < prevPageCount, "Expected pages count to decrease after deletion, %d < %d", len(p), prevPageCount)
		}
	}

}
