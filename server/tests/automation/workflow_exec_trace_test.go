package automation

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// an empty workflow can fail in other ways, it just must get past the access control
func assertNoAccessControlError(rsp *http.Response, _ *http.Request) error {
	tmp := helpers.StdErrorResponse{}
	if err := helpers.DecodeBody(rsp, &tmp); err != nil {
		return err
	}

	if strings.Contains(tmp.Error.Message, "notAllowedTo") {
		return fmt.Errorf("unexpected access control error: %s", tmp.Error.Message)
	}

	return nil
}

// creates workflow over the API so that it gets registered with the workflow service
func (h helper) apiMakeWorkflow() *types.Workflow {
	handle := "handle_" + rs()
	helpers.AllowMe(h, types.ComponentRbacResource(), "workflow.create")

	h.apiInit().
		Post("/workflows/").
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", handle).
		FormData("enabled", "true").
		Header("Accept", "application/json").
		Expect(h.t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res, err := store.LookupAutomationWorkflowByHandle(context.Background(), service.DefaultStore, handle)
	h.noError(err)
	return res
}

// Tracing is for workflow designers: it runs disabled workflows, skips trigger
// checks and returns the whole scope; executing the workflow is not enough for that
func TestWorkflowExecTraceForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	wf := h.apiMakeWorkflow()

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "execute")
	helpers.DenyMe(h, types.WorkflowRbacResource(0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/workflows/%d/exec", wf.ID)).
		Header("Accept", "application/json").
		JSON(`{"trace": true, "wait": true}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("workflow.errors.notAllowedToUpdate")).
		End()
}

// Executing without the trace needs nothing but the execute permission
func TestWorkflowExecWithoutTrace(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	wf := h.apiMakeWorkflow()

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "execute")
	helpers.DenyMe(h, types.WorkflowRbacResource(0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/workflows/%d/exec", wf.ID)).
		Header("Accept", "application/json").
		JSON(`{"wait": true}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(assertNoAccessControlError).
		End()
}

// Designers can still trace
func TestWorkflowExecTraceAsDesigner(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	wf := h.apiMakeWorkflow()

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "execute")
	helpers.AllowMe(h, types.WorkflowRbacResource(0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/workflows/%d/exec", wf.ID)).
		Header("Accept", "application/json").
		JSON(`{"trace": true, "wait": true}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(assertNoAccessControlError).
		End()
}
