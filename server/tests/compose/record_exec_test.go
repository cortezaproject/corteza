package compose

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/crusttech/human/server/compose/dalutils"
	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/tests/helpers"
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

// A card dropped onto a position another record already holds must land in
// front of it. The moved record is written to that position before the
// records it displaced are renumbered, so a sweep that includes the moved
// record leaves the two tied and the iterator picks the winner by ID — and an
// older neighbour wins, leaving the card one slot behind where it was dropped.
func TestRecordExecOrganizeDropsOntoAnOccupiedPosition(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read", "update")

	module := h.repoMakeRecordModuleWithFields(
		"organize drop position",
		&types.ModuleField{Name: "position", Kind: "Number"},
		&types.ModuleField{Name: "handle"},
	)

	makeRecord := func(position int, handle string) *types.Record {
		return h.makeRecord(module,
			&types.RecordValue{Name: "position", Value: strconv.Itoa(position)},
			&types.RecordValue{Name: "handle", Value: handle},
		)
	}

	order := func() string {
		sorting, _ := filter.NewSorting("position ASC")
		set, _, err := dalutils.ComposeRecordsList(context.Background(), defDal, module, types.RecordFilter{
			ModuleID:    module.ID,
			NamespaceID: module.NamespaceID,
			Sorting:     sorting,
		})
		h.noError(err)

		out := ""
		_ = set.Walk(func(r *types.Record) error {
			out += r.Values.FilterByName("handle")[0].Value
			return nil
		})
		return out
	}

	_ = makeRecord(1, "a")
	_ = makeRecord(2, "b")
	cRec := makeRecord(3, "c")

	h.a.Equal("abc", order())

	// 'c' is dropped in front of 'b', onto the position 'b' holds. 'c' is the
	// younger of the two, so this is the tie the sweep used to decide wrongly.
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{Name: "recordID", Value: strconv.FormatUint(cRec.ID, 10)},
		{Name: "positionField", Value: "position"},
		{Name: "position", Value: "2"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.Equal("acb", order())
}

// A drop is an update to the record, so the automation an author bound to
// record update has to hear it: organize used to dispatch only its own events,
// leaving a card moved on a board invisible to every update handler.
func TestRecordExecOrganizeFiresRecordUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read", "update")

	module := h.repoMakeRecordModuleWithFields(
		"organize update events",
		&types.ModuleField{Name: "position", Kind: "Number"},
		&types.ModuleField{Name: "handle"},
	)

	makeRecord := func(position int, handle string) *types.Record {
		return h.makeRecord(module,
			&types.RecordValue{Name: "position", Value: strconv.Itoa(position)},
			&types.RecordValue{Name: "handle", Value: handle},
		)
	}

	_ = makeRecord(1, "a")
	_ = makeRecord(2, "b")
	cRec := makeRecord(3, "c")

	var seen []string
	watch := func(event string) uintptr {
		return eventBus.Register(
			func(ctx context.Context, ev eventbus.Event) error {
				rec := ev.(interface{ Record() *types.Record }).Record()
				seen = append(seen, event+":"+rec.Values.FilterByName("handle")[0].Value)
				return nil
			},
			eventbus.For("compose:record"),
			eventbus.On(event),
		)
	}
	defer eventBus.Unregister(watch("beforeUpdate"), watch("afterUpdate"))

	// 'c' is dropped in front of 'b', which pushes 'b' to a new position too.
	h.apiSendRecordExec(module.NamespaceID, module.ID, "organize", request.ProcedureArgs{
		{Name: "recordID", Value: strconv.FormatUint(cRec.ID, 10)},
		{Name: "positionField", Value: "position"},
		{Name: "position", Value: "2"}}).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Both halves of the update contract, and only for the record that was
	// dragged — 'b' was renumbered by the same call and must stay quiet, or one
	// move on a full column becomes a run per card.
	h.a.Equal([]string{"beforeUpdate:c", "afterUpdate:c"}, seen)
}

// A beforeUpdate handler gates a save, so it gates a drop: the record must be
// left where it was when one refuses.
func TestRecordExecOrganizeRefusedByBeforeUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read", "update")

	module := h.repoMakeRecordModuleWithFields(
		"organize update veto",
		&types.ModuleField{Name: "position", Kind: "Number"},
		&types.ModuleField{Name: "handle"},
	)

	makeRecord := func(position int, handle string) *types.Record {
		return h.makeRecord(module,
			&types.RecordValue{Name: "position", Value: strconv.Itoa(position)},
			&types.RecordValue{Name: "handle", Value: handle},
		)
	}

	_ = makeRecord(1, "a")
	cRec := makeRecord(2, "c")

	ptr := eventBus.Register(
		func(ctx context.Context, ev eventbus.Event) error {
			return fmt.Errorf("not on my board")
		},
		eventbus.For("compose:record"),
		eventbus.On("beforeUpdate"),
	)
	defer eventBus.Unregister(ptr)

	payload, err := json.Marshal(request.RecordExec{Args: request.ProcedureArgs{
		{Name: "recordID", Value: strconv.FormatUint(cRec.ID, 10)},
		{Name: "positionField", Value: "position"},
		{Name: "position", Value: "0"},
	}})
	h.noError(err)

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/exec/organize", module.NamespaceID, module.ID)).
		Header("Accept", "application/json").
		JSON(string(payload)).
		Expect(h.t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not on my board")).
		End()

	lRec := h.lookupRecordByID(module, cRec.ID)
	h.a.Equal("2", lRec.Values.FilterByName("position")[0].Value)
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
