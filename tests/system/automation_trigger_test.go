package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

// We're using pkg/automation service layer for low-level tasks
func (h helper) svcMakeAutomationTrigger(script *automation.Script, event string) *automation.Trigger {
	trigger := &automation.Trigger{
		ScriptID: script.ID,
		Event:    event,
	}

	h.a.NoError(service.DefaultInternalAutomationManager.CreateTrigger(
		context.Background(),
		script,
		trigger,
	))

	return trigger
}

func TestAutomationTriggerRead(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")

	script := h.svcMakeAutomationScript("dummy")
	trigger := h.svcMakeAutomationTrigger(script, "manual")

	h.apiInit().
		Get(fmt.Sprintf("/automation/script/%d/trigger/%d", script.ID, trigger.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.event`, trigger.Event)).
		Assert(jsonpath.Equal(`$.response.triggerID`, fmt.Sprintf("%d", trigger.ID))).
		End()
}

func TestAutomationTriggerList(t *testing.T) {
	h := newHelper(t)

	script := h.svcMakeAutomationScript("dummy")

	h.svcMakeAutomationTrigger(script, "app")
	h.svcMakeAutomationTrigger(script, "app")

	h.apiInit().
		Get(fmt.Sprintf("/automation/script/%d/trigger/", script.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAutomationTriggerCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")

	script := h.svcMakeAutomationScript("dummy")

	h.apiInit().
		Post(fmt.Sprintf("/automation/script/%d/trigger/", script.ID)).
		FormData("event", "my-event").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoTriggerManagementPermissions")).
		End()
}

func TestAutomationTriggerCreate(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")
	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "update")

	script := h.svcMakeAutomationScript("dummy")

	h.apiInit().
		Post(fmt.Sprintf("/automation/script/%d/trigger/", script.ID)).
		FormData("event", "my-event").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAutomationTriggerUpdateForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")

	script := h.svcMakeAutomationScript("dummy")
	trigger := h.svcMakeAutomationTrigger(script, "my little automation trigger")

	h.apiInit().
		Post(fmt.Sprintf("/automation/script/%d/trigger/%d", script.ID, trigger.ID)).
		FormData("name", "manual").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoTriggerManagementPermissions")).
		End()
}

func TestAutomationTriggerUpdate(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")
	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "update")

	script := h.svcMakeAutomationScript("dummy")
	trigger := h.svcMakeAutomationTrigger(script, "my little automation trigger")

	h.apiInit().
		Post(fmt.Sprintf("/automation/script/%d/trigger/%d", script.ID, trigger.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	trigger, err := service.DefaultInternalAutomationManager.FindTriggerByID(context.Background(), trigger.ID)
	h.a.NoError(err)
	h.a.NotNil(trigger)
}

func TestAutomationTriggerDeleteForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")

	script := h.svcMakeAutomationScript("dummy")
	trigger := h.svcMakeAutomationTrigger(script, "my little automation trigger")

	h.apiInit().
		Delete(fmt.Sprintf("/automation/script/%d/trigger/%d", script.ID, trigger.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoTriggerManagementPermissions")).
		End()
}

func TestAutomationTriggerDelete(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")
	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "update")

	script := h.svcMakeAutomationScript("dummy")
	trigger := h.svcMakeAutomationTrigger(script, "my little automation trigger")

	h.apiInit().
		Delete(fmt.Sprintf("/automation/script/%d/trigger/%d", script.ID, trigger.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	trigger, err := service.DefaultInternalAutomationManager.FindTriggerByID(context.Background(), trigger.ID)
	h.a.NoError(err)
	h.a.NotNil(trigger)
	h.a.NotNil(trigger.DeletedAt)
}
