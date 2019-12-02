package compose

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/steinfletcher/apitest"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) apiSendRecordExec(nsID, modID uint64, proc string, args []request.ProcedureArg) *apitest.Response {
	payload, err := json.Marshal(request.RecordExec{Args: args})
	h.a.NoError(err)

	return h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/exec/%s", nsID, modID, proc)).
		JSON(string(payload)).
		Expect(h.t)
}

func TestRecordExecUnknownProcedure(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/namespace/0/module/0/record/exec/test-unexisting-proc").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("unknown procedure")).
		End()
}

func TestRecordExec(t *testing.T) {
	h := newHelper(t)

	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.update")

	module := h.repoMakeRecordModuleWithFields(
		"record testing module",
		&types.ModuleField{Name: "position", Kind: "Number"},
		&types.ModuleField{Name: "handle"},
		&types.ModuleField{Name: "category"},
	)

	makeRecord := func(position int, handle, cat string) *types.Record {
		return h.repoMakeRecord(module,
			&types.RecordValue{Name: "position", Value: strconv.Itoa(position)},
			&types.RecordValue{Name: "handle", Value: handle},
			&types.RecordValue{Name: "category", Value: cat},
		)
	}

	assertSort := func(expectedHandles, expectedCats string) {
		// Using record service for fetching to avoid value pre-fetching etc..
		set, _, err := service.DefaultRecord.With(h.secCtx()).Find(types.RecordFilter{
			ModuleID:    module.ID,
			NamespaceID: module.NamespaceID,
			Sort:        "position ASC",
			PageFilter:  rh.PageFilter{},
		})

		h.a.NoError(err)
		h.a.NotNil(set)

		actualHandles := ""
		actualCats := ""

		_ = set.Walk(func(r *types.Record) error {
			v := r.Values.FilterByName("handle")

			if len(v) == 1 {
				actualHandles += v[0].Value
			} else {
				actualHandles += strconv.Itoa(len(v))
			}

			actualCats += r.Values.FilterByName("category")[0].Value[3:]

			return nil
		})

		h.a.Equal(expectedHandles, actualHandles)
		h.a.Equal(expectedCats, actualCats)
	}

	var (
		aRec = makeRecord(1, "a", "CAT1")
		bRec = makeRecord(2, "b", "CAT1")
		cRec = makeRecord(3, "c", "CAT1")
		dRec = makeRecord(4, "d", "CAT2")
		eRec = makeRecord(5, "e", "CAT2")
		fRec = makeRecord(6, "f", "CAT2")
		gRec = makeRecord(7, "g", "CAT3")
		hRec = makeRecord(8, "h", "CAT3")
		iRec = makeRecord(9, "i", "CAT3")
	)

	// map handle to record ID so we can use it for reordering
	rr := map[string]string{
		"a": strconv.FormatUint(aRec.ID, 10),
		"b": strconv.FormatUint(bRec.ID, 10),
		"c": strconv.FormatUint(cRec.ID, 10),
		"d": strconv.FormatUint(dRec.ID, 10),
		"e": strconv.FormatUint(eRec.ID, 10),
		"f": strconv.FormatUint(fRec.ID, 10),
		"g": strconv.FormatUint(gRec.ID, 10),
		"h": strconv.FormatUint(hRec.ID, 10),
		"i": strconv.FormatUint(iRec.ID, 10),
	}

	assertSort("abcdefghi", "111222333")

	// Move a to the middle
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{"recordID", rr["a"]},
		{"positionField", "position"},
		{"position", "5"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	assertSort("bcdeafghi", "112212333")

	// Move i to the beginning
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{"recordID", rr["i"]},
		{"positionField", "position"},
		{"position", "0"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	//                            bcdeafghi
	//                            v<------^
	assertSort("ibcdeafgh", "311221233")

	// Move b to the 5th place
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{"recordID", rr["b"]},
		{"filter", "category = 'CAT1'"},
		{"positionField", "position"},
		{"position", "5"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	//                            ibcdeafgh
	//                             ^->v
	assertSort("icdebfagh", "312212133")

	// This will keep order of letters but move b to category-2
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{"recordID", rr["b"]},
		{"groupField", "category"},
		{"group", "CAT2"},
		{"position", "5"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	//                            icdebfagh
	//                                ^
	assertSort("icdebfagh", "312222133")

	rsv, err := h.repoRecord().LoadValues([]string{"category"}, []uint64{bRec.ID})
	h.a.NoError(err)
	h.a.NotNil(rsv)
	h.a.Len(rsv.FilterByName("category"), 1)
	h.a.Equal("CAT2", rsv.FilterByName("category")[0].Value)
}
