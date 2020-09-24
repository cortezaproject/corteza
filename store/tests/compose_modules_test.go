package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testComposeModules(t *testing.T, s store.ComposeModules) {
	var (
		ctx = context.Background()
		req = require.New(t)

		namespaceID = id.Next()

		makeNew = func(name, handle string) *types.Module {
			// minimum data set for new composeModule
			return &types.Module{
				ID:          id.Next(),
				NamespaceID: namespaceID,
				CreatedAt:   time.Now(),
				Name:        name,
				Handle:      handle,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Module) {
			req := require.New(t)
			req.NoError(s.TruncateComposeModules(ctx))
			res := makeNew(string(rand.Bytes(10)), string(rand.Bytes(10)))
			req.NoError(s.CreateComposeModule(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		composeModule := makeNew("ComposeModuleCRUD", "compose-module-crud")
		req.NoError(s.CreateComposeModule(ctx, composeModule))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup", func(t *testing.T) {
		t.Run("by ID", func(t *testing.T) {
			composeModule := makeNew("look up by id", "look-up-by-id")
			req.NoError(s.CreateComposeModule(ctx, composeModule))
			fetched, err := s.LookupComposeModuleByID(ctx, composeModule.ID)
			req.NoError(err)
			req.Equal(composeModule.Name, fetched.Name)
			req.Equal(composeModule.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})

		t.Run("by NamespaceID, Name", func(t *testing.T) {
			composeModule := makeNew("look up by namespaceIDName", "look-up-by-namespaceIDName")
			req.NoError(s.CreateComposeModule(ctx, composeModule))
			fetched, err := s.LookupComposeModuleByNamespaceIDName(ctx, composeModule.NamespaceID, composeModule.Name)
			req.NoError(err)
			req.Equal(composeModule.Name, fetched.Name)
			req.Equal(composeModule.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})

		t.Run("by Handle", func(t *testing.T) {
			composeModule := makeNew("look up by namespaceIDHandle", "look-up-by-namespaceIDHandle")
			req.NoError(s.CreateComposeModule(ctx, composeModule))
			fetched, err := s.LookupComposeModuleByNamespaceIDHandle(ctx, composeModule.NamespaceID, composeModule.Handle)
			req.NoError(err)
			req.Equal(composeModule.Name, fetched.Name)
			req.Equal(composeModule.ID, fetched.ID)
			req.NotNil(fetched.CreatedAt)
			req.Nil(fetched.UpdatedAt)
			req.Nil(fetched.DeletedAt)
		})
	})

	t.Run("update", func(t *testing.T) {
		req, composeModule := truncAndCreate(t)
		composeModule.Name = "ComposeModuleCRUD+2"

		req.NoError(s.UpdateComposeModule(ctx, composeModule))

		updated, err := s.LookupComposeModuleByID(ctx, composeModule.ID)
		req.NoError(err)
		req.Equal(composeModule.Name, updated.Name)
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, composeModule := truncAndCreate(t)
			composeModule.Name = "ComposeModuleCRUD+2"
	
			req.NoError(s.UpsertComposeModule(ctx, composeModule))
	
			updated, err := s.LookupComposeModuleByID(ctx, composeModule.ID)
			req.NoError(err)
			req.Equal(composeModule.Name, updated.Name)
		})

		t.Run("new", func(t *testing.T) {
			composeModule := makeNew("upsert me", "upsert-me")
			composeModule.Name = "ComposeChartCRUD+2"

			req.NoError(s.UpsertComposeModule(ctx, composeModule))
	
			upserted, err := s.LookupComposeModuleByID(ctx, composeModule.ID)
			req.NoError(err)
			req.Equal(composeModule.Name, upserted.Name)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Module", func(t *testing.T) {
			req, composeModule := truncAndCreate(t)
			req.NoError(s.DeleteComposeModule(ctx, composeModule))
			_, err := s.LookupComposeModuleByID(ctx, composeModule.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, composeModule := truncAndCreate(t)
			req.NoError(s.DeleteComposeModuleByID(ctx, composeModule.ID))
			_, err := s.LookupComposeModuleByID(ctx, composeModule.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.Module{
			makeNew("/one-one", "module-1-1"),
			makeNew("/one-two", "module-1-2"),
			makeNew("/two-one", "module-2-1"),
			makeNew("/two-two", "module-2-2"),
			makeNew("/two-deleted", "module-2-d"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateComposeModules(ctx))
		req.NoError(s.CreateComposeModule(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchComposeModules(ctx, types.ModuleFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchComposeModules(ctx, types.ModuleFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchComposeModules(ctx, types.ModuleFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, f, err = s.SearchComposeModules(ctx, types.ModuleFilter{Handle: "module-2-1"})
		req.NoError(err)
		req.Len(set, 1)

		// find all prefixed
		set, f, err = s.SearchComposeModules(ctx, types.ModuleFilter{Query: "/two-"})
		req.NoError(err)
		req.Len(set, 2)

		_ = f // dummy
	})
}
