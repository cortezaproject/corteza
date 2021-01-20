package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"

	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	Permissions struct {
		ac permissionsAccessController
	}

	permissionsAccessController interface {
		Effective(context.Context) rbac.EffectiveSet
		Whitelist() rbac.Whitelist
		FindRulesByRoleID(context.Context, uint64) (rbac.RuleSet, error)
		Grant(ctx context.Context, rr ...*rbac.Rule) error
	}
)

func (Permissions) New() *Permissions {
	return &Permissions{
		ac: service.DefaultAccessControl,
	}
}

func (ctrl Permissions) Effective(ctx context.Context, r *request.PermissionsEffective) (interface{}, error) {
	return ctrl.ac.Effective(ctx), nil
}

func (ctrl Permissions) List(ctx context.Context, r *request.PermissionsList) (interface{}, error) {
	return ctrl.ac.Whitelist().Flatten(), nil
}

func (ctrl Permissions) Read(ctx context.Context, r *request.PermissionsRead) (interface{}, error) {
	return ctrl.ac.FindRulesByRoleID(ctx, r.RoleID)
}

func (ctrl Permissions) Delete(ctx context.Context, r *request.PermissionsDelete) (interface{}, error) {
	rr, err := ctrl.ac.FindRulesByRoleID(ctx, r.RoleID)
	if err != nil {
		return nil, err
	}

	_ = rr.Walk(func(rule *rbac.Rule) error {
		// Setting access to "inherit" will make Grant remove the rule
		rule.Access = rbac.Inherit
		return nil
	})

	return api.OK(), ctrl.ac.Grant(ctx, rr...)
}

func (ctrl Permissions) Update(ctx context.Context, r *request.PermissionsUpdate) (interface{}, error) {
	rr := r.Rules
	_ = rr.Walk(func(rule *rbac.Rule) error {
		// Make sure everything is properly set
		rule.RoleID = r.RoleID
		return nil
	})

	return api.OK(), ctrl.ac.Grant(ctx, rr...)
}
