package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/crusttech/crust/crm/types"
)

func TestModule(t *testing.T) {
	repository := Module().With(context.Background())

	fields, err := json.Marshal([]types.Field{
		types.Field{
			Name: "name",
			Type: "input",
		},
		types.Field{
			Name: "email",
			Type: "email",
		},
		types.Field{
			Name: "options",
			Type: "select_multi",
		},
		types.Field{
			Name: "description",
			Type: "text",
		},
	})
	assert(t, err == nil, "Error when encoding JSON fields: %+v", err)

	// the module object we're working with
	module := &types.Module{
		Name: "Test",
	}
	(&module.Fields).Scan(fields)

	prevModuleCount := 0

	{
		{
			m, err := repository.Update(module)
			assert(t, m == nil, "Expected empty return for ivalid update, got %#v", m)
			assert(t, err != nil, "Expected error when updating invalid content")
		}

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

		// fetch module fields
		{
			fl, err := repository.FieldNames(m)
			assert(t, err == nil, "Error when retrieving module fields by module: %+v", err)
			assert(t, len(fl) == 4, "Expected 4 fields, got %d", len(fl))
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
			assert(t, len(ms) >= 1, "Expected at least one module, got %d", len(ms))
			prevModuleCount = len(ms)
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
			assert(t, len(ms) < prevModuleCount, "Expected modules count to decrease after deletion, %d < %d", len(ms), prevModuleCount)
		}
	}
}
