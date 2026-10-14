package compose

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/compose/service"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// Bulk edit sends an empty value for a number field that should be cleared;
// that must not end up stored as 0 (or [0] for multi-value fields).
func TestRecordPatch_emptyNumberIsCleared(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record patch module",
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "amount", Kind: "Number"},
		&types.ModuleField{Name: "nums", Kind: "Number", Multi: true},
	)
	record := h.makeRecord(module,
		&types.RecordValue{Name: "name", Value: "keep me"},
		&types.RecordValue{Name: "amount", Value: "42"},
		&types.RecordValue{Name: "nums", Value: "1", Place: 0},
		&types.RecordValue{Name: "nums", Value: "2", Place: 1},
	)

	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	h.apiInit().
		Patch(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"query": "recordID = %d", "values": [{"name": "amount", "value": ""}, {"name": "nums"}]}`, record.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r)

	h.a.Equal("keep me", r.Values.Get("name", 0).Value, "untouched fields stay")

	if v := r.Values.Get("amount", 0); v != nil {
		h.a.Equal("", v.Value, "cleared number must not become 0")
	}

	for _, v := range r.Values.FilterByName("nums") {
		h.a.NotEqual("0", v.Value, "cleared multi-value number must not become [0]")
		h.a.Equal("", v.Value)
	}
}

// A bulk patch by query reaches every matching record, however many there
// are, and ends even when the patched records still match the query
func TestRecordPatch_byQueryWalksAllPages(t *testing.T) {
	const items = 1100

	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record patch module",
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "amount", Kind: "Number"},
	)
	for i := 0; i < items; i++ {
		h.makeRecord(module, &types.RecordValue{Name: "name", Value: "a"})
	}

	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	ctx, cancel := context.WithTimeout(h.secCtx(), 20*time.Second)
	defer cancel()

	err := service.DefaultRecord.BulkModifyByFilter(ctx, types.RecordFilter{
		Query:       "name = 'a'",
		NamespaceID: module.NamespaceID,
		ModuleID:    module.ID,
	}, types.RecordValueSet{{Name: "amount", Value: "1"}}, types.OperationTypePatch)
	h.noError(err)

	set, _, err := service.DefaultRecord.Find(h.secCtx(), types.RecordFilter{
		Query:       "amount = 1",
		NamespaceID: module.NamespaceID,
		ModuleID:    module.ID,
		Paging:      filter.Paging{Limit: items + 1},
	})
	h.noError(err)
	h.a.Len(set, items, "every matching record is patched")
}
