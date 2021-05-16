package system

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	ft "github.com/cortezaproject/corteza-server/pkg/flag/types"
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

func (h helper) repoFlagApplication(appID, owner uint64, flag string, active bool) *ft.Flag {
	res := &ft.Flag{
		Kind:       "system:application",
		ResourceID: appID,
		OwnedBy:    owner,
		Name:       flag,
		Active:     active,
	}

	h.a.NoError(store.CreateFlag(context.Background(), service.DefaultStore, res))
	return res
}

func (h helper) lookupApplicationByID(ID uint64) *types.Application {
	res, err := store.LookupApplicationByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func (h helper) searchApplicationFlags(ID, owner uint64, flag string) ft.FlagSet {
	res, _, err := store.SearchFlags(context.Background(), service.DefaultStore, ft.FlagFilter{
		Kind:       "system:application",
		ResourceID: []uint64{ID},
		OwnedBy:    []uint64{owner},
		Name:       []string{flag},
	})
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

	h.deny(f.RbacResource(), "read")

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
	h.allow(types.ComponentRbacResource(), "application.create")

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
	h.allow(types.ComponentRbacResource(), "application.create")
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
	h.allow(types.ApplicationRbacResource(0), "update")

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
	h.allow(types.ApplicationRbacResource(0), "update")

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
	h.allow(types.ApplicationRbacResource(0), "update")
	a := h.repoMakeApplication()
	b := h.repoMakeApplication()
	c := h.repoMakeApplication()
	h.deny(b.RbacResource(), "update")

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
	h.allow(types.ApplicationRbacResource(0), "update")
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
	h.allow(types.ApplicationRbacResource(0), "delete")

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

func TestApplicationUndelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ApplicationRBACResource.AppendWildcard(), "delete")

	res := h.repoMakeApplication()

	h.apiInit().
		Post(fmt.Sprintf("/application/%d/undelete", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupApplicationByID(res.ID)
	h.a.NotNil(res)
	h.a.Nil(res.DeletedAt)
}

func TestApplicationLabels(t *testing.T) {
	h := newHelper(t)
	h.clearApplications()

	h.allow(types.ComponentRbacResource(), "application.create")
	h.allow(types.ApplicationRbacResource(0), "read")
	h.allow(types.ApplicationRbacResource(0), "update")
	h.allow(types.ApplicationRbacResource(0), "delete")

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

func TestApplicationFlags(t *testing.T) {
	h := newHelper(t)
	h.clearApplications()

	h.allow(types.ComponentRbacResource(), "application.create")

	t.Run("create", func(t *testing.T) {
		h.allow(types.ComponentRbacResource(), "application.flag.global")
		res := h.repoMakeApplication()

		h.apiInit().
			Post(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, 0)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		ff := h.searchApplicationFlags(res.ID, 0, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)
	})

	t.Run("create; not allowed", func(t *testing.T) {
		h.deny(types.ComponentRbacResource(), "application.flag.global")
		res := h.repoMakeApplication()

		h.apiInit().
			Post(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, 0)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertError("not allowed to manage global flags for applications")).
			End()
	})

	t.Run("create own", func(t *testing.T) {
		h.allow(types.ComponentRbacResource(), "application.flag.self")
		res := h.repoMakeApplication()
		h.repoFlagApplication(res.ID, 0, "testFlag", true)

		h.apiInit().
			Post(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, 10)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		ff := h.searchApplicationFlags(res.ID, 0, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)

		ff = h.searchApplicationFlags(res.ID, 10, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)
	})

	t.Run("create own; not allowed", func(t *testing.T) {
		h.deny(types.ComponentRbacResource(), "application.flag.self")
		res := h.repoMakeApplication()
		h.repoFlagApplication(res.ID, 0, "testFlag", true)

		h.apiInit().
			Post(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, 10)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertError("not allowed to manage flags for applications")).
			End()
	})

	t.Run("read application", func(t *testing.T) {
		h.allow(types.ApplicationRbacResource(0), "read")
		res := h.repoMakeApplication()
		h.repoFlagApplication(res.ID, 0, "testFlag", true)

		h.apiInit().
			Get(fmt.Sprintf("/application/%d", res.ID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Present(`$.response.flags`)).
			Assert(jsonpath.Len(`$.response.flags`, 1)).
			Assert(jsonpath.Equal(`$.response.flags[0]`, "testFlag")).
			End()
	})

	t.Run("list applications", func(t *testing.T) {
		h.allow(types.ApplicationRbacResource(0), "read")
		h.clearApplications()
		res := h.repoMakeApplication()
		h.repoFlagApplication(res.ID, 0, "testFlag", true)

		h.apiInit().
			Get("/application/").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Len(`$.response.set`, 1)).
			Assert(jsonpath.Len(`$.response.set[0].flags`, 1)).
			Assert(jsonpath.Equal(`$.response.set[0].flags[0]`, "testFlag")).
			End()
	})

	t.Run("read application; with own flag", func(t *testing.T) {
		h.allow(types.ApplicationRbacResource(0), "read")
		res := h.repoMakeApplication()
		h.repoFlagApplication(res.ID, 0, "testFlag", true)
		h.repoFlagApplication(res.ID, h.cUser.ID, "testFlagOwn", true)

		h.apiInit().
			Get(fmt.Sprintf("/application/%d", res.ID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Present(`$.response.flags`)).
			Assert(jsonpath.Len(`$.response.flags`, 2)).
			End()
	})

	t.Run("read application; overwrite global", func(t *testing.T) {
		h.allow(types.ApplicationRbacResource(0), "read")
		res := h.repoMakeApplication()
		h.repoFlagApplication(res.ID, 0, "testFlag", true)
		h.repoFlagApplication(res.ID, h.cUser.ID, "testFlag", false)

		h.apiInit().
			Get(fmt.Sprintf("/application/%d", res.ID)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.NotPresent(`$.response.flags`)).
			End()
	})

	t.Run("filter by flags", func(t *testing.T) {
		flag := rs()
		h.allow(types.ApplicationRbacResource(0), "read")
		h.repoMakeApplication()
		h.repoMakeApplication()
		res := h.repoMakeApplication()
		h.repoFlagApplication(res.ID, 0, flag, true)

		h.apiInit().
			Get(fmt.Sprintf("/application/")).
			QueryCollection(url.Values{"flags": []string{flag}}).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Len("$.response.set", 1)).
			End()
	})

	t.Run("filter by flags; self inactive", func(t *testing.T) {
		flag := rs()
		h.allow(types.ApplicationRbacResource(0), "read")
		h.repoMakeApplication()
		h.repoMakeApplication()
		res := h.repoMakeApplication()
		h.repoFlagApplication(res.ID, 0, flag, true)
		h.repoFlagApplication(res.ID, h.cUser.ID, flag, false)

		h.apiInit().
			Get(fmt.Sprintf("/application/")).
			QueryCollection(url.Values{"flags": []string{flag}}).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Len("$.response.set", 0)).
			End()
	})
}

