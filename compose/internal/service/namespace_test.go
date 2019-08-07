// +build integration

package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestNamespace(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing", true)

	// Set Identity (required for permission checks).
	ctx = auth.SetIdentityToContext(ctx, auth.NewIdentity(1337))

	svc := Namespace().With(ctx)

	// the namespace object we're working with
	namespace := &types.Namespace{
		Name: "Test",
	}

	prevNamespaceCount := uint(0)

	{
		{
			m, err := svc.Update(namespace)
			test.Assert(t, m == nil, "Expected empty return for invalid update, got %#v", m)
			test.Assert(t, err != nil, "Expected error when updating invalid content")
		}

		// create namespace
		m, err := svc.Create(namespace)
		test.Assert(t, err == nil, "Error when creating namespace: %+v", err)
		test.Assert(t, m.ID > 0, "Expected auto generated ID")

		// fetch created namespace
		{
			ms, err := svc.FindByID(m.ID)
			test.Assert(t, err == nil, "Error when retrieving namespace by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			test.Assert(t, ms.Name == m.Name, "Expected Name from database to match, %s != %s", m.Name, ms.Name)
		}

		// update created namespace
		{
			m.Name = "Updated test"
			m.UpdatedAt = nil
			_, err := svc.Update(m)
			test.Assert(t, err == nil, "Error when updating namespace, %+v", err)
		}

		// re-fetch namespace
		{
			ms, err := svc.FindByID(m.ID)
			test.Assert(t, err == nil, "Error when retrieving namespace by id: %+v", err)
			test.Assert(t, ms.ID == m.ID, "Expected ID from database to match, %d != %d", m.ID, ms.ID)
			test.Assert(t, ms.Name == m.Name, "Expected Name from database to match, %s != %s", m.Name, ms.Name)
		}

		// fetch all namespaces
		{
			_, f, err := svc.Find(types.NamespaceFilter{})
			test.Assert(t, err == nil, "Error when retrieving namespaces: %+v", err)
			test.Assert(t, f.Count > 0, "Expected at least one namespace, got %d", f.Count)
			prevNamespaceCount = f.Count
		}

		// re-fetch namespace
		{
			err := svc.DeleteByID(m.ID)
			test.Assert(t, err == nil, "Error when deleting namespace by id: %+v", err)
		}

		// fetch all namespaces
		{
			_, f, err := svc.Find(types.NamespaceFilter{})
			test.Assert(t, err == nil, "Error when retrieving namespaces: %+v", err)
			test.Assert(t, f.Count < prevNamespaceCount, "Expected namespaces count to decrease after deletion, %d < %d", f.Count, prevNamespaceCount)
		}
	}
}
