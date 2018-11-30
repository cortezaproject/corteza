package service

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/crm/types"
)

func TestPage(t *testing.T) {
	repository := Page().With(context.Background())

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
			assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			assert(t, err != nil, "Expected error when updating invalid content")
		}

		// create page
		m, err := repository.Create(page)
		assert(t, err == nil, "Error when creating page: %+v", err)
		assert(t, m.ID > 0, "Expected auto generated ID")

		page.SelfID = m.ID

		{
			_, err := repository.Create(page)
			assert(t, err != nil, "%+v", errors.Errorf("Expected error when creating duplicate moduleID page"))
		}

		{
			page.ModuleID = 0
			_, err := repository.Create(page)
			assert(t, err == nil, "Unexpected error when creating page, %+v", err)
		}
		{
			_, err := repository.Create(page)
			assert(t, err == nil, "Unexpected error when creating page, %+v", err)
		}

		// fetch created page
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %+v", errors.Errorf("%d != %d", m.ID, ms.ID))
			assert(t, ms.Title == m.Title, "Expected Title from database to match, %+v", errors.Errorf("%s != %s", m.Title, ms.Title))
		}

		// update created page
		{
			m.Title = "Updated test"
			_, err := repository.Update(m)
			assert(t, err == nil, "Error when updating page, %+v", err)
		}

		// re-fetch page
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			assert(t, ms.ID == m.ID, "re-fetch: Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.Title == m.Title, "Expected Title from database to match, %s != %s", m.Title, ms.Title)
		}

		// re-fetch page with moduleID
		{
			ms, err := repository.FindByModuleID(m.ModuleID)
			assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			assert(t, ms.ID == m.ID, "fetch-module: Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.Title == m.Title, "Expected Title from database to match, %s != %s", m.Title, ms.Title)
		}

		// fetch all pages
		{
			ms, err := repository.FindBySelfID(0)
			assert(t, err == nil, "Error when retrieving pages: %+v", err)
			assert(t, len(ms) >= 1, "Expected at least one page, got %d", len(ms))
			prevPageCount = len(ms)
		}

		// fetch all record pages
		{
			ms, err := repository.FindRecordPages()
			assert(t, err == nil, "Error when retrieving pages: %+v", err)
			assert(t, len(ms) >= 1, "Expected at least one page, got %d", len(ms))
			prevPageCount = len(ms)
		}

		// fetch all pages
		{
			ms, err := repository.FindBySelfID(m.ID)
			assert(t, err == nil, "Error when retrieving pages: %+v", err)
			assert(t, len(ms) == 2, "Expected two pages with selfID=%d, got %v", m.ID, spew.Sdump(ms))
			prevPageCount = len(ms)

			parent := m.ID
			ids := []uint64{ms[0].ID, ms[1].ID}

			{
				err := repository.Reorder(parent, ids)
				assert(t, err == nil, "Error when reordering pages: %+v", err)

				ms, err = repository.FindBySelfID(m.ID)
				assert(t, err == nil, "Error when retrieving pages: %+v", err)
				assert(t, len(ms) == 2, "Expected two pages with selfID=%d, got %v", m.ID, spew.Sdump(ms))
				assert(t, ms[0].Weight < ms[1].Weight, "Expected ascending order, %+v", errors.Errorf("%d < %d", ms[0].Weight, ms[1].Weight))
			}
		}

		// re-fetch page
		{
			err := repository.DeleteByID(m.ID)
			assert(t, err == nil, "Error when deleting page by id: %+v", err)
		}

		// fetch all pages
		{
			ms, err := repository.FindBySelfID(0)
			assert(t, err == nil, "Error when retrieving pages: %+v", err)
			assert(t, len(ms) < prevPageCount, "Expected pages count to decrease after deletion, %d < %d", len(ms), prevPageCount)
		}
	}

}
