package compose

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ChartRBACResource.AppendWildcard(), "read")
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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ChartRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	c := h.makeChart(ns, "some-chart")

	cbh, err := service.DefaultChart.FindByHandle(h.secCtx(), ns.ID, c.Handle)

	h.noError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
}

func TestChartList(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.makeChart(ns, "chart")
	f := h.makeChart(ns, "chart_forbidden")

	h.deny(types.ChartRBACResource.AppendID(f.ID), "read")

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
		Header("Accept", "application/json").
		FormData("name", "some-chart").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create charts")).
		End()
}

func TestChartCreate(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.NamespaceRBACResource.AppendWildcard(), "chart.create")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeChart(ns, "some-chart")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this chart")).
		End()
}

func TestChartUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	res := h.makeChart(ns, "some-chart")
	h.allow(types.ChartRBACResource.AppendWildcard(), "update")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ChartRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeChart(ns, "some-chart")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this chart")).
		End()
}

func TestChartDelete(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ChartRBACResource.AppendWildcard(), "read")
	h.allow(types.ChartRBACResource.AppendWildcard(), "delete")

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

func TestChartLabels(t *testing.T) {
	h := newHelper(t)
	h.clearCharts()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.NamespaceRBACResource.AppendWildcard(), "chart.create")
	h.allow(types.ChartRBACResource.AppendWildcard(), "read")
	h.allow(types.ChartRBACResource.AppendWildcard(), "update")
	h.allow(types.ChartRBACResource.AppendWildcard(), "delete")

	var (
		ns = h.makeNamespace("some-namespace")
		ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Chart{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/chart/", ns.ID),
			types.Chart{Labels: map[string]string{"foo": "bar", "bar": "42"}},
			payload,
		)
		req.NotZero(payload.ID)

		h.a.Equal(payload.Labels["foo"], "bar",
			"labels must contain foo with value bar")
		h.a.Equal(payload.Labels["bar"], "42",
			"labels must contain bar with value 42")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")

		ID = payload.ID
	})

	t.Run("update", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req     = require.New(t)
			payload = &types.Chart{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, ID),
			types.Chart{Labels: map[string]string{"foo": "baz", "baz": "123"}},
			payload,
		)
		req.NotZero(payload.ID)
		req.Nil(payload.UpdatedAt, "updatedAt must not change after changing labels")

		req.Equal(payload.Labels["foo"], "baz",
			"labels must contain foo with value baz")
		req.NotContains(payload.Labels, "bar",
			"labels must not contain bar")
		req.Equal(payload.Labels["baz"], "123",
			"labels must contain baz with value 123")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")
	})

	t.Run("search", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req = require.New(t)
			set = types.ChartSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/chart/", ns.ID),
			&set, url.Values{"labels": []string{"baz=123"}},
		)
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
