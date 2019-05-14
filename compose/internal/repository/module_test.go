// +build unit integration

package repository

import (
	"context"
	"testing"

	"github.com/crusttech/crust/compose/types"
	"github.com/crusttech/crust/internal/test"

	"github.com/titpetric/factory"
)

func TestModule_updateFields(t *testing.T) {
	tx(t, func(ctx context.Context, db *factory.DB, ns *types.Namespace) (err error) {
		var (
			m    *types.Module
			repo = Module(ctx, db)
		)

		m, err = repo.Create(&types.Module{
			NamespaceID: ns.ID,
			Name:        "test-module",
		})

		test.NoError(t, err, "unexpected error on module creation")
		test.Assert(t, len(m.Fields) == 0, "unexpected fields found in the fresh module")

		m, err = repo.Create(&types.Module{
			NamespaceID: ns.ID,
			Name:        "test-module",
			Fields: types.ModuleFieldSet{
				&types.ModuleField{Name: "one"},
				&types.ModuleField{Name: "two"},
			},
		})

		test.NoError(t, err, "unexpected error on module creation")
		test.Assert(t, len(m.Fields) == 2, "expecting to find two fields in the new module")

		m.Fields[0].Name = "one-v2"
		m.Fields[1] = &types.ModuleField{Name: "three"}
		m, err = repo.Update(m)

		test.NoError(t, err, "unexpected error on module update")
		test.Assert(t, len(m.Fields) == 2, "expecting to find two fields in the new module")
		test.Assert(t, m.Fields[0].Name == "one-v2", "expecting to find field 'one'")
		test.Assert(t, m.Fields[0].Place == 0, "expecting Place=0")
		test.Assert(t, m.Fields[1].Name == "three", "expecting to find field 'three'")
		test.Assert(t, m.Fields[1].Place == 1, "expecting Place=1")

		return
	})
}
