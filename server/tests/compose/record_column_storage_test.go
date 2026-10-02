package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// replaceModel registers the module's model and applies the schema
// alterations it needs, the way the admin does through the UI
func (h helper) replaceModel(ns *types.Namespace, module *types.Module) {
	ctx := context.Background()

	models, err := service.ModulesToModelSet(defDal, ns, module)
	h.noError(err)
	h.a.Len(models, 1)

	alts, err := defDal.ReplaceModel(ctx, nil, models[0])
	h.noError(err)

	if len(alts) > 0 {
		errs, err := defDal.ApplyAlteration(ctx, alts...)
		h.noError(err)
		for _, e := range errs {
			h.noError(e)
		}

		alts, err = defDal.ReloadModel(ctx, nil, models[0])
		h.noError(err)
		h.a.Empty(alts, "alterations should be applied")
	}
}

func plainField(name, kind string) *types.ModuleField {
	f := &types.ModuleField{Name: name, Kind: kind}
	f.Config.DAL.EncodingStrategy = &types.EncodingStrategy{EncodingStrategyPlain: &types.EncodingStrategyPlain{}}
	return f
}

// A module stored in its own table with every field in its own column must
// still accept records, also when the fields were first stored as JSON and
// moved to columns afterwards (the "values" column then stays behind).
func TestRecordCreate_allFieldsInColumns(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "record.create")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read", "record.value.update")

	ns := h.makeNamespace("column storage namespace")

	create := func(t *testing.T, module *types.Module) {
		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
			JSON(`{"values": [{"name": "name", "value": "val"}, {"name": "amount", "value": "3"}]}`).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	}

	t.Run("fresh module", func(t *testing.T) {
		module := &types.Module{Name: "fresh columns", Handle: "fresh_columns", NamespaceID: ns.ID,
			Fields: types.ModuleFieldSet{plainField("name", "String"), plainField("amount", "Number")},
		}
		module.Config.DAL.Ident = "compose_record_fresh_columns"
		module = h.createModule(ns, module)
		h.replaceModel(ns, module)
		create(t, module)
	})

	t.Run("fields moved from JSON to columns", func(t *testing.T) {
		module := &types.Module{Name: "moved columns", Handle: "moved_columns", NamespaceID: ns.ID,
			Fields: types.ModuleFieldSet{{Name: "name", Kind: "String"}, {Name: "amount", Kind: "Number"}},
		}
		module.Config.DAL.Ident = "compose_record_moved_columns"
		module = h.createModule(ns, module)
		h.replaceModel(ns, module)

		for _, f := range module.Fields {
			f.Config.DAL.EncodingStrategy = &types.EncodingStrategy{EncodingStrategyPlain: &types.EncodingStrategyPlain{}}
		}
		h.noError(store.UpdateComposeModuleField(context.Background(), service.DefaultStore, module.Fields...))
		h.replaceModel(ns, module)

		create(t, module)
	})
}
