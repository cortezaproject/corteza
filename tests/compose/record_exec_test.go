package compose

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/dalutils"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest"
)

func (h helper) apiSendRecordExec(nsID, modID uint64, proc string, args []request.ProcedureArg) *apitest.Response {
	payload, err := json.Marshal(request.RecordExec{Args: args})
	h.noError(err)

	return h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/exec/%s", nsID, modID, proc)).
		JSON(string(payload)).
		Expect(h.t)
}

func TestRecordExecUnknownProcedure(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/namespace/0/module/0/record/exec/test-unexisting-proc").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("unknown procedure")).
		End()
}

func TestRecordExecOrganize(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	//helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	//helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read", "update")
	//helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read", "record.value.update")

	module := h.repoMakeRecordModuleWithFields(
		"record testing module",
		&types.ModuleField{Name: "position", Kind: "Number"},
		&types.ModuleField{Name: "handle"},
		&types.ModuleField{Name: "category"},
	)

	makeRecord := func(position int, handle, cat string) *types.Record {
		return h.makeRecord(module,
			&types.RecordValue{Name: "position", Value: strconv.Itoa(position)},
			&types.RecordValue{Name: "handle", Value: handle},
			&types.RecordValue{Name: "category", Value: cat},
		)
	}

	assertSort := func(expectedHandles, expectedCats string) {
		// Using record service for fetching to avoid value pre-fetching etc..
		sorting, _ := filter.NewSorting("position ASC")
		set, _, err := dalutils.ComposeRecordsList(context.Background(), defDal, module, types.RecordFilter{
			ModuleID:    module.ID,
			NamespaceID: module.NamespaceID,
			Sorting:     sorting,
		})

		h.noError(err)
		h.a.NotNil(set)

		actualHandles := ""
		actualCats := ""

		_ = set.Walk(func(r *types.Record) error {
			//fmt.Printf("%d\t%s\t%s\t%s\n", r.ID,
			//	r.Values.FilterByName("position")[0].Value,
			//	r.Values.FilterByName("handle")[0].Value,
			//	r.Values.FilterByName("category")[0].Value,
			//)

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

	t.Logf("seeding records")

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

	t.Logf("testing initial order")
	assertSort("abcdefghi", "111222333")

	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **
	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **

	t.Logf("moving 'a' to position '6'")
	// Move a to the middle
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{Name: "recordID", Value: rr["a"]},
		{Name: "positionField", Value: "position"},
		{Name: "position", Value: "6"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	//                            abcdefghi
	//                            ^---v
	assertSort("bcdeafghi", "112212333")

	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **
	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **

	t.Logf("moving 'i' to position '0'")
	// Move i to the beginning
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{Name: "recordID", Value: rr["i"]},
		{Name: "positionField", Value: "position"},
		{Name: "position", Value: "0"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	//                            bcdeafghi
	//                            v<------^
	assertSort("ibcdeafgh", "311221233")

	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **
	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **

	t.Logf("moving 'b' to position '5'")
	// Move b to the 5th place
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{Name: "recordID", Value: rr["b"]},
		{Name: "filter", Value: "category = 'CAT1'"},
		{Name: "positionField", Value: "position"},
		{Name: "position", Value: "5"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	//                            ibcdeafgh
	//                             ^->v
	assertSort("icdebafgh", "312211233")

	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **
	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **

	t.Logf("moving 'b' to category CAT2")
	// This will keep order of letters but move b to category-2
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{Name: "recordID", Value: rr["b"]},
		{Name: "groupField", Value: "category"},
		{Name: "group", Value: "CAT2"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	//                            icdebfagh
	//                                ^
	assertSort("icdebafgh", "312221233")

	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **
	// ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **

	lRec := h.lookupRecordByID(module, bRec.ID)
	h.a.NotNil(lRec.Values)
	h.a.Len(lRec.Values.FilterByName("category"), 1)
	h.a.Equal("CAT2", lRec.Values.FilterByName("category")[0].Value)
}
