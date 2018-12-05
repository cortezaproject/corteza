package service

import (
	"context"
	"testing"

	"encoding/json"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/auth"
	systemRepository "github.com/crusttech/crust/system/repository"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestContent(t *testing.T) {
	user := &systemTypes.User{
		ID:       1337,
		Username: "TestUser",
	}
	{
		err := user.GeneratePassword("Mary had a little lamb, little lamb, little lamb")
		assert(t, err == nil, "Error generating password: %+v", err)
	}

	{
		userAPI := systemRepository.User(context.Background(), nil)
		_, err := userAPI.Create(user)
		assert(t, err == nil, "Error when inserting user: %+v", err)
	}

	ctx := auth.SetIdentityToContext(context.Background(), auth.NewIdentity(user.Identity()))
	repository := Content().With(ctx)

	module := &types.Module{
		Name: "Test",
		Fields: []types.ModuleField{
			types.ModuleField{
				Name: "name",
				Kind: "input",
			},
			types.ModuleField{
				Name: "email",
				Kind: "email",
			},
			types.ModuleField{
				Name: "options",
				Kind: "select_multi",
			},
		},
	}

	// set up a module
	{
		_, err := Module().With(context.Background()).Create(module)
		assert(t, err == nil, "Error when creating module: %+v", err)
		assert(t, module.ID > 0, "Expected auto generated ID")
	}

	columns := []types.ContentColumn{
		types.ContentColumn{
			Name:  "name",
			Value: "Tit Petric",
		},
		types.ContentColumn{
			Name:  "email",
			Value: "tit.petric@example.com",
		},
		types.ContentColumn{
			Name:    "options",
			Related: []string{"1", "2", "3"},
		},
	}

	content := &types.Content{
		ModuleID: module.ID,
	}
	(&content.Fields).Scan(func() []byte {
		b, _ := json.Marshal(columns)
		return b
	}())

	// now work with content
	{
		{
			m, err := repository.Update(content)
			assert(t, m == nil, "Expected empty return for ivalid update, got %#v", m)
			assert(t, err != nil, "Expected error when updating invalid content")
		}

		// create content
		m, err := repository.Create(content)
		assert(t, err == nil, "Error when creating content: %+v", err)
		assert(t, m.ID > 0, "Expected auto generated ID")
		assert(t, m.User != nil, "Expected non-nil user when creating content")
		assert(t, m.User.Username == "TestUser", "Expected 'TestUser' as username, got '%s'", m.User.Username)

		// fetch created content
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving content by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.ModuleID == m.ModuleID, "Expected Module ID from database to match, %d != %d", m.ModuleID, ms.ModuleID)

			{
				fields, err := repository.Fields(ms)
				// fields := make([]testContentRow, 0)
				// err = json.Unmarshal(ms.Fields, &fields)
				assert(t, err == nil, "%+v", errors.Wrap(err, "Didn't expect error when unmarshalling"))
				assert(t, len(fields) == len(columns), "Expected different field count: %d != %d", 2, len(fields))
				for k, v := range columns {
					assert(t, fields[k].Name == v.Name, "Expected fields[%d].Name = %s, got %s", k, fields[k].Name, v.Name)
				}
			}
			{
				fields := make([]types.ContentColumn, 0)
				err := json.Unmarshal(ms.Fields, &fields)
				assert(t, err == nil, "%+v", errors.Wrap(err, "Didn't expect error when unmarshalling"))
				assert(t, len(fields) == len(columns), "Expected different field count: %d != %d", 2, len(fields))
				for k, v := range columns {
					assert(t, fields[k].Name == v.Name, "Expected fields[%d].Name = %s, got %s", k, fields[k].Name, v.Name)
				}
			}
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
			mr, err := repository.Find(module.ID, "", 0, 20)
			assert(t, err == nil, "Error when retrieving contents: %+v", err)
			assert(t, len(mr.Contents) == 1, "Expected one content, got %d", len(mr.Contents))
			assert(t, mr.Meta.Count == 1, "Expected Meta.Count == 1, got %d", mr.Meta.Count)
			assert(t, mr.Contents[0].ModuleID == m.ModuleID, "Expected content module to match, %d != %d", m.ModuleID, mr.Contents[0].ModuleID)
		}

		// fetch all contents by query
		{
			mr, err := repository.Find(module.ID, "petric", 0, 20)
			assert(t, err == nil, "Error when retrieving contents: %+v", err)
			assert(t, len(mr.Contents) == 1, "Expected one content, got %d", len(mr.Contents))
			assert(t, mr.Meta.Count == 1, "Expected Meta.Count == 1, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Page == 0, "Expected Meta.Page == 0, got %d", mr.Meta.Page)
			assert(t, mr.Meta.PerPage == 20, "Expected Meta.PerPage == 20, got %d", mr.Meta.PerPage)
			assert(t, mr.Meta.Query == "petric", "Expected Meta.Query == petric, got '%s'", mr.Meta.Query)
		}

		// fetch all contents by query
		{
			mr, err := repository.Find(module.ID, "niall", 0, 20)
			assert(t, err == nil, "Error when retrieving contents: %+v", err)
			assert(t, len(mr.Contents) == 0, "Expected no contents, got %d", len(mr.Contents))
		}

		// re-fetch content
		{
			err := repository.DeleteByID(m.ID)
			assert(t, err == nil, "Error when retrieving content by id: %+v", err)
		}

		// fetch all contents
		{
			mr, err := repository.Find(module.ID, "", 0, 20)
			assert(t, err == nil, "Error when retrieving contents: %+v", err)
			assert(t, len(mr.Contents) == 0, "Expected no content, got %d", len(mr.Contents))
		}
	}

}
