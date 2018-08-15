package repository

import (
	"context"
	"encoding/json"
	"github.com/crusttech/crust/crm/types"
	"testing"
)

type testContentRow struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}

func TestContent(t *testing.T) {
	repository := NewContent(context.TODO()).With(context.Background())

	// clean up tables
	{
		for _, name := range []string{"crm_module", "crm_module_content"} {
			_, err := db().Exec("truncate " + name)
			assert(t, err == nil, "Error when clearing "+name+": %+v", err)
		}
	}

	fields, err := json.Marshal([]types.Field{
		types.Field{
			Name: "name",
			Type: "input",
		},
		types.Field{
			Name: "email",
			Type: "email",
		},
	})
	assert(t, err == nil, "Error when encoding JSON fields: %+v", err)

	module := &types.Module{
		Name: "Test",
	}
	(&module.Fields).Scan(fields)

	// set up a module
	{
		_, err := NewModule(context.TODO()).With(context.Background()).Create(module)
		assert(t, err == nil, "Error when creating module: %+v", err)
		assert(t, module.ID > 0, "Expected auto generated ID")
	}

	content := &types.Content{
		ModuleID: module.ID,
	}
	(&content.Fields).Scan(func() []byte {
		b, _ := json.Marshal([]testContentRow{
			testContentRow{"name", "Tit Petric"},
			testContentRow{"email", "tit.petric@example.com"},
		})
		return b
	}())

	// now work with content
	{
		// create content
		m, err := repository.Create(content)
		assert(t, err == nil, "Error when creating content: %+v", err)
		assert(t, m.ID > 0, "Expected auto generated ID")

		// fetch created content
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving content by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.ModuleID == m.ModuleID, "Expected Module ID from database to match, %d != %d", m.ModuleID, ms.ModuleID)

			fields := make([]testContentRow, 0)
			err = json.Unmarshal(ms.Fields, &fields)
			assert(t, err == nil, "Didn't expect error when unmarshalling: %+v", err)
			assert(t, len(fields) == 2, "Expected different field count: %d != %d", 2, len(fields))
			assert(t, fields[0].Name == "name", "Expected field.0 type = name, got %s", fields[0].Name)
			assert(t, fields[1].Name == "email", "Expected field.1 type = email, got %s", fields[1].Name)
		}

		// update created content
		{
			_, err := repository.Update(m)
			assert(t, err == nil, "Error when updating content, %+v", err)
		}

		// re-fetch content
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving content by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.ModuleID == m.ModuleID, "Expected ID from database to match, %d != %d", m.ModuleID, ms.ModuleID)
		}

		// fetch all contents
		{
			ms, err := repository.Find()
			assert(t, err == nil, "Error when retrieving contents: %+v", err)
			assert(t, len(ms) == 1, "Expected one content, got %d", len(ms))
			assert(t, ms[0].ModuleID == m.ModuleID, "Expected content module to match, %s != %s", m.ModuleID, ms[0].ModuleID)
		}

		// re-fetch content
		{
			err := repository.DeleteByID(m.ID)
			assert(t, err == nil, "Error when retrieving content by id: %+v", err)
		}

		// fetch all contents
		{
			ms, err := repository.Find()
			assert(t, err == nil, "Error when retrieving contents: %+v", err)
			assert(t, len(ms) == 0, "Expected one content, got %d", len(ms))
		}
	}

}
