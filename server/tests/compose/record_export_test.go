package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
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
