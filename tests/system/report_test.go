package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) createReport(report *types.Report) *types.Report {
	if report.ID == 0 {
		report.ID = id.Next()
	}

	if report.CreatedAt.IsZero() {
		report.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateReport(context.Background(), report))
	return report
}

func (h helper) clearReports() {
	h.noError(store.TruncateReports(context.Background(), service.DefaultStore))
}

func (h helper) lookupReportByHandle(handle string) *types.Report {
	res, err := store.LookupReportByHandle(context.Background(), service.DefaultStore, handle)
	h.noError(err)
	return res
}

func TestReportScenarios_create(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ComponentRbacResource(), "report.create")

	h.apiInit().
		Post("/reports/").
		Header("Accept", "application/json").
		JSON(`{
			"handle": "test_report",
			"scenarios": [{ "label": "scenario 1", "filters": { "ds_1": { "raw": "field == 'value'" } } }]
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupReportByHandle("test_report")
	h.a.NotNil(r)

	h.a.Len(r.Scenarios, 1)
	sc := r.Scenarios[0]
	h.a.NotNil(sc.Filters)
	h.a.Contains(sc.Filters, "ds_1")
}

func TestReportScenarios_update(t *testing.T) {
	h := newHelper(t)
	h.clearReports()
	helpers.AllowMe(h, types.ReportRbacResource(0), "update")

	r := h.createReport(&types.Report{
		Handle: "test_report",
	})

	h.apiInit().
		Put(fmt.Sprintf("/reports/%d", r.ID)).
		Header("Accept", "application/json").
		JSON(`{
			"handle": "test_report",
			"scenarios": [{ "label": "scenario 1", "filters": { "ds_1": { "raw": "field == 'value'" } } }]
		}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r = h.lookupReportByHandle("test_report")
	h.a.NotNil(r)

	h.a.Len(r.Scenarios, 1)
	sc := r.Scenarios[0]
	h.a.NotNil(sc.Filters)
	h.a.Contains(sc.Filters, "ds_1")
}
