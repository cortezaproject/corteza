package compose

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearCharts() {
	h.clearNamespaces()
	h.noError(store.TruncateComposeCharts(context.Background(), service.DefaultStore))
}

func (h helper) makeChart(ns *types.Namespace, name string) *types.Chart {
	res := &types.Chart{
		ID:          id.Next(),
		CreatedAt:   time.Now(),
		Name:        name,
		Handle:      name,
		NamespaceID: ns.ID,
	}

	h.noError(store.CreateComposeChart(context.Background(), service.DefaultStore, res))
	return res
}

func (h helper) lookupChartByID(ID uint64) *types.Chart {
	res, err := store.LookupComposeChartByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestChartRead(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeChart(ns, "some-chart")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, m.Name)).
		Assert(jsonpath.Equal(`$.response.chartID`, fmt.Sprintf("%d", m.ID))).
		End()
}

func TestChartReadByHandle(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	c := h.makeChart(ns, "some-chart")

	cbh, err := service.DefaultChart.With(h.secCtx()).FindByHandle(ns.ID, c.Handle)

	h.noError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
}

func TestChartList(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.makeChart(ns, "chart1")
	h.makeChart(ns, "chart2")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/chart/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestChartList_filterForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.makeChart(ns, "chart")
	f := h.makeChart(ns, "chart_forbidden")

	h.deny(types.ChartPermissionResource.AppendID(f.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/chart/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.name=="chart_forbiden"]`)).
		End()
}

func TestChartCreateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/", ns.ID)).
		FormData("name", "some-chart").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create charts")).
		End()
}

func TestChartCreate(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "chart.create")

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/", ns.ID)).
		FormData("name", "some-chart").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestChartUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeChart(ns, "some-chart")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this chart")).
		End()
}

func TestChartUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	res := h.makeChart(ns, "some-chart")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, res.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupChartByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(res.Name, "changed-name")
}

func TestChartDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeChart(ns, "some-chart")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this chart")).
		End()
}

func TestChartDelete(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "delete")

	ns := h.makeNamespace("some-namespace")
	res := h.makeChart(ns, "some-chart")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, res.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupChartByID(res.ID)
	h.a.NotNil(res.DeletedAt)
}
