package seeder

import (
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/seeder"
	"testing"
)

func TestCreateUser(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	limit := 10
	gen := seeder.Seeder(h.ctx, seeder.DefaultStore, seeder.Faker())

	userIDs, err := gen.CreateUser(seeder.Params{Limit: limit})
	h.noError(err)
	h.a.Equal(limit, len(userIDs))
}

func TestClearAllUser(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	limit := 10
	gen := seeder.Seeder(h.ctx, seeder.DefaultStore, seeder.Faker())

	userIDs, err := gen.CreateUser(seeder.Params{Limit: limit})
	h.noError(err)
	h.a.Equal(limit, len(userIDs))

	err = gen.DeleteAllUser()
	h.noError(err)
}

func TestCreateRecord(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()
	h.clearModules()

	n := h.makeNamespace("fake-data-namespace")
	m := h.makeModule(n, "fake-data-module",
		setModuleField("String", "str1", true),
		setModuleField("String", "str2", false),
		setModuleField("Number", "number1", true),
		setModuleField("DateTime", "dt1", true),
	)

	limit := 10
	gen := seeder.Seeder(h.ctx, seeder.DefaultStore, seeder.Faker())

	recIDs, err := gen.CreateRecord(seeder.RecordParams{
		NamespaceID:     n.ID,
		NamespaceHandle: n.Slug,
		ModuleID:        m.ID,
		ModuleHandle:    m.Handle,
		Params:          seeder.Params{Limit: limit},
	})
	h.noError(err)
	h.a.Equal(limit, len(recIDs))
}

func TestClearAllRecord(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	n := h.makeNamespace("fake-data-namespace")
	m := h.makeModule(n, "fake-data-module",
		setModuleField("String", "str1", true),
		setModuleField("String", "str2", false),
		setModuleField("Number", "number1", true),
		setModuleField("DateTime", "dt1", true),
	)

	limit := 10
	gen := seeder.Seeder(h.ctx, seeder.DefaultStore, seeder.Faker())

	recIDs, err := gen.CreateRecord(seeder.RecordParams{
		NamespaceID:     n.ID,
		NamespaceHandle: n.Slug,
		ModuleID:        m.ID,
		ModuleHandle:    m.Handle,
		Params:          seeder.Params{Limit: limit},
	})
	h.noError(err)
	h.a.Equal(limit, len(recIDs))

	err = gen.DeleteAllRecord(m)
	h.noError(err)
}

func setModuleField(kind, name string, required bool) *types.ModuleField {
	return &types.ModuleField{Kind: kind, Name: name, Required: required}
}
