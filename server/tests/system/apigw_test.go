package system

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) createRouteWithFilter(s string, fkind string) (*types.ApigwRoute, *types.ApigwFilter) {
	return h.createRouteAndFilterWithEnabled(s, fkind, true, true)
}

func (h helper) createRouteWithFilterEnabled(s string, fkind string, enable bool) (*types.ApigwRoute, *types.ApigwFilter) {
	return h.createRouteAndFilterWithEnabled(s, fkind, true, enable)
}

func (h helper) createRouteAndFilterWithEnabled(s string, fkind string, rEnable, fEnable bool) (*types.ApigwRoute, *types.ApigwFilter) {
	r := h.createRoute(&types.ApigwRoute{Endpoint: "/" + s, Method: "GET", Enabled: rEnable})
	f := h.createFilters(&types.ApigwFilter{Kind: fkind, Enabled: fEnable}, r.ID)

	return r, f
}

func (h helper) createRoute(route *types.ApigwRoute) *types.ApigwRoute {
	if route.ID == 0 {
		route.ID = id.Next()
	}

	if route.CreatedAt.IsZero() {
		route.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateApigwRoute(context.Background(), route))
	return route
}

func (h helper) createFilters(f *types.ApigwFilter, routeID uint64) *types.ApigwFilter {
	f.Route = routeID

	if f.ID == 0 {
		f.ID = id.Next()
	}

	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateApigwFilter(context.Background(), f))
	return f
}

func (h helper) clearRoutes() {
	h.noError(store.TruncateApigwFilters(context.Background(), service.DefaultStore))
	h.noError(store.TruncateApigwRoutes(context.Background(), service.DefaultStore))
}

func TestApigwRouteRead(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test", "")

	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/route/%d", r.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.routeID`, strconv.FormatUint(r.ID, 10))).
		Assert(jsonpath.Equal(`$.response.endpoint`, "/test")).
		End()
}

func TestApigwRouteRead_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test", "")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/route/%d", r.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToRead")).
		End()
}

func TestApigwRouteSearch(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	h.createRouteWithFilter("test1", "")
	h.createRouteWithFilter("test2", "")

	helpers.AllowMe(h, types.ComponentRbacResource(), "apigw-routes.search")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/route/")).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		Assert(jsonpath.Equal(`$.response.set[0].endpoint`, "/test1")).
		Assert(jsonpath.Equal(`$.response.set[1].endpoint`, "/test2")).
		End()
}

func TestApigwRouteSearch_includeDisabled(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	h.createRouteWithFilter("test1", "")
	h.createRouteWithFilter("test2", "")
	h.createRouteAndFilterWithEnabled("test3", "", false, true)
	h.createRouteAndFilterWithEnabled("test4", "", false, false)

	helpers.AllowMe(h, types.ComponentRbacResource(), "apigw-routes.search")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/route/")).
		Query("disabled", "1").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 4)).
		Assert(jsonpath.Equal(`$.response.set[0].endpoint`, "/test1")).
		Assert(jsonpath.Equal(`$.response.set[1].endpoint`, "/test2")).
		End()
}

func TestApigwRouteSearch_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	h.createRouteWithFilter("test1", "")
	h.createRouteWithFilter("test2", "")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/route/")).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToSearch")).
		End()
}

func TestApigwRouteSearch_forbidenSpecific(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test1", "")
	h.createRouteWithFilter("test2", "")

	helpers.AllowMe(h, types.ComponentRbacResource(), "apigw-routes.search")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(0), "read")
	helpers.DenyMe(h, types.ApigwRouteRbacResource(r.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/route/")).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 1)).
		Assert(jsonpath.Equal(`$.response.set[0].endpoint`, "/test2")).
		End()
}

func TestApigwRouteCreate(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	helpers.AllowMe(h, types.ComponentRbacResource(), "apigw-route.create")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/route")).
		Header("Accept", "application/json").
		FormData("endpoint", "/test").
		FormData("method", "GET").
		FormData("enabled", "false").
		FormData("group", "g1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.routeID`)).
		Assert(jsonpath.Equal(`$.response.endpoint`, "/test")).
		Assert(jsonpath.Equal(`$.response.enabled`, false)).
		End()
}

func TestApigwRouteCreate_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	h.apiInit().
		Post(fmt.Sprintf("/apigw/route")).
		Header("Accept", "application/json").
		FormData("endpoint", "/test").
		FormData("method", "GET").
		FormData("enabled", "false").
		FormData("group", "g1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToCreate")).
		End()
}

func TestApigwRouteUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test", "")

	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "update")

	h.apiInit().
		Put(fmt.Sprintf("/apigw/route/%d", r.ID)).
		Header("Accept", "application/json").
		FormData("endpoint", "/test-edited").
		FormData("enabled", "false").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.routeID`)).
		Assert(jsonpath.Equal(`$.response.endpoint`, "/test-edited")).
		Assert(jsonpath.Equal(`$.response.enabled`, false)).
		End()
}

func TestApigwRouteUpdate_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test", "")

	h.apiInit().
		Put(fmt.Sprintf("/apigw/route/%d", r.ID)).
		Header("Accept", "application/json").
		FormData("endpoint", "/test-edited").
		FormData("enabled", "false").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToUpdate")).
		End()
}

func TestApigwRouteDelete(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test", "")

	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/apigw/route/%d", r.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApigwRouteDelete_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test", "")

	h.apiInit().
		Delete(fmt.Sprintf("/apigw/route/%d", r.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToDelete")).
		End()
}

func TestApigwRouteUnDelete(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test", "")

	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "delete")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/route/%d/undelete", r.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApigwRouteUnDelete_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test", "")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/route/%d/undelete", r.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToUndelete")).
		End()
}

// Filters

func TestApigwFilterRead(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilter("test", "test")

	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/filter/%d", f.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.filterID`, strconv.FormatUint(f.ID, 10))).
		Assert(jsonpath.Equal(`$.response.routeID`, strconv.FormatUint(r.ID, 10))).
		End()
}

func TestApigwFilterRead_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	_, f := h.createRouteWithFilter("test", "test")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/filter/%d", f.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToRead")).
		End()
}

func TestApigwFilterSearch(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilter("test1", "")

	helpers.AllowMe(h, types.ComponentRbacResource(), "apigw-routes.search")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/filter/")).
		Query("routeID", strconv.FormatUint(r.ID, 10)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 1)).
		Assert(jsonpath.Equal(`$.response.set[0].filterID`, strconv.FormatUint(f.ID, 10))).
		Assert(jsonpath.Equal(`$.response.set[0].routeID`, strconv.FormatUint(r.ID, 10))).
		End()
}

func TestApigwFilterSearch_includeDisabled(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilterEnabled("test1", "", false)

	helpers.AllowMe(h, types.ComponentRbacResource(), "apigw-routes.search")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/filter/")).
		Query("routeID", strconv.FormatUint(r.ID, 10)).
		Query("disabled", "1").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 1)).
		Assert(jsonpath.Equal(`$.response.set[0].filterID`, strconv.FormatUint(f.ID, 10))).
		Assert(jsonpath.Equal(`$.response.set[0].routeID`, strconv.FormatUint(r.ID, 10))).
		End()
}

func TestApigwFilterSearch_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test1", "")

	h.apiInit().
		Get(fmt.Sprintf("/apigw/filter/")).
		Query("routeID", strconv.FormatUint(r.ID, 10)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToRead")).
		End()
}

