// +build integration

package service

import (
	"context"
	"testing"

	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/test"
	systemService "github.com/crusttech/crust/system/service"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestRecord(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing", true)

	user := &systemTypes.User{
		ID:       1337,
		Username: "TestUser",
	}
	{
		err := user.GeneratePassword("Mary had a little lamb, little lamb, little lamb")
		test.Assert(t, err == nil, "Error generating password: %+v", err)
	}

	{
		userSvc := systemService.TestUser(t, ctx)
		_, err := userSvc.Create(user, nil, "")
		test.NoError(t, err, "expected no error creating user, got %v", err)
	}

	ctx = auth.SetIdentityToContext(ctx, auth.NewIdentity(user.Identity()))

	svc := Record().With(ctx)

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
	module, err = Module().With(ctx).Create(module)
	test.Assert(t, err == nil, "Error when creating module: %+v", err)
	test.Assert(t, module.ID > 0, "Expected auto generated ID")

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
			m, err := svc.Update(record1)
			test.Assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			test.Assert(t, err != nil, "Expected error when updating invalid record")
		}

		// create record
		m1, err := svc.Create(record1)
		test.Assert(t, err == nil, "Error when creating record: %+v", err)
		test.Assert(t, m1.ID > 0, "Expected auto generated ID")

		// create record
		m2, err := svc.Create(record2)
		test.Assert(t, err == nil, "Error when creating record: %+v", err)
		test.Assert(t, m2.ID > 0, "Expected auto generated ID")

		// fetch created record
		{
			ms, err := svc.FindByID(m1.ID)
			test.Assert(t, err == nil, "Error when retrieving record by id: %+v", err)
			test.Assert(t, ms.ID == m1.ID, "Expected ID from database to match, %d != %d", m1.ID, ms.ID)
			test.Assert(t, ms.ModuleID == m1.ModuleID, "Expected Module ID from database to match, %d != %d", m1.ModuleID, ms.ModuleID)
		}

		// update created record
		{
			_, err := svc.Update(m1)
			test.Assert(t, err == nil, "Error when updating record, %+v", err)
		}

		// re-fetch record
		{
			ms, err := svc.FindByID(m1.ID)
			test.Assert(t, err == nil, "Error when retrieving record by id: %+v", err)
			test.Assert(t, ms.ID == m1.ID, "Expected ID from database to match, %d != %d", m1.ID, ms.ID)
			test.Assert(t, ms.ModuleID == m1.ModuleID, "Expected ID from database to match, %d != %d", m1.ModuleID, ms.ModuleID)
		}

		// fetch all records
		{
			mr, err := svc.Find(module.ID, "", "id desc", 0, 20)
			test.Assert(t, err == nil, "Error when retrieving records: %+v", err)
			test.Assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			test.Assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			test.Assert(t, mr.Meta.Sort == "id desc", "Expected Meta.Sort == id desc, got '%s'", mr.Meta.Sort)
			test.Assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)
			test.Assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be descending")
		}

		// fetch all records
		{
			mr, err := svc.Find(module.ID, "", "name asc, email desc", 0, 20)
			test.Assert(t, err == nil, "Error when retrieving records: %+v", err)
			test.Assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			test.Assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			test.Assert(t, mr.Meta.Sort == "name asc, email desc", "Expected Meta.Sort == 'name asc, email desc' '%s'", mr.Meta.Sort)
			test.Assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)

			// @todo sort is not stable
			// test.Assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be ascending")
		}

		// fetch all records
		{
			mr, err := svc.Find(module.ID, "", "created_at desc", 0, 20)
			test.Assert(t, err == nil, "Error when retrieving records: %+v", err)
			test.Assert(t, len(mr.Records) == 2, "Expected two record, got %d", len(mr.Records))
			test.Assert(t, mr.Meta.Count == 2, "Expected Meta.Count == 2, got %d", mr.Meta.Count)
			test.Assert(t, mr.Meta.Sort == "created_at desc", "Expected Meta.Sort == created_at desc, got '%s'", mr.Meta.Sort)
			test.Assert(t, mr.Records[0].ModuleID == m1.ModuleID, "Expected record module to match, %d != %d", m1.ModuleID, mr.Records[0].ModuleID)

			// @todo sort is not stable
			// test.Assert(t, mr.Records[0].ID > mr.Records[1].ID, "Expected order to be ascending")
		}

		// fetch all records by query
		{
			filter := "name='John Doe' AND email='john.doe@example.com'"
			sort := "id desc"

			mr, err := svc.Find(module.ID, filter, sort, 0, 20)
			test.Assert(t, err == nil, "Error when retrieving records: %+v", err)
			test.Assert(t, len(mr.Records) == 1, "Expected one record, got %d", len(mr.Records))
			test.Assert(t, mr.Meta.Count == 1, "Expected Meta.Count == 1, got %d", mr.Meta.Count)
			test.Assert(t, mr.Meta.Page == 0, "Expected Meta.Page == 0, got %d", mr.Meta.Page)
			test.Assert(t, mr.Meta.PerPage == 20, "Expected Meta.PerPage == 20, got %d", mr.Meta.PerPage)
			test.Assert(t, mr.Meta.Filter == filter, "Expected Meta.Filter == %q, got %q", filter, mr.Meta.Filter)
			test.Assert(t, mr.Meta.Sort == sort, "Expected Meta.Sort == %q, got %q", sort, mr.Meta.Sort)
		}

		// fetch all records by query
		{
			mr, err := svc.Find(module.ID, "name='niall'", "id asc", 0, 20)
			test.Assert(t, err == nil, "Error when retrieving records: %+v", err)
			test.Assert(t, len(mr.Records) == 0, "Expected no records, got %d", len(mr.Records))
		}

		// delete record
		{
			err := svc.DeleteByID(m1.ID)
			test.Assert(t, err == nil, "Error when retrieving record by id: %+v", err)

			err = svc.DeleteByID(m2.ID)
			test.Assert(t, err == nil, "Error when retrieving record by id: %+v", err)
		}

		// fetch all records
		{
			mr, err := svc.Find(module.ID, "", "", 0, 20)
			test.Assert(t, err == nil, "Error when retrieving records: %+v", err)
			test.Assert(t, len(mr.Records) == 0, "Expected no record, got %d", len(mr.Records))
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
