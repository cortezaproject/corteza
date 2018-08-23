package repository

import (
	"context"
	"testing"
)

func TestField(t *testing.T) {

	repository := NewField(context.TODO()).With(context.Background())

	{
		// fetch all fields
		{
			ms, err := repository.Find()
			must(t, err, "Error when retrieving fields")
			assert(t, len(ms) > 1, "Expected more than one field")
		}

		// fetch named field
		{
			m, err := repository.FindByType("email")
			must(t, err, "Error when retrieving field by name")
			assert(t, m != nil, "Unexpected nil value for field by name")
			assert(t, m.Type == "email", "Unexpected type, expected email, got %s", m.Type)
		}
	}
}
