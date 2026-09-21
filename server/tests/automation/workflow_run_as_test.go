package automation

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/crusttech/human/server/automation/service"
	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
)

func (h helper) repoMakeRunAsUser() *sysTypes.User {
	u := &sysTypes.User{
		ID:        id.Next(),
		Email:     rs() + "@test.tld",
		Handle:    "admin_" + rs(),
		CreatedAt: time.Now(),
	}

	h.noError(store.CreateUser(context.Background(), service.DefaultStore, u))
	return u
}

func (h helper) repoMakeWorkflowRunAs(runAs uint64) *types.Workflow {
	res := h.repoMakeWorkflow()
	res.RunAs = runAs
	h.noError(store.UpdateAutomationWorkflow(context.Background(), service.DefaultStore, res))
	return res
}

func TestWorkflowCreateRunAsForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	admin := h.repoMakeRunAsUser()

	helpers.AllowMe(h, types.ComponentRbacResource(), "workflow.create")

	h.apiInit().
		Post("/workflows/").
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", "handle_"+rs()).
		FormData("runAs", fmt.Sprintf("%d", admin.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("workflow.errors.notAllowedToSetRunAs")).
		End()
}

func TestWorkflowCreateRunAs(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	admin := h.repoMakeRunAsUser()

	helpers.AllowMe(h, types.ComponentRbacResource(), "workflow.create")
	helpers.AllowMe(h, admin.RbacResource(), "impersonate")

	handle := "handle_" + rs()

	h.apiInit().
		Post("/workflows/").
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", handle).
		FormData("runAs", fmt.Sprintf("%d", admin.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res, err := store.LookupAutomationWorkflowByHandle(context.Background(), service.DefaultStore, handle)
	h.noError(err)
	h.a.Equal(admin.ID, res.RunAs)
}

func TestWorkflowCreateRunAsMyself(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()

	// current user needs to exist to be used as run-as user
	h.cUser.Email = rs() + "@test.tld"
	h.cUser.Handle = "myself_" + rs()
	h.cUser.CreatedAt = time.Now()
	h.noError(store.CreateUser(context.Background(), service.DefaultStore, h.cUser))

	helpers.AllowMe(h, types.ComponentRbacResource(), "workflow.create")

	handle := "handle_" + rs()

	h.apiInit().
		Post("/workflows/").
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", handle).
		FormData("runAs", fmt.Sprintf("%d", h.cUser.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res, err := store.LookupAutomationWorkflowByHandle(context.Background(), service.DefaultStore, handle)
	h.noError(err)
	h.a.Equal(h.cUser.ID, res.RunAs)
}

func TestWorkflowUpdateRunAsForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	admin := h.repoMakeRunAsUser()
	res := h.repoMakeWorkflow()

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "update")

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", res.ID)).
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", res.Handle).
		FormData("runAs", fmt.Sprintf("%d", admin.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("workflow.errors.notAllowedToSetRunAs")).
		End()

	h.a.Zero(h.lookupWorkflowByID(res.ID).RunAs)
}

func TestWorkflowUpdateRunAs(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	admin := h.repoMakeRunAsUser()
	res := h.repoMakeWorkflow()

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "update")
	helpers.AllowMe(h, admin.RbacResource(), "impersonate")

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", res.ID)).
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", res.Handle).
		FormData("runAs", fmt.Sprintf("%d", admin.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.Equal(admin.ID, h.lookupWorkflowByID(res.ID).RunAs)
}

// Workflows that already run as someone else can still be updated
func TestWorkflowUpdateKeepRunAs(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	admin := h.repoMakeRunAsUser()
	res := h.repoMakeWorkflowRunAs(admin.ID)

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "update")

	newHandle := "updated_" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", res.ID)).
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", newHandle).
		FormData("runAs", fmt.Sprintf("%d", admin.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupWorkflowByID(res.ID)
	h.a.Equal(newHandle, res.Handle)
	h.a.Equal(admin.ID, res.RunAs)
}

func TestWorkflowUpdateClearRunAs(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	admin := h.repoMakeRunAsUser()
	res := h.repoMakeWorkflowRunAs(admin.ID)

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "update")

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", res.ID)).
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", res.Handle).
		FormData("runAs", "0").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.Zero(h.lookupWorkflowByID(res.ID).RunAs)
}

func TestWorkflowUpdateOwnerForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	res := h.repoMakeWorkflow()
	owner := res.OwnedBy

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "update")
	helpers.DenyMe(h, types.ComponentRbacResource(), "grant")

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", res.ID)).
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", res.Handle).
		FormData("ownedBy", fmt.Sprintf("%d", h.cUser.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("workflow.errors.notAllowedToChangeOwner")).
		End()

	h.a.Equal(owner, h.lookupWorkflowByID(res.ID).OwnedBy)
}

func TestWorkflowUpdateOwner(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	res := h.repoMakeWorkflow()

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "update")
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", res.ID)).
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", res.Handle).
		FormData("ownedBy", fmt.Sprintf("%d", h.cUser.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.Equal(h.cUser.ID, h.lookupWorkflowByID(res.ID).OwnedBy)
}

// Updates that do not send the owner must not reset or change it
func TestWorkflowUpdateKeepOwner(t *testing.T) {
	h := newHelper(t)
	h.clearWorkflows()
	res := h.repoMakeWorkflow()
	res.OwnedBy = id.Next()
	h.noError(store.UpdateAutomationWorkflow(context.Background(), service.DefaultStore, res))

	helpers.AllowMe(h, types.WorkflowRbacResource(0), "update")
	helpers.DenyMe(h, types.ComponentRbacResource(), "grant")

	h.apiInit().
		Put(fmt.Sprintf("/workflows/%d", res.ID)).
		FormData("meta", fmt.Sprintf(`{"name": "%s"}`, rs())).
		FormData("handle", res.Handle).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.Equal(res.OwnedBy, h.lookupWorkflowByID(res.ID).OwnedBy)
}
