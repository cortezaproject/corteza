package compose

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
)

// An import whose rows carry no record IDs must create the rows without
// listing the whole module first (it used to run an unfiltered record
// search to find rows to update).
func TestRecordImportRun_withoutRecordIDs(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import run module",
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "email", Kind: "Email"},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create", "records.search")

	// rows the import must not touch
	h.makeRecord(module, &types.RecordValue{Name: "name", Value: "existing"})

	url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	api := h.apiInit()

	rsp := &rImportSession{}
	h.apiInitRecordImport(api, url, "rows.csv", []byte("fname,femail\nv1,v1@example.tld\nv2,v2@example.tld\n")).
		End().
		JSON(rsp)
	h.a.NotEmpty(rsp.Response.SessionID)

	h.apiRunRecordImport(api, fmt.Sprintf("%s/%s", url, rsp.Response.SessionID), `{"fields":{"fname":"name","femail":"email"},"onError":"fail"}`).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.progress")).
		End()

	set, _, err := dalutils.ComposeRecordsList(context.Background(), defDal, module, types.RecordFilter{
		ModuleID:    module.ID,
		NamespaceID: module.NamespaceID,
	})
	h.a.NoError(err)
	h.a.Len(set, 3, "two imported rows next to the existing one")
}

// Multi-line text and quotes survive the import as they are; a number with a
// decimal comma fails the row instead of being stored as 0.
func TestRecordImportRun_cellValues(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import cell values module",
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "num", Kind: "Number", Options: types.ModuleFieldOptions{"precision": 2}},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create", "records.search")

	url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	api := h.apiInit()

	importRows := func(csv string) {
		rsp := &rImportSession{}
		h.apiInitRecordImport(api, url, "rows.csv", []byte(csv)).End().JSON(rsp)
		h.a.NotEmpty(rsp.Response.SessionID)

		h.apiRunRecordImport(api, fmt.Sprintf("%s/%s", url, rsp.Response.SessionID), `{"fields":{"fname":"name","fnum":"num"},"onError":"fail","multiValueDelimiter":";"}`).
			End()
	}

	records := func() types.RecordSet {
		set, _, err := dalutils.ComposeRecordsList(context.Background(), defDal, module, types.RecordFilter{
			ModuleID:    module.ID,
			NamespaceID: module.NamespaceID,
		})
		h.a.NoError(err)
		return set
	}

	importRows("fname,fnum\n\"line 1\nline 2\",1.5\n\"5\"\" pipe\",2\n")

	set := records()
	h.a.Len(set, 2)
	names := []string{set[0].Values.Get("name", 0).Value, set[1].Values.Get("name", 0).Value}
	h.a.ElementsMatch([]string{"line 1\nline 2", `5" pipe`}, names)

	importRows("fname,fnum\ncomma,\"1,5\"\n")
	h.a.Len(records(), 2, "row with a decimal comma must not be imported")
}
