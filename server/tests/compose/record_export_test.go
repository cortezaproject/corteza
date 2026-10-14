package compose

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// A failed export must not reach the browser as an attachment, or the
// browser shows its own error page instead of the error message.
func TestRecordExport_failureIsNotAnAttachment(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record export module",
		&types.ModuleField{Name: "name", Kind: "String"},
	)
	helpers.AllowMe(h, module.RbacResource(), "records.search")
	h.makeRecord(module, &types.RecordValue{Name: "name", Value: "a"})

	r := h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/exporttest.csv", module.NamespaceID, module.ID)).
		Query("fields", "name").
		Query("timezone", "Invalid/Zone").
		Expect(t).
		Status(http.StatusInternalServerError).
		End()

	h.a.Empty(r.Response.Header.Get("Content-Disposition"))
}

// A reference field without a module to resolve against exports its raw IDs.
func TestRecordExport_resolveRefsWithoutReferencedModule(t *testing.T) {
	for name, moduleID := range map[string]string{
		"missing module": strconv.FormatUint(id.Next(), 10),
		"no module set":  "",
	} {
		t.Run(name, func(t *testing.T) {
			h := newHelper(t)
			h.clearRecords()

			opts := types.ModuleFieldOptions{"labelField": "name"}
			if moduleID != "" {
				opts["moduleID"] = moduleID
			}

			module := h.repoMakeRecordModuleWithFields("record export module",
				&types.ModuleField{Name: "name", Kind: "String"},
				&types.ModuleField{Name: "ref", Kind: "Record", Options: opts},
			)
			helpers.AllowMe(h, module.RbacResource(), "records.search")
			refID := strconv.FormatUint(id.Next(), 10)
			rec := h.makeRecord(module,
				&types.RecordValue{Name: "name", Value: "a"},
				&types.RecordValue{Name: "ref", Value: refID},
			)

			h.apiInit().
				Get(fmt.Sprintf("/namespace/%d/module/%d/record/exporttest.csv", module.NamespaceID, module.ID)).
				Query("fields", "name,ref").
				Query("resolveRefs", "true").
				Expect(t).
				Status(http.StatusOK).
				Body(fmt.Sprintf("ID,name,ref\n%d,a,%s\n", rec.ID, refID)).
				End()
		})
	}
}

// Export holds only the records the user may read, with or without resolved
// references, and a referenced record the user may not read keeps its raw ID
func TestRecordExport_onlyReadableRecords(t *testing.T) {
	for _, resolve := range []string{"false", "true"} {
		t.Run("resolveRefs="+resolve, func(t *testing.T) {
			h := newHelper(t)
			h.clearRecords()

			target := h.repoMakeRecordModuleWithFields("record export target",
				&types.ModuleField{Name: "name", Kind: "String"},
			)
			hidden := h.makeRecord(target, &types.RecordValue{Name: "name", Value: "hidden target"})

			module := h.makeRecordModuleWithFieldsOnNs("record export module", h.lookupNamespaceByID(target.NamespaceID),
				&types.ModuleField{Name: "name", Kind: "String"},
				&types.ModuleField{Name: "ref", Kind: "Record", Options: types.ModuleFieldOptions{
					"moduleID":   strconv.FormatUint(target.ID, 10),
					"labelField": "name",
				}},
			)
			helpers.AllowMe(h, module.RbacResource(), "records.search")

			readable := h.makeRecord(module,
				&types.RecordValue{Name: "name", Value: "readable"},
				&types.RecordValue{Name: "ref", Value: strconv.FormatUint(hidden.ID, 10)},
			)
			h.makeRecord(module, &types.RecordValue{Name: "name", Value: "unreadable"})

			helpers.DenyMe(h, types.RecordRbacResource(0, 0, 0), "read")
			helpers.AllowMe(h, types.RecordRbacResource(module.NamespaceID, module.ID, readable.ID), "read")

			h.apiInit().
				Get(fmt.Sprintf("/namespace/%d/module/%d/record/exporttest.csv", module.NamespaceID, module.ID)).
				Query("fields", "name,ref").
				Query("resolveRefs", resolve).
				Expect(t).
				Status(http.StatusOK).
				Body(fmt.Sprintf("ID,name,ref\n%d,readable,%d\n", readable.ID, hidden.ID)).
				End()
		})
	}
}

// Export needs the same permission as searching the module's records
func TestRecordExport_needsRecordSearch(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record export module",
		&types.ModuleField{Name: "name", Kind: "String"},
	)
	h.makeRecord(module, &types.RecordValue{Name: "name", Value: "a"})
	helpers.DenyMe(h, module.RbacResource(), "records.search")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/exporttest.csv", module.NamespaceID, module.ID)).
		Query("fields", "name").
		Header("Accept", "application/json").
		Expect(t).
		Assert(helpers.AssertError("record.errors.notAllowedToSearch")).
		End()
}
