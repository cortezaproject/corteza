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

func TestRecord(t *testing.T) {
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
	repository := Record().With(ctx)

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
			types.ModuleField{
				Name: "description",
				Kind: "text",
			},
		},
	}

	// set up a module
	{
		_, err := Module().With(context.Background()).Create(module)
		assert(t, err == nil, "Error when creating module: %+v", err)
		assert(t, module.ID > 0, "Expected auto generated ID")
	}

	columns := []types.RecordColumn{
		types.RecordColumn{
			Name:  "name",
			Value: "Tit Petric",
		},
		types.RecordColumn{
			Name:  "email",
			Value: "tit.petric@example.com",
		},
		types.RecordColumn{
			Name:    "options",
			Related: []string{"1", "2", "3"},
		},
		types.RecordColumn{
			Name:  "description",
			Value: "jack of all trades",
		},
	}

	record1 := &types.Record{
		ModuleID: module.ID,
	}
	(&record1.Fields).Scan(func() []byte {
		b, _ := json.Marshal(columns)
		return b
	}())

	columns2 := []types.RecordColumn{
		types.RecordColumn{
			Name:  "name",
			Value: "Marko Novak",
		},
		types.RecordColumn{
			Name:  "email",
			Value: "marko.n@example.com",
		},
		types.RecordColumn{
			Name:    "options",
			Related: []string{"1", "2", "3"},
		},
		types.RecordColumn{
			Name:  "description",
			Value: "persona non grata",
		},
	}

	record2 := &types.Record{
		ModuleID: module.ID,
	}
	(&record2.Fields).Scan(func() []byte {
		b, _ := json.Marshal(columns2)
		return b
	}())

	// now work with records
	{
		{
			m, err := repository.Update(record1)
			assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			assert(t, err != nil, "Expected error when updating invalid record")
		}

		// create record
		m1, err := repository.Create(record1)
		assert(t, err == nil, "Error when creating record: %+v", err)
		assert(t, m1.ID > 0, "Expected auto generated ID")
		assert(t, m1.User != nil, "Expected non-nil user when creating record")
		assert(t, m1.User.Username == "TestUser", "Expected 'TestUser' as username, got '%s'", m1.User.Username)

		// create record
		m2, err := repository.Create(record2)
		assert(t, err == nil, "Error when creating record: %+v", err)
		assert(t, m2.ID > 0, "Expected auto generated ID")
		assert(t, m2.User != nil, "Expected non-nil user when creating record")
		assert(t, m2.User.Username == "TestUser", "Expected 'TestUser' as username, got '%s'", m2.User.Username)

		// fetch created record
		{
			ms, err := repository.FindByID(m1.ID)
			assert(t, err == nil, "Error when retrieving record by id: %+v", err)
			assert(t, ms.ID == m1.ID, "Expected ID from database to match, %d != %d", m1.ID, ms.ID)
			assert(t, ms.ModuleID == m1.ModuleID, "Expected Module ID from database to match, %d != %d", m1.ModuleID, ms.ModuleID)

			{
				fields, err := repository.Fields(ms)
				// fields := make([]testRecordRow, 0)
				// err = json.Unmarshal(ms.Fields, &fields)
				assert(t, err == nil, "%+v", errors.Wrap(err, "Didn't expect error when unmarshalling"))
				assert(t, len(fields) == len(columns), "Expected different field count: %d != %d", 2, len(fields))
				for k, v := range columns {
					assert(t, fields[k].Name == v.Name, "Expected fields[%d].Name = %s, got %s", k, fields[k].Name, v.Name)
				}
			}
			{
				fields := make([]types.RecordColumn, 0)
				err := json.Unmarshal(ms.Fields, &fields)
				assert(t, err == nil, "%+v", errors.Wrap(err, "Didn't expect error when unmarshalling"))
				assert(t, len(fields) == len(columns), "Expected different field count: %d != %d", 2, len(fields))
				for k, v := range columns {
					assert(t, fields[k].Name == v.Name, "Expected fields[%d].Name = %s, got %s", k, fields[k].Name, v.Name)
				}
			}
		}

		// update created record
		{
			_, err := repository.Update(m1)
			assert(t, err == nil, "Error when updating record, %+v", err)
		}

		// re-fetch record
		{
			ms, err := repository.FindByID(m1.ID)
			assert(t, err == nil, "Error when retrieving record by id: %+v", err)
			assert(t, ms.ID == m1.ID, "Expected ID from database to match, %d != %d", m1.ID, ms.ID)
			assert(t, ms.ModuleID == m1.ModuleID, "Expected ID from database to match, %d != %d", m1.ModuleID, ms.ModuleID)
		}

		// fetch all records
		{
			mr, err := repository.Find(module.ID, "", 0, 20, "id desc")
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Sort == "id desc", "Expected Meta.Sort == id desc, got '%s'", mr.Meta.Sort)
			assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)
			assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be descending")
		}

		// fetch all records
		{
			mr, err := repository.Find(module.ID, "", 0, 20, "name asc, email desc")
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Sort == "name asc, email desc", "Expected Meta.Sort == 'name asc, email desc' '%s'", mr.Meta.Sort)
			assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)
			assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be ascending")
		}

		// fetch all records
		{
			mr, err := repository.Find(module.ID, "", 0, 20, "created_at desc")
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Sort == "created_at desc", "Expected Meta.Sort == created_at desc, got '%s'", mr.Meta.Sort)
			assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)
			assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be ascending")
		}

		// fetch all records by query
		{
			mr, err := repository.Find(module.ID, "petric", 0, 20, "id desc")
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 1, "Expected one record, got %d", len(mr.Records))
			assert(t, mr.Meta.Count == 1, "Expected Meta.Count == 1, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Page == 0, "Expected Meta.Page == 0, got %d", mr.Meta.Page)
			assert(t, mr.Meta.PerPage == 20, "Expected Meta.PerPage == 20, got %d", mr.Meta.PerPage)
			assert(t, mr.Meta.Query == "petric", "Expected Meta.Query == petric, got '%s'", mr.Meta.Query)
			assert(t, mr.Meta.Sort == "id desc", "Expected Meta.Sort == id desc, got '%s'", mr.Meta.Sort)
		}

		// fetch all records by query
		{
			mr, err := repository.Find(module.ID, "niall", 0, 20, "id asc")
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 0, "Expected no records, got %d", len(mr.Records))
		}

		// delete record
		{
			err := repository.DeleteByID(m1.ID)
			assert(t, err == nil, "Error when retrieving record by id: %+v", err)

			err = repository.DeleteByID(m2.ID)
			assert(t, err == nil, "Error when retrieving record by id: %+v", err)
		}

		// fetch all records
		{
			mr, err := repository.Find(module.ID, "", 0, 20, "")
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 0, "Expected no record, got %d", len(mr.Records))
		}
	}
}
