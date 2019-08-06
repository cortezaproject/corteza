// +build integration

package service

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestChart(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing", true)

	// Set Identity (required for permission checks).
	ctx = auth.SetIdentityToContext(ctx, auth.NewIdentity(1337))

	ns1, _ := createTestNamespaces(ctx, t)

	svc := Chart().With(ctx)

	// the chart object we're working with
	chart := &types.Chart{
		NamespaceID: ns1.ID,
		Name:        "Test",
	}

	{
		{
			m, err := svc.Update(chart)
			test.Assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			test.Assert(t, err != nil, "Expected error when updating invalid content")
		}

		// create chart
		m, err := svc.Create(chart)
		test.Assert(t, err == nil, "Error when creating chart: %+v", err)
		test.Assert(t, m.ID > 0, "Expected auto generated ID")

		{
			_, err := svc.Create(chart)
			test.Assert(t, err == nil, "Unexpected error when creating chart, %+v", err)
		}

		// fetch created chart
		{
			ms, err := svc.FindByID(m.NamespaceID, m.ID)
			test.Assert(t, err == nil, "Error when retrieving chart by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "Expected ID from database to match, %+v", errors.Errorf("%d != %d", m.ID, ms.ID))
			test.Assert(t, ms.Name == m.Name, "Expected Name from database to match, %+v", errors.Errorf("%s != %s", m.Name, ms.Name))
		}

		// update created chart
		{
			m.UpdatedAt = nil
			m.Name = "Updated test"
			_, err := svc.Update(m)
			test.Assert(t, err == nil, "Error when updating chart, %+v", err)
		}

		// re-fetch chart
		{
			ms, err := svc.FindByID(m.NamespaceID, m.ID)
			test.Assert(t, err == nil, "Error when retrieving chart by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "re-fetch: Expected ID from database to match, %d != %d", m.ID, ms.ID)
			test.Assert(t, ms.Name == m.Name, "Expected Name from database to match, %s != %s", m.Name, ms.Name)
		}

		// delete chart
		{
			err := svc.DeleteByID(m.NamespaceID, m.ID)
			test.Assert(t, err == nil, "Error when deleting chart by id: %+v", err)
		}
	}

}
