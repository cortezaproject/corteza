package compose

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/crusttech/human/server/compose/dalutils"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) softDeleteRecord(module *types.Module, rec *types.Record) {
	now := time.Now()
	rec.DeletedAt = &now
	h.noError(dalutils.ComposeRecordUpdate(context.Background(), defDal, module, rec))
}

func TestRecordListByID(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record by id module")
	helpers.AllowMe(h, module.RbacResource(), "records.search")

	a := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "a"})
	h.makeRecord(module, &types.RecordValue{Name: "name", Value: "b"})
	c := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "c"})
	h.softDeleteRecord(module, c)

	url := fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)

	t.Run("only the named records", func(t *testing.T) {
		h.apiInit().
			Get(url).
			Query("recordID[]", strconv.FormatUint(a.ID, 10)).
			Query("recordID[]", strconv.FormatUint(c.ID, 10)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Len(`$.response.set`, 1)).
			Assert(jsonpath.Equal(`$.response.set[0].recordID`, strconv.FormatUint(a.ID, 10))).
			End()
	})

	t.Run("deleted ones when asked for", func(t *testing.T) {
		h.apiInit().
			Get(url).
			Query("recordID[]", strconv.FormatUint(a.ID, 10)).
			Query("recordID[]", strconv.FormatUint(c.ID, 10)).
			Query("deleted", "1").
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Len(`$.response.set`, 2)).
			End()
	})

	t.Run("nothing for an ID that is not a number", func(t *testing.T) {
		h.apiInit().
			Get(url).
			Query("recordID[]", "not-an-id").
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Len(`$.response.set`, 0)).
			End()
	})
}

