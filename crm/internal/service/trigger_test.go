package service

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/crm/types"
)

func TestTrigger(t *testing.T) {
	repository := Trigger().With(context.Background())

	// the trigger object we're working with
	trigger := &types.Trigger{
		Name:     "Test",
		ModuleID: 123,
	}

	{
		{
			m, err := repository.Update(trigger)
			assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			assert(t, err != nil, "Expected error when updating invalid content")
		}

		// create trigger
		m, err := repository.Create(trigger)
		assert(t, err == nil, "Error when creating trigger: %+v", err)
		assert(t, m.ID > 0, "Expected auto generated ID")

		{
			_, err := repository.Create(trigger)
			assert(t, err == nil, "Unexpected error when creating trigger, %+v", err)
		}

		// fetch created trigger
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving trigger by id: %+v", err)
			assert(t, ms.ID == m.ID, "Expected ID from database to match, %+v", errors.Errorf("%d != %d", m.ID, ms.ID))
			assert(t, ms.Name == m.Name, "Expected Name from database to match, %+v", errors.Errorf("%s != %s", m.Name, ms.Name))
		}

		// update created trigger
		{
			m.Name = "Updated test"
			_, err := repository.Update(m)
			assert(t, err == nil, "Error when updating trigger, %+v", err)
		}

		// re-fetch trigger
		{
			ms, err := repository.FindByID(m.ID)
			assert(t, err == nil, "Error when retrieving trigger by id: %+v", err)
			assert(t, ms.ID == m.ID, "re-fetch: Expected ID from database to match, %d != %d", m.ID, ms.ID)
			assert(t, ms.Name == m.Name, "Expected Name from database to match, %s != %s", m.Name, ms.Name)
		}

		// delete trigger
		{
			err := repository.DeleteByID(m.ID)
			assert(t, err == nil, "Error when deleting trigger by id: %+v", err)
		}
	}

}
