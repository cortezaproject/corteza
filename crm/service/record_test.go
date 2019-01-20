package service

import (
	"context"
	"testing"

	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/test"
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
		Fields: types.ModuleFieldSet{
			&types.ModuleField{
				Name: "name",
			},
			&types.ModuleField{
				Name: "email",
			},
			&types.ModuleField{
				Name:  "options",
				Multi: true,
			},
			&types.ModuleField{
				Name: "description",
			},
			&types.ModuleField{
				Name: "another_record",
				Kind: "Record",
			},
		},
	}

	// set up a module
	var err error
	module, err = Module().With(context.Background()).Create(module)
	assert(t, err == nil, "Error when creating module: %+v", err)
	assert(t, module.ID > 0, "Expected auto generated ID")

	record1 := &types.Record{
		ModuleID: module.ID,
	}

	record2 := &types.Record{
		ModuleID: module.ID,
		Values: types.RecordValueSet{
			&types.RecordValue{
				Name:  "name",
				Value: "John Doe",
			},
			&types.RecordValue{
				Name:  "email",
				Value: "john.doe@example.com",
			},
			&types.RecordValue{
				Name:  "options",
				Value: "1",
			},
			&types.RecordValue{
				Name:  "options",
				Value: "2",
			},
			&types.RecordValue{
				Name:  "options",
				Value: "3",
			},
			&types.RecordValue{
				Name:  "description",
				Value: "just an example",
			},
			&types.RecordValue{
				Name:  "another_record",
				Value: "918273645",
			},
		},
	}

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

		// create record
		m2, err := repository.Create(record2)
		assert(t, err == nil, "Error when creating record: %+v", err)
		assert(t, m2.ID > 0, "Expected auto generated ID")

		// fetch created record
		{
			ms, err := repository.FindByID(m1.ID)
			assert(t, err == nil, "Error when retrieving record by id: %+v", err)
			assert(t, ms.ID == m1.ID, "Expected ID from database to match, %d != %d", m1.ID, ms.ID)
			assert(t, ms.ModuleID == m1.ModuleID, "Expected Module ID from database to match, %d != %d", m1.ModuleID, ms.ModuleID)
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
			mr, err := repository.Find(module.ID, "", "id desc", 0, 20)
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Sort == "id desc", "Expected Meta.Sort == id desc, got '%s'", mr.Meta.Sort)
			assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)
			assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be descending")
		}

		// fetch all records
		{
			mr, err := repository.Find(module.ID, "", "name asc, email desc", 0, 20)
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Sort == "name asc, email desc", "Expected Meta.Sort == 'name asc, email desc' '%s'", mr.Meta.Sort)
			assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)

			// @todo sort is not stable
			// assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be ascending")
		}

		// fetch all records
		{
			mr, err := repository.Find(module.ID, "", "created_at desc", 0, 20)
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Sort == "created_at desc", "Expected Meta.Sort == created_at desc, got '%s'", mr.Meta.Sort)
			assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)

			// @todo sort is not stable
			// assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be ascending")
		}

		// fetch all records by query
		{
			filter := "name='John Doe' AND email='john.doe@example.com'"
			sort := "id desc"

			mr, err := repository.Find(module.ID, filter, sort, 0, 20)
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 1, "Expected one record, got %d", len(mr.Records))
			assert(t, mr.Meta.Count == 1, "Expected Meta.Count == 1, got %d", mr.Meta.Count)
			assert(t, mr.Meta.Page == 0, "Expected Meta.Page == 0, got %d", mr.Meta.Page)
			assert(t, mr.Meta.PerPage == 20, "Expected Meta.PerPage == 20, got %d", mr.Meta.PerPage)
			assert(t, mr.Meta.Filter == filter, "Expected Meta.Filter == %q, got %q", filter, mr.Meta.Filter)
			assert(t, mr.Meta.Sort == sort, "Expected Meta.Sort == %q, got %q", sort, mr.Meta.Sort)
		}

		// fetch all records by query
		{
			mr, err := repository.Find(module.ID, "name='niall'", "id asc", 0, 20)
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
			mr, err := repository.Find(module.ID, "", "", 0, 20)
			assert(t, err == nil, "Error when retrieving records: %+v", err)
			assert(t, len(mr.Records) == 0, "Expected no record, got %d", len(mr.Records))
		}
	}
}

func TestValueSanitizer(t *testing.T) {
	var (
		svc    = record{}
		module = &types.Module{
			Fields: types.ModuleFieldSet{
				&types.ModuleField{Name: "single1"},
				&types.ModuleField{Name: "multi1", Multi: true},
				&types.ModuleField{Name: "ref1", Kind: "Record"},
				&types.ModuleField{Name: "multiRef1", Kind: "Record", Multi: true},
			},
		}
		rvs types.RecordValueSet
	)

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}}
	test.NoError(t, svc.sanitizeValues(module, rvs), "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(rvs) == 1, "expecting 1 record value after sanitization, got %d", len(rvs))

	rvs = types.RecordValueSet{{Name: "unknown", Value: "single"}}
	test.Assert(t, svc.sanitizeValues(module, rvs) != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}, {Name: "single1", Value: "single2"}}
	test.Assert(t, svc.sanitizeValues(module, rvs) != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "multi1", Value: "multi1"}, {Name: "multi1", Value: "multi1"}}
	test.NoError(t, svc.sanitizeValues(module, rvs), "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(rvs) == 2, "expecting 2 record values after sanitization, got %d", len(rvs))
	test.Assert(t, rvs[0].Place == 0, "expecting first value to have place value 0, got %d", rvs[0].Place)
	test.Assert(t, rvs[1].Place == 1, "expecting second value to have place value 1, got %d", rvs[1].Place)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "multi1"}}
	test.Assert(t, svc.sanitizeValues(module, rvs) != nil, "expecting sanitizeValues() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "ref1", Value: "12345"}}
	test.NoError(t, svc.sanitizeValues(module, rvs), "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(rvs) == 1, "expecting 1 record values after sanitization, got %d", len(rvs))
	test.Assert(t, rvs[0].Ref == 12345, "expecting parsed ref value to match, got %d", rvs[0].Ref)

	rvs = types.RecordValueSet{{Name: "multiRef1", Value: "12345"}, {Name: "multiRef1", Value: "67890"}}
	test.NoError(t, svc.sanitizeValues(module, rvs), "unexpected error for sanitizeValues() call: %v")
	test.Assert(t, len(rvs) == 2, "expecting 2 record values after sanitization, got %d", len(rvs))
	test.Assert(t, rvs[0].Ref == 12345, "expecting parsed ref value to match, got %d", rvs[0].Ref)
	test.Assert(t, rvs[1].Ref == 67890, "expecting parsed ref value to match, got %d", rvs[1].Ref)
}
