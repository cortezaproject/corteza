package helpers

import (
	"context"

	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	composeTypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	myRoleGetter interface {
		MyRole() uint64
	}
)

func UpdateRBAC(rr ...uint64) {
	// make sure event listener for role changes is removed
	eventbus.Service().UnregisterByResource("system:role")

	// convert given list of role IDs to RBAC roles
	// and update service
	ccr := []*rbac.Role{}

	// Make sure all bypass roles are set on RBAC service
	for _, r := range auth.ServiceUser().Roles() {
		ccr = append(ccr, rbac.BypassRole.Make(r, auth.BypassRoleHandle))
	}

	for _, r := range rr {
		ccr = append(ccr, rbac.CommonRole.Make(r, "integration-tests"))
	}
	rbac.Global().Clear()
	rbac.Global().UpdateRoles(ccr...)
}

func AllowMe(mrg myRoleGetter, r string, oo ...string) {
	for _, o := range oo {
		Grant(rbac.AllowRule(mrg.MyRole(), r, o))
	}
}

func Allow(role *types.Role, r string, oo ...string) {
	for _, o := range oo {
		Grant(rbac.AllowRule(role.ID, r, o))
	}
}

func DenyMe(mrg myRoleGetter, r string, oo ...string) {
	for _, o := range oo {
		Grant(rbac.DenyRule(mrg.MyRole(), r, o))
	}
}

func Deny(role *types.Role, r string, oo ...string) {
	for _, o := range oo {
		Grant(rbac.DenyRule(role.ID, r, o))
	}
}

func Grant(rr ...*rbac.Rule) {
	cli.HandleError(rbac.Global().Grant(context.Background(), rr...))
}

// Common RBAC presets

func AllowMeModuleCreate(mrg myRoleGetter) {
	AllowMe(mrg, composeTypes.NamespaceRbacResource(0), "module.create")
}

func AllowMeWorkflowSearch(mrg myRoleGetter) {
	AllowMe(mrg, automationTypes.WorkflowRbacResource(0), "workflows.search")
	AllowMe(mrg, automationTypes.WorkflowRbacResource(0), "triggers.search")
}

func AllowMeModuleSearch(mrg myRoleGetter) {
	AllowMe(mrg, composeTypes.NamespaceRbacResource(0), "modules.search")
	AllowMe(mrg, composeTypes.ModuleRbacResource(0, 0), "read")
}

func AllowMeModuleCRUD(mrg myRoleGetter) {
	AllowMe(mrg, composeTypes.NamespaceRbacResource(0), "module.create")
	AllowMe(mrg, composeTypes.NamespaceRbacResource(0), "modules.search")
	AllowMe(mrg, composeTypes.ModuleRbacResource(0, 0), "read")
	AllowMe(mrg, composeTypes.ModuleRbacResource(0, 0), "update")
	AllowMe(mrg, composeTypes.ModuleRbacResource(0, 0), "delete")
}

func AllowMeRecordCRUD(mrg myRoleGetter) {
	AllowMe(mrg, composeTypes.ModuleRbacResource(0, 0), "record.create")
	AllowMe(mrg, composeTypes.ModuleRbacResource(0, 0), "records.search")
	AllowMe(mrg, composeTypes.RecordRbacResource(0, 0, 0), "read")
	AllowMe(mrg, composeTypes.RecordRbacResource(0, 0, 0), "update")
	AllowMe(mrg, composeTypes.RecordRbacResource(0, 0, 0), "delete")
	AllowMe(mrg, composeTypes.ModuleFieldRbacResource(0, 0, 0), "record.value.update")
	AllowMe(mrg, composeTypes.ModuleFieldRbacResource(0, 0, 0), "record.value.read")
}

func AllowMeDalConnectionSearch(mrg myRoleGetter) {
	AllowMe(mrg, types.ComponentRbacResource(), "dal-connections.search")
	AllowMe(mrg, types.DalConnectionRbacResource(0), "read")
}

func AllowMeDalConnectionCRUD(mrg myRoleGetter) {
	AllowMe(mrg, types.ComponentRbacResource(),
		"dal-connection.create",
		"dal-connections.search",
	)

	AllowMe(mrg, types.DalConnectionRbacResource(0),
		"read",
		"update",
		"delete",
		"dal-config.manage",
	)
}
