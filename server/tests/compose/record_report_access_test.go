package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/tests/helpers"
)

// Reports aggregate record values, fields that can not be read must not be used in them
func TestRecordReportUnreadableField(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields(
		"record testing module",
		&types.ModuleField{Name: "public", Kind: "Number"},
		&types.ModuleField{Name: "salary", Kind: "Number"},
	)

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.DenyMe(h, types.ModuleFieldRbacResource(0, 0, module.Fields.FindByName("salary").ID), "record.value.read")

	for name, q := range map[string][3]string{
		"metric":    {"SUM(salary) AS total", "created_at", ""},
		"dimension": {"COUNT(ID) AS total", "salary", ""},
		"filter":    {"COUNT(ID) AS total", "created_at", "salary > 1000"},
	} {
		t.Run(name, func(t *testing.T) {
			h.apiInit().
				Get(fmt.Sprintf("/namespace/%d/module/%d/record/report", module.NamespaceID, module.ID)).
				Query("metrics", q[0]).
				Query("dimensions", q[1]).
				Query("filter", q[2]).
				Header("Accept", "application/json").
				Expect(t).
				Status(http.StatusOK).
				Assert(helpers.AssertError("record.errors.notAllowedToRead")).
				End()
		})
	}
}
