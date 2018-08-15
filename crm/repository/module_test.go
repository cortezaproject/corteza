package repository

import (
	"context"
	"github.com/crusttech/crust/crm/types"
	"testing"
)

func TestModule(t *testing.T) {

	repository := NewModule(context.TODO()).With(context.Background())

	// clean up tables
	{
		_, err := db().Exec("truncate crm_module")
		assert(t, err == nil, "Error when clearing crm_module: %+v", err)
	}

	// the module object we're working with
	module := &types.Module{
		Name: "Test",
	}

	{
		// create module
		m, err := repository.Create(module)
		assert(t, err == nil, "Error when creating module: %+v", err)
		assert(t, m.ID > 0, "Expected auto generated ID")

		// fetch created module
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving module by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.Name == m.Name, "Expected Name from database to match, %s != %s", m.Name, ms.Name)
		}

		// update created module
		{
			m.Name = "Updated test"
			_, err := repository.Update(m)
			assert(t, err == nil, "Error when updating module, %+v", err)
		}

		// re-fetch module
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving module by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.Name == m.Name, "Expected Name from database to match, %s != %s", m.Name, ms.Name)
		}

		// fetch all modules
		{
			ms, err := repository.Find()
			assert(t, err == nil, "Error when retrieving modules: %+v", err)
			assert(t, len(ms) == 1, "Expected one module, got %d", len(ms))
			assert(t, ms[0].Name == m.Name, "Expected module name to match, %s != %s", m.Name, ms[0].Name)
		}

		// re-fetch module
		{
			err := repository.DeleteByID(m.ID)
			assert(t, err == nil, "Error when deleting module by id: %+v", err)
		}

		// fetch all modules
		{
			ms, err := repository.Find()
			assert(t, err == nil, "Error when retrieving modules: %+v", err)
			assert(t, len(ms) == 0, "Expected no modules, got %d", len(ms))
		}
	}

}
