package helpers

import (
	"context"

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
