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
	"github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func (h helper) clearWorkflows() {
	h.noError(store.TruncateAutomationWorkflows(context.Background(), service.DefaultStore))
}

func (h helper) createWorkflow(res *types.Workflow) *types.Workflow {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateAutomationWorkflow(context.Background(), res))
	return res
}

func (h helper) repoMakeWorkflow(ss ...string) *types.Workflow {
	var r = &types.Workflow{
		ID:        id.Next(),
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	if len(ss) > 1 {
		r.Handle = ss[1]
	} else {
		r.Handle = "h_" + rs()

	}

	h.a.NoError(store.CreateAutomationWorkflow(context.Background(), service.DefaultStore, r))

	return r
}

func (h helper) lookupWorkflowByID(ID uint64) *types.Workflow {
	res, err := store.LookupAutomationWorkflowByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestWorkflowRead(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()

	wf := h.repoMakeWorkflow()
	h.allow(types.WorkflowRBACResource.AppendID(wf.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/workflows/%d", wf.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.handle`, wf.Handle)).
		Assert(jsonpath.Equal(`$.response.workflowID`, fmt.Sprintf("%d", wf.ID))).
		End()
}

func TestWorkflowList(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()

	h.allow(types.WorkflowRBACResource.AppendWildcard(), "read")

	h.repoMakeWorkflow()
	h.repoMakeWorkflow()

	h.apiInit().
		Get("/workflows/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestWorkflowList_filterForbidden(t *testing.T) {
	h := newHelper(t)

	// @todo this can be a problematic test because it leaves
	//       behind workflows that are not denied this context
	//       db purge might be needed

	h.repoMakeWorkflow("workflow")
	f := h.repoMakeWorkflow()

	h.deny(types.WorkflowRBACResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/workflows/").
		Query("handle", f.Handle).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(fmt.Sprintf(`$.response.set[? @.handle=="%s"]`, f.Handle))).
		End()
}

func TestWorkflowCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/workflows/").
		Header("Accept", "application/json").
		FormData("name", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create workflows")).
		End()
}

func TestWorkflowCreateNotUnique(t *testing.T) {
	h := newHelper(t)
	h.allow(types.AutomationRBACResource, "workflow.create")

	workflow := h.repoMakeWorkflow()
	h.apiInit().
		Post("/workflows/").
		Header("Accept", "application/json").
		FormData("name", rs()).
		FormData("handle", workflow.Handle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("workflow handle not unique")).
		End()
}

func TestWorkflowCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.AutomationRBACResource, "workflow.create")

	h.apiInit().
		Post("/workflows/").
		FormData("name", rs()).
		FormData("handle", "handle_"+rs()).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestWorkflowCreateFull(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationRBACResource, "workflow.create")

	h.clearWorkflows()
	var (
		output = &types.Workflow{}
		stored = &types.Workflow{}
		input  = &types.Workflow{
			Handle: "wf-full-test",
			Meta: &types.WorkflowMeta{
				Name:        "name",
				Description: "desc",
				Visual:      map[string]interface{}{"foo": "bar"},
			},
			Enabled:      true,
			Trace:        true,
			KeepSessions: 10000,
			Scope:        expr.RVars{"foo": expr.Must(expr.NewString("bar"))}.Vars(),
			Steps: types.WorkflowStepSet{
				{ID: 1, Kind: types.WorkflowStepKindVisual, Meta: types.WorkflowStepMeta{Visual: map[string]interface{}{"foo": "bar"}}},
				{ID: 2, Kind: types.WorkflowStepKindVisual},
			},
			Paths: types.WorkflowPathSet{
				{ParentID: 1, ChildID: 2, Meta: types.WorkflowPathMeta{Visual: map[string]interface{}{"foo": "bar"}}},
			},
			RunAs:   42,
			OwnedBy: 42,
		}
	)

	h.apiInit().
		Post("/workflows/").
		Header("Accept", "application/json").
		JSON(helpers.JSON(input)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&struct{ Response *types.Workflow }{output})

	h.a.NotZero(output.ID)
	h.a.NotZero(output.OwnedBy)
	h.a.NotNil(output.CreatedAt)
	h.a.NotZero(output.CreatedBy)

	// Make comparing input&output a bit easier
	input.ID = output.ID
	input.OwnedBy = output.OwnedBy
	input.CreatedBy = output.CreatedBy
	input.CreatedAt = output.CreatedAt
	h.a.NoError(output.Scope.ResolveTypes(service.Registry().Type))

	h.a.Equal(input, output)

	h.allow(types.WorkflowRBACResource.AppendID(output.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/workflows/%d", output.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&struct{ Response *types.Workflow }{stored})

	h.a.NoError(stored.Scope.ResolveTypes(service.Registry().Type))

	h.a.Equal(input, stored)
}

func TestWorkflowUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeWorkflow()

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", u.ID)).
		Header("Accept", "application/json").
		FormData("email", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this workflow")).
		End()
}

func TestWorkflowUpdate(t *testing.T) {
	h := newHelper(t)
	res := h.repoMakeWorkflow()
	h.allow(types.WorkflowRBACResource.AppendWildcard(), "update")

	newName := "updated-" + rs()
	newHandle := "updated-" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", res.ID)).
		FormData("name", newName).
		FormData("handle", newHandle).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupWorkflowByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(newHandle, res.Handle)
}

func TestWorkflowDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeWorkflow()

	h.apiInit().
		Delete(fmt.Sprintf("/workflows/%d", u.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this workflow")).
		End()
}

func TestWorkflowDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.WorkflowRBACResource.AppendWildcard(), "delete")

	res := h.repoMakeWorkflow()

	h.apiInit().
		Delete(fmt.Sprintf("/workflows/%d", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupWorkflowByID(res.ID)
	h.a.NotNil(res)
	h.a.NotNil(res.DeletedAt)
}

func TestWorkflowLabels(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()

	h.allow(types.AutomationRBACResource, "workflow.create")
	h.allow(types.WorkflowRBACResource.AppendWildcard(), "read")
	h.allow(types.WorkflowRBACResource.AppendWildcard(), "update")
	h.allow(types.WorkflowRBACResource.AppendWildcard(), "delete")

	var (
		ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Workflow{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			"/workflows/",
			types.Workflow{Labels: map[string]string{"foo": "bar", "bar": "42"}},
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
			payload = &types.Workflow{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("PUT /workflows/%d", ID),
			&types.Workflow{Labels: map[string]string{"foo": "baz", "baz": "123"}},
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
			set    = types.WorkflowSet{}
			params = url.Values{}
		)

		params.Add("labels", "baz=123")
		params.Add("disabled", filter.StateInclusive.String())

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t, "/workflows/", &set, params)
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}

func TestWorkflowStepsPayload(t *testing.T) {
	wf := &types.Workflow{}
	h := newHelper(t)
	h.allow(types.AutomationRBACResource, "workflow.create")

	h.apiInit().
		Post("/workflows/").
		JSON(`{
	"scope": {
		"foo": { "@type": "String", "@value": "" },
		"bar": { "@type": "String", "@value": "" },
		"baz": { "@type": "String", "@value": "" },
		"hgtg": { "@type": "Integer", "@value": 42 }
	},
	
	"meta": {
		"name": "name of the workdflow",
		"description": "lorem ipsum dolor sit...",
		"visual":  { "arbitrary key": "arbitrary value" }
	},
	
	"steps": [
		{ "stepID": "100",
		  "kind":   "expressions",
		  "arguments": [
			  { "name":  "foo", "expr": "\"foo\"" },
			  { "name":  "bar", "expr": "\"bar\"" }
		  ],
		  "meta": { "name": "step name text", "description": "step description text", "visual":  { "arbitrary key": "arbitrary value" }}
		},
		{ "stepID": "101",
		  "kind":   "expressions",
		  
		  "arguments": [
			  { "name":  "foo", "expr": "\"foo\"" },
			  { "name":  "bar", "expr": "\"bar\"" }
		  ] 
		},
		
		{ "stepID": "200", "kind": "gateway", "ref": "join" },
		{ "stepID": "201", "kind": "gateway", "ref": "fork" },
		{ "stepID": "202", "kind": "gateway", "ref": "excl" },
		{ "stepID": "203", "kind": "gateway", "ref": "incl" },
		
		{ "stepID": "300",
		  "kind":   "function",
		  "ref":    "serverSideFuncRef",
		  "meta": { "description": "ref needs to come from a predefined list of registered functions; arguments pick values from scope and need to be compatible with registered function's parameters; same goes for results" },
		  
		  "arguments": [
			  { "name":  "arg1", "expr": "foo" },
			  { "name":  "arg2", "expr": "hgtg * 2" }
		  ],
		  
		  "results": [
			  { "name":  "scopeVar1", "expr": "nameOfVarFromFunctionsResult1" },
			  { "name":  "scopeVar2", "expr": "nameOfVarFromFunctionsResult2" }
		  ] 
		},
		{ "stepID": "400",
		  "kind":   "prompt",
		  "ref":    "clientPromptRef",
		  "meta": { "description": "backed does not need to know about ref, it only cares that input data (when resuming session) is compatible with arguments (and tests there)" },
		  
		  "arguments": [
			  { "name":  "foo", "expr": "inputVar1",
				"tests": [
					{ "expr": "inputVar1 != \"\"", "error": "expecting inputVar1 to be non-empty string" }
				]
			  },
			  { "name":  "bar", "expr": "inputVar2" }
		  ]
		}
	],
	
	"paths": [
		{ "parentID":  "201", "childID":  "300", "expr": "hgtg == 42",
		  "meta": { "name": "...", "description": "only if equals 42" } },
		{ "parentID":  "201", "childID":  "400",
		  "meta": { "name": "...", "description": "else" } }
	]
}`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&struct{ Response *types.Workflow }{Response: wf})
}