func TestApplicationFlags_Flow1(t *testing.T) {
	h := newHelper(t)
	h.clearApplications()

	h.allow(types.ComponentRbacResource(), "application.create")

	t.Run("create", func(t *testing.T) {
		h.allow(types.ComponentRbacResource(), "application.flag.global")
		h.allow(types.ComponentRbacResource(), "application.flag.self")
		res := h.repoMakeApplication()
		a := h.apiInit()

		a.Post(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, h.cUser.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		a.Post(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, 0)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		a.Delete(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, h.cUser.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		ff := h.searchApplicationFlags(res.ID, 0, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)
		h.a.True(ff[0].Active)

		ff = h.searchApplicationFlags(res.ID, h.cUser.ID, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)
		h.a.False(ff[0].Active)

		a.Delete(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, h.cUser.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		ff = h.searchApplicationFlags(res.ID, 0, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)
		h.a.True(ff[0].Active)

		ff = h.searchApplicationFlags(res.ID, h.cUser.ID, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)
		h.a.False(ff[0].Active)

		a.Post(fmt.Sprintf("/application/%d/flag/%d/testFlag", res.ID, h.cUser.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		ff = h.searchApplicationFlags(res.ID, 0, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)
		h.a.True(ff[0].Active)

		ff = h.searchApplicationFlags(res.ID, h.cUser.ID, "testFlag")
		h.a.NotNil(ff)
		h.a.Len(ff, 1)
		h.a.True(ff[0].Active)
	})
}
