package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

// We're using pkg/automation service layer for low-level tasks
func (h helper) svcMakeAutomationScript(namespace *types.Namespace, name string) *automation.Script {
	script := &automation.Script{
		NamespaceID: namespace.ID,
		Name:        name,
	}

	h.a.NoError(service.DefaultInternalAutomationManager.CreateScript(
		context.Background(),
		script,
	))

	return script
}

func TestAutomationScriptRead(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "read")

	ns := h.repoMakeNamespace("some-namespace")
	script := h.svcMakeAutomationScript(ns, "myLittleAutomationScript"+rs())

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/automation/script/%d", ns.ID, script.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, script.Name)).
		Assert(jsonpath.Equal(`$.response.scriptID`, fmt.Sprintf("%d", script.ID))).
		End()
}

func TestAutomationScriptList(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")

	ns := h.repoMakeNamespace("some-namespace")

	h.svcMakeAutomationScript(ns, "myLittleAutomationScript"+rs())
	h.svcMakeAutomationScript(ns, "myLittleAutomationScript"+rs())

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/automation/script/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAutomationScriptCreateForbidden(t *testing.T) {
	h := newHelper(t)

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/automation/script/", ns.ID)).
		FormData("name", "myLittleAutomationScript"+rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoCreatePermissions")).
		End()
}

func TestAutomationScriptCreate(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "automation-script.create")

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/automation/script/", ns.ID)).
		FormData("name", "myLittleAutomationScript"+rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAutomationScriptUpdateForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")

	ns := h.repoMakeNamespace("some-namespace")
	script := h.svcMakeAutomationScript(ns, "myLittleAutomationScript"+rs())

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/automation/script/%d", ns.ID, script.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoUpdatePermissions")).
		End()
}

func TestAutomationScriptUpdate(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "update")

	ns := h.repoMakeNamespace("some-namespace")
	script := h.svcMakeAutomationScript(ns, "myLittleAutomationScript"+rs())

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/automation/script/%d", ns.ID, script.ID)).
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

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")

	ns := h.repoMakeNamespace("some-namespace")
	script := h.svcMakeAutomationScript(ns, "myLittleAutomationScript"+rs())

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/automation/script/%d", ns.ID, script.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoDeletePermissions")).
		End()
}

func TestAutomationScriptDelete(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.AutomationScriptPermissionResource.AppendWildcard(), "delete")

	ns := h.repoMakeNamespace("some-namespace")
	script := h.svcMakeAutomationScript(ns, "myLittleAutomationScript"+rs())

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/automation/script/%d", ns.ID, script.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	script, err := service.DefaultInternalAutomationManager.FindScriptByID(context.Background(), script.ID)
	h.a.NoError(err)
	h.a.NotNil(script)
	h.a.NotNil(script.DeletedAt)
}
