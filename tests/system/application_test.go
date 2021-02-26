package system

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
)

func (h helper) clearApplications() {
	h.noError(store.TruncateApplications(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeApplication(ss ...string) *types.Application {
	var res = &types.Application{
		ID:        id.Next(),
		CreatedAt: time.Now(),
		Unify:     &types.ApplicationUnify{},
	}

	if len(ss) > 0 {
		res.Name = ss[0]
	} else {
		res.Name = "n_" + rs()
	}

	h.a.NoError(store.CreateApplication(context.Background(), service.DefaultStore, res))

	return res
}

func (h helper) lookupApplicationByID(ID uint64) *types.Application {
	res, err := store.LookupApplicationByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func (h helper) lookupApplicationByName(name string) *types.Application {
	res, _, err := store.SearchApplications(context.Background(), service.DefaultStore, types.ApplicationFilter{
		Name: name,
	})
	h.noError(err)
	if len(res) == 0 {
		return nil
	}
	return res[0]
}

func TestApplicationRead(t *testing.T) {
	h := newHelper(t)

	u := h.repoMakeApplication()

	h.apiInit().
		Get(fmt.Sprintf("/application/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, u.Name)).
		Assert(jsonpath.Equal(`$.response.applicationID`, fmt.Sprintf("%d", u.ID))).
		End()
}

func TestApplicationList(t *testing.T) {
	h := newHelper(t)

	h.repoMakeApplication(h.randEmail())
	h.repoMakeApplication(h.randEmail())

	h.apiInit().
		Get("/application/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApplicationList_filterForbidden(t *testing.T) {
	h := newHelper(t)

	// @todo this can be a problematic test because it leaves
	//       behind applications that are not denied this context
	//       db purge might be needed

	h.repoMakeApplication("application")
	f := h.repoMakeApplication()

	h.deny(types.ApplicationRBACResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/application/").
		Query("name", f.Name).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(fmt.Sprintf(`$.response.set[? @.name=="%s"]`, f.Name))).
		End()
}

func TestApplicationCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/application/").
		Header("Accept", "application/json").
		FormData("name", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create applications")).
		End()
}

func TestApplicationCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemRBACResource, "application.create")

	h.apiInit().
		Post("/application/").
		FormData("name", rs()).
		FormData("handle", "handle_"+rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApplicationCreate_weight(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemRBACResource, "application.create")
	name := "name_weight_create_" + rs()

	h.apiInit().
		Post("/application/").
		FormData("name", name).
		FormData("weight", "10").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res := h.lookupApplicationByName(name)
	h.a.NotNil(res)
	h.a.Equal(10, res.Weight)
}

func TestApplicationUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeApplication()

	h.apiInit().
		Put(fmt.Sprintf("/application/%d", u.ID)).
		Header("Accept", "application/json").
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this application")).
		End()
}

func TestApplicationUpdate(t *testing.T) {
	h := newHelper(t)
	res := h.repoMakeApplication()
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "update")

	newName := "updated-" + rs()
	newHandle := "updated-" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/application/%d", res.ID)).
		Header("Accept", "application/json").
		FormData("name", newName).
		FormData("handle", newHandle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupApplicationByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(newName, res.Name)
}

func TestApplicationUpdate_weight(t *testing.T) {
	h := newHelper(t)
	res := h.repoMakeApplication()
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "update")

	newName := "updated-" + rs()
	newHandle := "updated-" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/application/%d", res.ID)).
		Header("Accept", "application/json").
		FormData("name", newName).
		FormData("handle", newHandle).
		FormData("weight", "20").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupApplicationByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(20, res.Weight)
}

func TestApplicationReorder_forbiden(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "update")
	a := h.repoMakeApplication()
	b := h.repoMakeApplication()
	c := h.repoMakeApplication()
	h.deny(types.ApplicationRBACResource.AppendID(b.ID), "update")

	h.apiInit().
		Post("/application/reorder").
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{ "applicationIDs": ["%d", "%d", "%d"] }`, b.ID, a.ID, c.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this application")).
		End()
}

func TestApplicationReorder(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "update")
	a := h.repoMakeApplication()
	b := h.repoMakeApplication()
	c := h.repoMakeApplication()

	h.apiInit().
		Post("/application/reorder").
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{ "applicationIDs": ["%d", "%d", "%d"] }`, b.ID, a.ID, c.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	check := func(app uint64, w int) {
		res := h.lookupApplicationByID(app)
		h.a.NotNil(res)
		h.a.Equal(w, res.Weight)
	}

	check(b.ID, 1)
	check(a.ID, 2)
	check(c.ID, 3)
}

func TestApplicationDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeApplication()

	h.apiInit().
		Delete(fmt.Sprintf("/application/%d", u.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this application")).
		End()
}

func TestApplicationDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "delete")

	res := h.repoMakeApplication()

	h.apiInit().
		Delete(fmt.Sprintf("/application/%d", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupApplicationByID(res.ID)
	h.a.NotNil(res)
	h.a.NotNil(res.DeletedAt)
}

func TestApplicationLabels(t *testing.T) {
	h := newHelper(t)
	h.clearApplications()

	h.allow(types.SystemRBACResource, "application.create")
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "read")
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "update")
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "delete")

	var (
		ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Application{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			"/application/",
			types.Application{Labels: map[string]string{"foo": "bar", "bar": "42"}},
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
			payload = &types.Application{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("PUT /application/%d", ID),
			&types.Application{Labels: map[string]string{"foo": "baz", "baz": "123"}},
			payload,
		)
		req.NotZero(payload.ID)
		//req.Nil(payload.UpdatedAt, "updatedAt must not change after changing labels")

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
			set = types.ApplicationSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t, "/application/", &set, url.Values{"labels": []string{"baz=123"}})
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
