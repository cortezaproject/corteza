package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// A required boolean field must be saveable as false; false is sent as an
// empty value and used to be rejected as a missing required field.
func TestRecordRequiredBool_falseIsAccepted(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("required bool module",
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "flag", Kind: "Bool", Required: true},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	url := fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)

	// create with the flag sent the way the webapp sends false: name only
	var created struct {
		Response *types.Record
	}
	h.apiInit().
		Post(url).
		JSON(`{"values": [{"name": "name", "value": "val"}, {"name": "flag"}]}`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&created)

	h.a.NotNil(created.Response)
	h.a.NotZero(created.Response.ID)

	// update, keeping the flag false and explicitly empty
	h.apiInit().
		Post(fmt.Sprintf("%s%d", url, created.Response.ID)).
		JSON(`{"values": [{"name": "name", "value": "changed"}, {"name": "flag", "value": ""}]}`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// a required field of another kind is still enforced
	other := h.makeRecordModuleWithFieldsOnNs("required string module", h.lookupNamespaceByID(module.NamespaceID),
		&types.ModuleField{Name: "req", Kind: "String", Required: true},
	)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", other.NamespaceID, other.ID)).
		JSON(`{"values": [{"name": "req", "value": ""}]}`).
		Header("Accept", "application/json").
		Expect(t).
		Assert(helpers.AssertRecordValueError(
			&types.RecordValueError{
				Kind:    "empty",
				Message: "record-field.errors.empty",
				Meta:    map[string]interface{}{"field": "req"},
			},
		)).
		End()
}
