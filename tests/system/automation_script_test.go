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
func (h helper) svcMakeAutomationScript(name string) *automation.Script {
	script := &automation.Script{
		Name: name,
	}

	h.a.NoError(service.DefaultInternalAutomationManager.CreateScript(
		context.Background(),
		script,
	))

	return script
}

func TestAutomationScriptRead(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")

	script := h.svcMakeAutomationScript("my little automation script")

	h.apiInit().
		Get(fmt.Sprintf("/automation/script/%d", script.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, script.Name)).
		Assert(jsonpath.Equal(`$.response.scriptID`, fmt.Sprintf("%d", script.ID))).
		End()
}

func TestAutomationScriptList(t *testing.T) {
	h := newHelper(t)

	h.svcMakeAutomationScript("app")
	h.svcMakeAutomationScript("app")

	h.apiInit().
		Get(fmt.Sprintf("/automation/script/")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAutomationScriptCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.deny(types.SystemPermissionResource, "automation-script.create")

	h.apiInit().
		Post(fmt.Sprintf("/automation/script/")).
		FormData("name", "my little automation script").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoCreatePermissions")).
		End()
}

func TestAutomationScriptCreate(t *testing.T) {
	h := newHelper(t)

	h.allow(types.SystemPermissionResource, "automation-script.create")

	h.apiInit().
		Post(fmt.Sprintf("/automation/script/")).
		FormData("name", "my little automation script").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAutomationScriptUpdateForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")

	script := h.svcMakeAutomationScript("my little automation script")

	h.apiInit().
		Post(fmt.Sprintf("/automation/script/%d", script.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoUpdatePermissions")).
		End()
}

func TestAutomationScriptUpdate(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")
	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "update")

	script := h.svcMakeAutomationScript("my little automation script")

	h.apiInit().
		Post(fmt.Sprintf("/automation/script/%d", script.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	script, err := service.DefaultInternalAutomationManager.FindScriptByID(context.Background(), script.ID)
	h.a.NoError(err)
	h.a.NotNil(script)
	h.a.Equal("changed-name", script.Name)
}

func TestAutomationScriptDeleteForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")

	script := h.svcMakeAutomationScript("my little automation script")

	h.apiInit().
		Delete(fmt.Sprintf("/automation/script/%d", script.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoDeletePermissions")).
		End()
}

func TestAutomationScriptDelete(t *testing.T) {
	h := newHelper(t)

	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")
	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "delete")

	script := h.svcMakeAutomationScript("my little automation script")

	h.apiInit().
		Delete(fmt.Sprintf("/automation/script/%d", script.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	script, err := service.DefaultInternalAutomationManager.FindScriptByID(context.Background(), script.ID)
	h.a.NoError(err)
	h.a.NotNil(script)
	h.a.NotNil(script.DeletedAt)
}