func TestApigwFilterCreate(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "update")

	h.apiInit().
		Put(fmt.Sprintf("/apigw/filter")).
		Header("Accept", "application/json").
		FormData("routeID", strconv.FormatUint(r.ID, 10)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.filterID`)).
		Assert(jsonpath.Equal(`$.response.routeID`, strconv.FormatUint(r.ID, 10))).
		End()
}

func TestApigwFilterCreate_enabled(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "update")

	h.apiInit().
		Put(fmt.Sprintf("/apigw/filter")).
		Header("Accept", "application/json").
		FormData("routeID", strconv.FormatUint(r.ID, 10)).
		FormData("enabled", strconv.FormatBool(true)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.filterID`)).
		Assert(jsonpath.Equal(`$.response.routeID`, strconv.FormatUint(r.ID, 10))).
		Assert(jsonpath.Equal(`$.response.enabled`, true)).
		End()
}

func TestApigwFilterCreate_disabled(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "update")

	h.apiInit().
		Put(fmt.Sprintf("/apigw/filter")).
		Header("Accept", "application/json").
		FormData("routeID", strconv.FormatUint(r.ID, 10)).
		FormData("enabled", strconv.FormatBool(false)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.filterID`)).
		Assert(jsonpath.Equal(`$.response.routeID`, strconv.FormatUint(r.ID, 10))).
		Assert(jsonpath.NotPresent(`$.response.enabled`)).
		End()
}

func TestApigwFilterCreate_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, _ := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")

	h.apiInit().
		Put(fmt.Sprintf("/apigw/filter")).
		Header("Accept", "application/json").
		FormData("routeID", strconv.FormatUint(r.ID, 10)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToUpdate")).
		End()
}

func TestApigwFilterUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "update")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/filter/%d", f.ID)).
		Header("Accept", "application/json").
		FormData("routeID", strconv.FormatUint(r.ID, 10)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.filterID`, strconv.FormatUint(f.ID, 10))).
		Assert(jsonpath.Equal(`$.response.routeID`, strconv.FormatUint(r.ID, 10))).
		End()
}

func TestApigwFilterUpdate_enabled(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilterEnabled("test1", "", false)
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "update")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/filter/%d", f.ID)).
		Header("Accept", "application/json").
		FormData("routeID", strconv.FormatUint(r.ID, 10)).
		FormData("enabled", strconv.FormatBool(true)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.filterID`, strconv.FormatUint(f.ID, 10))).
		Assert(jsonpath.Equal(`$.response.routeID`, strconv.FormatUint(r.ID, 10))).
		Assert(jsonpath.Equal(`$.response.enabled`, true)).
		End()
}

func TestApigwFilterUpdate_disabled(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilterEnabled("test1", "", true)
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "update")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/filter/%d", f.ID)).
		Header("Accept", "application/json").
		FormData("routeID", strconv.FormatUint(r.ID, 10)).
		FormData("enabled", strconv.FormatBool(false)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.filterID`, strconv.FormatUint(f.ID, 10))).
		Assert(jsonpath.Equal(`$.response.routeID`, strconv.FormatUint(r.ID, 10))).
		Assert(jsonpath.NotPresent(`$.response.enabled`)).
		End()
}

func TestApigwFilterUpdate_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/filter/%d", f.ID)).
		Header("Accept", "application/json").
		FormData("routeID", strconv.FormatUint(r.ID, 10)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToUpdate")).
		End()
}

func TestApigwFilterDelete(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/apigw/filter/%d", f.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApigwFilterDelete_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	_, f := h.createRouteWithFilter("test1", "")

	h.apiInit().
		Delete(fmt.Sprintf("/apigw/filter/%d", f.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToDelete")).
		End()
}

func TestApigwFilterUnDelete(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "delete")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/filter/%d/undelete", f.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApigwFilterUnDelete_forbiden(t *testing.T) {
	h := newHelper(t)
	h.clearRoutes()

	r, f := h.createRouteWithFilter("test1", "")
	helpers.AllowMe(h, types.ApigwRouteRbacResource(r.ID), "read")

	h.apiInit().
		Post(fmt.Sprintf("/apigw/filter/%d/undelete", f.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("apigwRoute.errors.notAllowedToUndelete")).
		End()
}
