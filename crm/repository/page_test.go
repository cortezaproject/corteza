package repository

import (
	"context"
	"github.com/crusttech/crust/crm/types"
	"testing"
)

func TestPage(t *testing.T) {
	repository := Page(context.TODO(), nil).With(context.Background(), nil)

	// the page object we're working with
	page := &types.Page{
		Title:    "Test",
		ModuleID: 123,
	}
	(&page.Blocks).Scan([]byte("[]"))

	prevPageCount := 0

	{
		// create page
		m, err := repository.Create(page)
		assert(t, err == nil, "Error when creating page: %+v", err)
		assert(t, m.ID > 0, "Expected auto generated ID")

		// fetch created page
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.Title == m.Title, "Expected Title from database to match, %s != %s", m.Title, ms.Title)
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
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.Title == m.Title, "Expected Title from database to match, %s != %s", m.Title, ms.Title)
		}

		// re-fetch page with moduleID
		{
			ms, err := repository.FindByModuleID(m.ModuleID)
			assert(t, err == nil, "Error when retrieving page by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.Title == m.Title, "Expected Title from database to match, %s != %s", m.Title, ms.Title)
		}

		// fetch all pages
		{
			ms, err := repository.Find()
			assert(t, err == nil, "Error when retrieving pages: %+v", err)
			assert(t, len(ms) >= 1, "Expected at least one page, got %d", len(ms))
			prevPageCount = len(ms)
		}

		// re-fetch page
		{
			err := repository.DeleteByID(m.ID)
			assert(t, err == nil, "Error when deleting page by id: %+v", err)
		}

		// fetch all pages
		{
			ms, err := repository.Find()
			assert(t, err == nil, "Error when retrieving pages: %+v", err)
			assert(t, len(ms) < prevPageCount, "Expected pages count to decrease after deletion, %d < %d", len(ms), prevPageCount)
		}
	}

}
