// +build integration

package service

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/test"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestPage(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing", true)

	user := &systemTypes.User{
		ID:       1337,
		Name:     "John Crm Doe",
		Username: "johndoe",
		SatosaID: "12345",
	}

	// Set Identity (required for permission checks).
	ctx = auth.SetIdentityToContext(ctx, user)

	repository := Page().With(ctx)

	// the page object we're working with
	page := &types.Page{
		Title:    "Test",
		ModuleID: 123,
	}
	(&page.Blocks).Scan([]byte("[]"))

	prevPageCount := 0

	{
		{
			m, err := repository.Update(page)
			test.Assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			test.Assert(t, err != nil, "Expected error when updating invalid content")
		}

		// create page
		m, err := repository.Create(page)
		test.Assert(t, err == nil, "Error when creating page: %+v", err)
		test.Assert(t, m.ID > 0, "Expected auto generated ID")

		page.SelfID = m.ID

		{
			_, err := repository.Create(page)
			test.Assert(t, err != nil, "%+v", errors.Errorf("Expected error when creating duplicate moduleID page"))
		}

		{
			page.ModuleID = 0
			_, err := repository.Create(page)
			test.Assert(t, err == nil, "Unexpected error when creating page, %+v", err)
		}
		{
			_, err := repository.Create(page)
			test.Assert(t, err == nil, "Unexpected error when creating page, %+v", err)
		}

		// fetch created page
		{
			ms, err := repository.FindByID(m.ID)
			test.Assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "Expected ID from database to match, %+v", errors.Errorf("%d != %d", m.ID, ms.ID))
			test.Assert(t, ms.Title == m.Title, "Expected Title from database to match, %+v", errors.Errorf("%s != %s", m.Title, ms.Title))
		}

		// update created page
		{
			m.Title = "Updated test"
			_, err := repository.Update(m)
			test.Assert(t, err == nil, "Error when updating page, %+v", err)
		}

		// re-fetch page
		{
			ms, err := repository.FindByID(m.ID)
			test.Assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "re-fetch: Expected ID from database to match, %d != %d", m.ID, ms.ID)
			test.Assert(t, ms.Title == m.Title, "Expected Title from database to match, %s != %s", m.Title, ms.Title)
		}

		// re-fetch page with moduleID
		{
			ms, err := repository.FindByModuleID(m.ModuleID)
			test.Assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "fetch-module: Expected ID from database to match, %d != %d", m.ID, ms.ID)
			test.Assert(t, ms.Title == m.Title, "Expected Title from database to match, %s != %s", m.Title, ms.Title)
		}

		// fetch all pages
		{
			ms, err := repository.FindBySelfID(0)
			test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
			test.Assert(t, len(ms) >= 1, "Expected at least one page, got %d", len(ms))
			prevPageCount = len(ms)
		}

		// fetch all record pages
		{
			ms, err := repository.FindRecordPages()
			test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
			test.Assert(t, len(ms) >= 1, "Expected at least one page, got %d", len(ms))
			prevPageCount = len(ms)
		}

		// fetch all pages
		{
			ms, err := repository.FindBySelfID(m.ID)
			test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
			test.Assert(t, len(ms) == 2, "Expected two pages with selfID=%d, got %v", m.ID, spew.Sdump(ms))
			prevPageCount = len(ms)

			parent := m.ID
			ids := []uint64{ms[0].ID, ms[1].ID}

			{
				err := repository.Reorder(parent, ids)
				test.Assert(t, err == nil, "Error when reordering pages: %+v", err)

				ms, err = repository.FindBySelfID(m.ID)
				test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
				test.Assert(t, len(ms) == 2, "Expected two pages with selfID=%d, got %v", m.ID, spew.Sdump(ms))
				test.Assert(t, ms[0].Weight < ms[1].Weight, "Expected ascending order, %+v", errors.Errorf("%d < %d", ms[0].Weight, ms[1].Weight))
			}
		}

		// re-fetch page
		{
			err := repository.DeleteByID(m.ID)
			test.Assert(t, err == nil, "Error when deleting page by id: %+v", err)
		}

		// fetch all pages
		{
			ms, err := repository.FindBySelfID(0)
			test.Assert(t, err == nil, "Error when retrieving pages: %+v", err)
			test.Assert(t, len(ms) < prevPageCount, "Expected pages count to decrease after deletion, %d < %d", len(ms), prevPageCount)
		}
	}

}
