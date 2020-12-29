package system

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func (h helper) clearTriggers() {
	h.noError(store.TruncateTriggers(context.Background(), service.DefaultStore))
}

func (h helper) createTrigger(res *types.Trigger) *types.Trigger {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateTrigger(context.Background(), res))
	return res
}

func (h helper) repoMakeTrigger(ss ...string) *types.Trigger {
	var r = &types.Trigger{
		ID:        id.Next(),
		CreatedAt: time.Now(),
	}

	if len(ss) > 1 {
		r.ResourceType = ss[1]
	} else {
		r.ResourceType = "h_" + rs()

	}

	h.a.NoError(store.CreateTrigger(context.Background(), service.DefaultStore, r))

	return r
}

func (h helper) lookupTriggerByID(ID uint64) *types.Trigger {
	res, err := store.LookupTriggerByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestTriggerRead(t *testing.T) {
	h := newHelper(t)
	h.clearTriggers()

	wf := h.repoMakeTrigger()
	h.allow(types.TriggerRBACResource.AppendID(wf.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/automation/triggers/%d", wf.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.resourceType`, wf.ResourceType)).
		Assert(jsonpath.Equal(`$.response.triggerID`, fmt.Sprintf("%d", wf.ID))).
		End()
}

func TestTriggerList(t *testing.T) {
	h := newHelper(t)
	h.clearTriggers()

	h.allow(types.TriggerRBACResource.AppendWildcard(), "read")

	h.repoMakeTrigger()
	h.repoMakeTrigger()

	h.apiInit().
		Get("/automation/triggers/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestTriggerList_filterForbidden(t *testing.T) {
	h := newHelper(t)

	// @todo this can be a problematic test because it leaves
	//       behind triggers that are not denied this context
	//       db purge might be needed

	h.repoMakeTrigger("trigger")
	f := h.repoMakeTrigger()

	h.deny(types.TriggerRBACResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/automation/triggers/").
		Query("resourceType", f.ResourceType).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(fmt.Sprintf(`$.response.set[? @.resourceType=="%s"]`, f.ResourceType))).
		End()
}

func TestTriggerCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/automation/triggers/").
		Header("Accept", "application/json").
		FormData("name", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create triggers")).
		End()
}

func TestTriggerCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemRBACResource, "trigger.create")

	h.apiInit().
		Post("/automation/triggers/").
		FormData("name", rs()).
		FormData("resourceType", "resourceType_"+rs()).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestTriggerCreateFull(t *testing.T) {
	h := newHelper(t)

	h.allow(types.SystemRBACResource, "trigger.create")

	h.clearTriggers()
	var (
		output = &types.Trigger{}
		stored = &types.Trigger{}
		input  = &types.Trigger{
			ResourceType: "wf-full-test",
			Enabled:      true,
			OwnedBy:      42,
		}
	)

	h.apiInit().
		Post("/automation/triggers/").
		Header("Accept", "application/json").
		JSON(helpers.JSON(input)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&struct{ Response *types.Trigger }{output})

	h.a.NotZero(output.ID)
	h.a.NotZero(output.OwnedBy)
	h.a.NotNil(output.CreatedAt)
	h.a.NotZero(output.CreatedBy)

	input.ID = output.ID
	input.OwnedBy = output.OwnedBy
	input.CreatedBy = output.CreatedBy
	input.CreatedAt = output.CreatedAt

	h.a.Equal(input, output)

	h.allow(types.TriggerRBACResource.AppendID(output.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/automation/triggers/%d", output.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&struct{ Response *types.Trigger }{stored})

	h.a.Equal(input, stored)

}

func TestTriggerUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeTrigger()

	h.apiInit().
		Put(fmt.Sprintf("/automation/triggers/%d", u.ID)).
		Header("Accept", "application/json").
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this trigger")).
		End()
}

func TestTriggerUpdate(t *testing.T) {
	h := newHelper(t)
	res := h.repoMakeTrigger()
	h.allow(types.TriggerRBACResource.AppendWildcard(), "update")

	newName := "updated-" + rs()
	newResourceType := "updated-" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/automation/triggers/%d", res.ID)).
		FormData("name", newName).
		FormData("resourceType", newResourceType).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupTriggerByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(newResourceType, res.ResourceType)
}

func TestTriggerDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeTrigger()

	h.apiInit().
		Delete(fmt.Sprintf("/automation/triggers/%d", u.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this trigger")).
		End()
}

func TestTriggerDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.TriggerRBACResource.AppendWildcard(), "delete")

	res := h.repoMakeTrigger()

	h.apiInit().
		Delete(fmt.Sprintf("/automation/triggers/%d", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupTriggerByID(res.ID)
	h.a.NotNil(res)
	h.a.NotNil(res.DeletedAt)
}

func TestTriggerLabels(t *testing.T) {
	h := newHelper(t)
	h.clearTriggers()

	h.allow(types.SystemRBACResource, "trigger.create")
	h.allow(types.TriggerRBACResource.AppendWildcard(), "read")
	h.allow(types.TriggerRBACResource.AppendWildcard(), "update")
	h.allow(types.TriggerRBACResource.AppendWildcard(), "delete")

	var (
		ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Trigger{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			"/automation/triggers/",
			types.Trigger{Labels: map[string]string{"foo": "bar", "bar": "42"}},
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
			payload = &types.Trigger{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("PUT /automation/triggers/%d", ID),
			&types.Trigger{Labels: map[string]string{"foo": "baz", "baz": "123"}},
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
			set = types.TriggerSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t, "/automation/triggers/", &set, url.Values{"labels": []string{"baz=123"}})
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
