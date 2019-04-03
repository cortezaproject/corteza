// +build integration

package service

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/test"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestTrigger(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing", true)

	user := &systemTypes.User{
		ID:       1337,
		Name:     "John Crm Doe",
		Username: "johndoe",
		SatosaID: "12345",
	}

	// Set Identity (required for permission checks).
	ctx = auth.SetIdentityToContext(ctx, user)

	svc := Trigger().With(ctx)

	// the trigger object we're working with
	trigger := &types.Trigger{
		Name:     "Test",
		ModuleID: 123,
	}

	{
		{
			m, err := svc.Update(trigger)
			test.Assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			test.Assert(t, err != nil, "Expected error when updating invalid content")
		}

		// create trigger
		m, err := svc.Create(trigger)
		test.Assert(t, err == nil, "Error when creating trigger: %+v", err)
		test.Assert(t, m.ID > 0, "Expected auto generated ID")

		{
			_, err := svc.Create(trigger)
			test.Assert(t, err == nil, "Unexpected error when creating trigger, %+v", err)
		}

		// fetch created trigger
		{
			ms, err := svc.FindByID(m.ID)
			test.Assert(t, err == nil, "Error when retrieving trigger by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "Expected ID from database to match, %+v", errors.Errorf("%d != %d", m.ID, ms.ID))
			test.Assert(t, ms.Name == m.Name, "Expected Name from database to match, %+v", errors.Errorf("%s != %s", m.Name, ms.Name))
		}

		// update created trigger
		{
			m.Name = "Updated test"
			_, err := svc.Update(m)
			test.Assert(t, err == nil, "Error when updating trigger, %+v", err)
		}

		// re-fetch trigger
		{
			ms, err := svc.FindByID(m.ID)
			test.Assert(t, err == nil, "Error when retrieving trigger by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "re-fetch: Expected ID from database to match, %d != %d", m.ID, ms.ID)
			test.Assert(t, ms.Name == m.Name, "Expected Name from database to match, %s != %s", m.Name, ms.Name)
		}

		// delete trigger
		{
			err := svc.DeleteByID(m.ID)
			test.Assert(t, err == nil, "Error when deleting trigger by id: %+v", err)
		}
	}

}