func TestRecordPatchByID(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record patch by id module")
	helpers.AllowMe(h, module.RbacResource(), "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	a := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "a"})
	b := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "b"})
	c := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "c"})

	h.apiInit().
		Patch(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "email", "value": "bulk@x.y"}], "recordID": ["%d", "%d"]}`, a.ID, b.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	for _, r := range []*types.Record{a, b} {
		v := h.lookupRecordByID(module, r.ID).Values.Get("email", 0)
		h.a.NotNil(v, "record %d must be patched", r.ID)
		h.a.Equal("bulk@x.y", v.Value)
	}

	h.a.Nil(h.lookupRecordByID(module, c.ID).Values.Get("email", 0), "a record not named must stay as it was")
}

func TestRecordBulkDeleteAndUndeleteByID(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record bulk by id module")
	helpers.AllowMe(h, module.RbacResource(), "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "delete", "undelete")

	a := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "a"})
	b := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "b"})
	c := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "c"})

	url := fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)
	body := fmt.Sprintf(`{"recordID": ["%d", "%d"]}`, a.ID, b.ID)

	h.apiInit().
		Delete(url).
		JSON(body).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotNil(h.lookupRecordByID(module, a.ID).DeletedAt)
	h.a.NotNil(h.lookupRecordByID(module, b.ID).DeletedAt)
	h.a.Nil(h.lookupRecordByID(module, c.ID).DeletedAt, "a record not named must not be deleted")

	h.apiInit().
		Patch(url+"undelete").
		JSON(fmt.Sprintf(`{"recordID": ["%d"]}`, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.Nil(h.lookupRecordByID(module, a.ID).DeletedAt)
	h.a.NotNil(h.lookupRecordByID(module, b.ID).DeletedAt, "a record not named must stay deleted")
}

func TestRecordExportResolvesRefs(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	ns := h.makeNamespace("record export refs namespace")
	company := h.makeRecordModuleWithFieldsOnNs("company", ns, &types.ModuleField{Name: "name", Kind: "String"})
	contact := h.makeRecordModuleWithFieldsOnNs("contact", ns,
		&types.ModuleField{Name: "name", Kind: "String"},
		&types.ModuleField{Name: "company", Kind: "Record", Options: types.ModuleFieldOptions{"moduleID": strconv.FormatUint(company.ID, 10), "labelField": "name"}},
		&types.ModuleField{Name: "tags", Kind: "Record", Multi: true, Options: types.ModuleFieldOptions{"moduleID": strconv.FormatUint(company.ID, 10), "labelField": "name"}},
	)
	helpers.AllowMe(h, contact.RbacResource(), "records.search")

	acme := h.makeRecord(company, &types.RecordValue{Name: "name", Value: "Acme"})
	birch := h.makeRecord(company, &types.RecordValue{Name: "name", Value: "Birch"})

	// tags name the later-created company first, so a label placed by the
	// lookup's result order instead of by ID lands against the wrong value
	h.makeRecord(contact,
		&types.RecordValue{Name: "name", Value: "Carol"},
		&types.RecordValue{Name: "company", Value: strconv.FormatUint(birch.ID, 10)},
		&types.RecordValue{Name: "tags", Value: strconv.FormatUint(birch.ID, 10), Place: 0},
		&types.RecordValue{Name: "tags", Value: strconv.FormatUint(acme.ID, 10), Place: 1},
	)

	rsp := h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/export.csv", ns.ID, contact.ID)).
		Query("fields", "name,company,tags").
		Query("resolveRefs", "true").
		Expect(t).
		Status(http.StatusOK).
		End()

	raw, err := io.ReadAll(rsp.Response.Body)
	h.noError(err)
	out := string(raw)

	h.a.Contains(out, "ID,name,company value,tags value,company,tags\n")
	h.a.Contains(out, ",Carol,Birch,Birch;Acme,", "each label must stand against its own value")
}

func TestRecordImportUpdatesExistingByID(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import by id module")
	helpers.AllowMe(h, module.RbacResource(), "records.search", "record.create")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	existing := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "old"})

	url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	csv := fmt.Sprintf("id,name\n%d,new\n", existing.ID)

	ses := &rImportSession{}
	h.apiInitRecordImport(h.apiInit(), url, "records.csv", []byte(csv)).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(ses)

	h.apiRunRecordImport(h.apiInit(), url+"/"+ses.Response.SessionID, `{"fields":{"id":"id","name":"name"},"onError":"fail"}`).
		Assert(helpers.AssertNoErrors).
		End()

	var progress struct {
		Response struct {
			Progress struct {
				FinishedAt *time.Time `json:"finishedAt"`
				Failed     uint64     `json:"failed"`
				FailReason string     `json:"failReason"`
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
	h.a.Zero(progress.Response.Progress.Failed, progress.Response.Progress.FailReason)

	set, _, err := dalutils.ComposeRecordsList(context.Background(), defDal, module, types.RecordFilter{
		ModuleID:    module.ID,
		NamespaceID: module.NamespaceID,
	})
	h.noError(err)
	h.a.Len(set, 1, "a row naming an existing record must update it, not add another")
	h.a.Equal("new", set[0].Values.Get("name", 0).Value)
}

func TestRecordExportByID(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record export by id module")
	helpers.AllowMe(h, module.RbacResource(), "records.search")

	a := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "a"})
	h.makeRecord(module, &types.RecordValue{Name: "name", Value: "b"})
	c := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "c"})

	export := func(filter string, ids ...string) string {
		req := h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/export.csv", module.NamespaceID, module.ID)).
			Query("fields", "name")
		if filter != "" {
			req = req.Query("filter", filter)
		}
		for _, id := range ids {
			req = req.Query("recordID[]", id)
		}

		raw, err := io.ReadAll(req.Expect(t).Status(http.StatusOK).End().Response.Body)
		h.noError(err)
		return string(raw)
	}

	h.a.Equal(
		fmt.Sprintf("ID,name\n%d,a\n%d,c\n", a.ID, c.ID),
		export("", strconv.FormatUint(a.ID, 10), strconv.FormatUint(c.ID, 10)),
		"only the named records are exported",
	)

	h.a.Equal(
		export("name = 'no such record'"),
		export("", "not-an-id"),
		"an ID that names nothing exports what an empty result does, never every record",
	)
}
