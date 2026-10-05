package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/crusttech/human/server/compose/dalutils"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/tests/helpers"
)

// Multi-line text and quotes survive the import as they are; a number with a
// decimal comma fails the row instead of being stored as 0.
func TestRecordImportCellValues(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import cell values module",
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "num", Kind: "Number", Options: types.ModuleFieldOptions{"precision": 2}},
	)
	helpers.AllowMe(h, module.RbacResource(), "records.search", "record.create")

	url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)

	importRows := func(csv string) (failed uint64) {
		ses := &rImportSession{}
		h.apiInitRecordImport(h.apiInit(), url, "records.csv", []byte(csv)).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(ses)

		h.apiRunRecordImport(h.apiInit(), url+"/"+ses.Response.SessionID, `{"fields":{"fname":"name","fnum":"num"},"onError":"fail","multiValueDelimiter":";"}`).
			End()

		var progress struct {
			Response struct {
				Progress struct {
					FinishedAt *time.Time `json:"finishedAt"`
					Failed     uint64     `json:"failed"`
				} `json:"progress"`
			} `json:"response"`
		}

		for i := 0; i < 50 && progress.Response.Progress.FinishedAt == nil; i++ {
			time.Sleep(100 * time.Millisecond)
			h.apiInit().
				Get(url+"/"+ses.Response.SessionID).
				Header("Accept", "application/json").
				Expect(t).
				Status(http.StatusOK).
				End().
				JSON(&progress)
		}

		h.a.NotNil(progress.Response.Progress.FinishedAt, "import must finish")
		return progress.Response.Progress.Failed
	}

	records := func() types.RecordSet {
		set, _, err := dalutils.ComposeRecordsList(context.Background(), defDal, module, types.RecordFilter{
			ModuleID:    module.ID,
			NamespaceID: module.NamespaceID,
		})
		h.noError(err)
		return set
	}

	h.a.Zero(importRows("fname,fnum\n\"line 1\nline 2\",1.5\n\"5\"\" pipe\",2\n"))

	set := records()
	h.a.Len(set, 2)
	names := []string{set[0].Values.Get("name", 0).Value, set[1].Values.Get("name", 0).Value}
	h.a.ElementsMatch([]string{"line 1\nline 2", `5" pipe`}, names)

	h.a.NotZero(importRows("fname,fnum\ncomma,\"1,5\"\n"))
	h.a.Len(records(), 2, "row with a decimal comma must not be imported")
}
