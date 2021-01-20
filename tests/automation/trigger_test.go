package automation

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func (h helper) clearTriggers() {
	h.noError(store.TruncateAutomationTriggers(context.Background(), service.DefaultStore))
}

func (h helper) createTrigger(res *types.Trigger) *types.Trigger {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateAutomationTrigger(context.Background(), res))
	return res
}

func (h helper) repoMakeTrigger(wf *types.Workflow, ss ...string) *types.Trigger {
	var r = &types.Trigger{
		ID:         id.Next(),
		CreatedAt:  time.Now(),
		WorkflowID: wf.ID,
	}

	if len(ss) > 1 {
		r.ResourceType = ss[1]
	} else {
		r.ResourceType = "h_" + rs()

	}

	r.Enabled = true

	h.a.NoError(store.CreateAutomationTrigger(context.Background(), service.DefaultStore, r))

	return r
}

func (h helper) lookupTriggerByID(ID uint64) *types.Trigger {
	res, err := store.LookupAutomationTriggerByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestTriggerRead(t *testing.T) {
	h := newHelper(t)
	h.clearTriggers()

	wf := h.repoMakeWorkflow()
	tg := h.repoMakeTrigger(wf)

	h.allow(types.AutomationRBACResource, "triggers.search")

	h.apiInit().
		Get(fmt.Sprintf("/triggers/%d", tg.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.resourceType`, tg.ResourceType)).
		Assert(jsonpath.Equal(`$.response.triggerID`, fmt.Sprintf("%d", tg.ID))).
		End()
}

func TestTriggerList(t *testing.T) {
	h := newHelper(t)
	h.clearTriggers()

	h.allow(types.AutomationRBACResource, "triggers.search")

	wf := h.repoMakeWorkflow()
	h.repoMakeTrigger(wf)
	h.repoMakeTrigger(wf)

	h.apiInit().
		Get("/triggers/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestTriggerCreate(t *testing.T) {
	h := newHelper(t)
	h.clearTriggers()
	h.clearWorkflows()

	var (
		wf  = h.repoMakeWorkflow()
		req = func() *apitest.Request {
			return h.apiInit().
				Post("/triggers/").
				Header("Accept", "application/json").
				FormData("name", rs()).
				FormData("workflowID", fmt.Sprintf("%d", wf.ID))

		}
	)

	t.Run("allowed", func(t *testing.T) {
		h.allow(types.WorkflowRBACResource.AppendID(wf.ID), "triggers.manage")
		req().Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	})

	t.Run("denied", func(t *testing.T) {
		h.deny(types.WorkflowRBACResource.AppendID(wf.ID), "triggers.manage")
		req().Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertError("not allowed to create triggers")).
			End()
	})
}

func TestTriggerCreateFull(t *testing.T) {
	h := newHelper(t)

	h.allow(types.WorkflowRBACResource.AppendWildcard(), "triggers.manage")

	h.clearTriggers()
	var (
		wf     = h.repoMakeWorkflow()
		output = &types.Trigger{}
		stored = &types.Trigger{}
		input  = &types.Trigger{
			WorkflowID:   wf.ID,
			ResourceType: "wf-full-test",
			Enabled:      true,
			Input:        expr.RVars{}.Vars(),
			OwnedBy:      42,
		}
	)

	h.apiInit().
		Post("/triggers/").
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
	h.a.NoError(output.Input.ResolveTypes(service.Registry().Type))

	input.Meta = output.Meta

	h.a.Equal(input, output)

	h.allow(types.AutomationRBACResource, "triggers.search")

	h.apiInit().
		Debug().
		Get(fmt.Sprintf("/triggers/%d", output.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&struct{ Response *types.Trigger }{stored})

	h.a.NoError(stored.Input.ResolveTypes(service.Registry().Type))
	input.Meta = stored.Meta

	h.a.Equal(input, stored)
}

func TestTriggerUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearTriggers()
	h.clearWorkflows()

	var (
		wf  = h.repoMakeWorkflow()
		tg1 = h.createTrigger(&types.Trigger{WorkflowID: wf.ID, ResourceType: "as-created"})
		tg2 = h.createTrigger(&types.Trigger{WorkflowID: wf.ID, ResourceType: "as-created"})

		req = func(id uint64, resType string) *apitest.Request {
			return h.apiInit().
				Put(fmt.Sprintf("/triggers/%d", id)).
				Header("Accept", "application/json").
				FormData("resourceType", resType)

		}
	)

	t.Run("allowed", func(t *testing.T) {
		h.allow(types.WorkflowRBACResource.AppendID(wf.ID), "triggers.manage")
		req(tg1.ID, "foo").Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		res := h.lookupTriggerByID(tg1.ID)
		h.a.NotNil(res)
		h.a.Equal("foo", res.ResourceType)
	})

	t.Run("denied", func(t *testing.T) {
		h.deny(types.WorkflowRBACResource.AppendID(wf.ID), "triggers.manage")
		req(tg2.ID, "bar").Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertError("not allowed to update this trigger")).
			End()

		res := h.lookupTriggerByID(tg2.ID)
		h.a.NotNil(res)
		h.a.NotEqual("bar", res.ResourceType)
	})
}

func TestTriggerDeleteForbidden(t *testing.T) {
	h := newHelper(t)

	wf := h.repoMakeWorkflow()
	res := h.repoMakeTrigger(wf)

	h.apiInit().
		Delete(fmt.Sprintf("/triggers/%d", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this trigger")).
		End()

	res = h.lookupTriggerByID(res.ID)
	h.a.NotNil(res)
	h.a.Nil(res.DeletedAt)
}

func TestTriggerDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.WorkflowRBACResource.AppendWildcard(), "triggers.manage")

	wf := h.repoMakeWorkflow()
	res := h.repoMakeTrigger(wf)

	h.apiInit().
		Delete(fmt.Sprintf("/triggers/%d", res.ID)).
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

	h.allow(types.WorkflowRBACResource.AppendWildcard(), "triggers.manage")
	h.allow(types.AutomationRBACResource, "triggers.search")

	var (
		ID uint64
		wf = h.repoMakeWorkflow()
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Trigger{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			"/triggers/",
			types.Trigger{Labels: map[string]string{"foo": "bar", "bar": "42"}, WorkflowID: wf.ID},
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
			fmt.Sprintf("PUT /triggers/%d", ID),
			&types.Trigger{Labels: map[string]string{"foo": "baz", "baz": "123"}, WorkflowID: wf.ID},
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
			req    = require.New(t)
			set    = types.TriggerSet{}
			params = url.Values{}
		)
		params.Add("labels", "baz=123")
		params.Add("disabled", filter.StateInclusive.String())

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t, "/triggers/", &set, params)
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
