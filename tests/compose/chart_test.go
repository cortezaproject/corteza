package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) repoChart() repository.ChartRepository {
	return repository.Chart(context.Background(), db())
}

func (h helper) repoMakeChart(ns *types.Namespace, name string) *types.Chart {
	m, err := h.
		repoChart().
		Create(&types.Chart{Name: name, Handle: name, NamespaceID: ns.ID})
	h.a.NoError(err)

	return m
}

func TestChartRead(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeChart(ns, "some-chart")

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

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	c := h.repoMakeChart(ns, "some-chart")

	cbh, err := service.DefaultChart.With(h.secCtx()).FindByHandle(ns.ID, c.Handle)

	h.a.NoError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
}

func TestChartList(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")

	h.repoMakeChart(ns, "app")
	h.repoMakeChart(ns, "app")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/chart/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestChartCreateForbidden(t *testing.T) {
	h := newHelper(t)

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/", ns.ID)).
		FormData("name", "some-chart").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoCreatePermissions")).
		End()
}

func TestChartCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "chart.create")

	ns := h.repoMakeNamespace("some-namespace")

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
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeChart(ns, "some-chart")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoUpdatePermissions")).
		End()
}

func TestChartUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeChart(ns, "some-chart")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m, err := h.repoChart().FindByID(ns.ID, m.ID)
	h.a.NoError(err)
	h.a.NotNil(m)
	h.a.Equal(m.Name, "changed-name")
}

func TestChartDeleteForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeChart(ns, "some-chart")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoDeletePermissions")).
		End()
}

func TestChartDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "read")
	h.allow(types.ChartPermissionResource.AppendWildcard(), "delete")

	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeChart(ns, "some-chart")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/chart/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m, err := h.repoChart().FindByID(ns.ID, m.ID)
	h.a.Error(err, "compose.repository.ChartNotFound")
}
