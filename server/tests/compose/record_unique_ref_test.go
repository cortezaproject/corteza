package compose

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// Saving a record with a value in a unique Record (reference) field used to
// hit an unimplemented lookup and panic; uniqueness of references is left to
// duplicate detection, so the save has to go through.
func TestRecordCreate_uniqueReferenceField(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	mRef := h.repoMakeRecordModuleWithFields("unique ref target module")
	rRef := h.makeRecord(mRef, &types.RecordValue{Name: "name", Value: "target"})

	module := h.makeRecordModuleWithFieldsOnNs("unique ref module", h.lookupNamespaceByID(mRef.NamespaceID),
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "ref", Kind: "Record", Options: types.ModuleFieldOptions{
			"moduleID": strconv.FormatUint(mRef.ID, 10),
			"isUnique": true,
		}},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")

	var created struct {
		Response *types.Record
	}

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "first"}, {"name": "ref", "value": "%d"}]}`, rRef.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&created)

	h.a.NotNil(created.Response)
	r := h.lookupRecordByID(module, created.Response.ID)
	h.a.Equal(strconv.FormatUint(rRef.ID, 10), r.Values.Get("ref", 0).Value)
}
