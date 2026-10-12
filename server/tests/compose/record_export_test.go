package compose

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

// A failed export must not reach the browser as an attachment, or the
// browser shows its own error page instead of the error message.
func TestRecordExport_failureIsNotAnAttachment(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record export module",
		&types.ModuleField{Name: "name", Kind: "String"},
	)
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
