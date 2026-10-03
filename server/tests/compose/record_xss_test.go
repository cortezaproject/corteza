package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
)

// A string field may opt out of XSS sanitization so that values like
// "Name <mail@example.tld>" are kept as entered; rich text fields never can.
func TestRecordStringField_xssSanitizationOptOut(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record xss module",
		&types.ModuleField{Name: "plain", Kind: "String", Options: types.ModuleFieldOptions{"sanitizeXSS": false}},
		&types.ModuleField{Name: "safe", Kind: "String"},
		&types.ModuleField{Name: "rich", Kind: "String", Options: types.ModuleFieldOptions{"sanitizeXSS": false, "useRichTextEditor": true}},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")

	const (
		address = "Application No-Reply <noreply@example.com>"
		markup  = `<span onerror=alert()>Title here</span>`
	)

	var created struct {
		Response *types.Record
	}

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "plain", "value": %q}, {"name": "safe", "value": %q}, {"name": "rich", "value": %q}]}`, address, address, markup)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&created)

	h.a.NotNil(created.Response)

	// values are also passed through sanitization when read
	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, created.Response.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Contains(`$.response.values[?(@.name=="plain")].value`, address)).
		Assert(jsonpath.Contains(`$.response.values[?(@.name=="safe")].value`, "Application No-Reply ")).
		Assert(jsonpath.Contains(`$.response.values[?(@.name=="rich")].value`, "<span>Title here</span>")).
		End()
}
