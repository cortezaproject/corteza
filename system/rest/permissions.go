package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

type (
	Permissions struct {
		ac permissionsAccessController
	}

	permissionsAccessController interface {
		Effective(context.Context) permissions.EffectiveSet
		Whitelist() permissions.Whitelist
		FindRulesByRoleID(context.Context, uint64) (permissions.RuleSet, error)
		Grant(ctx context.Context, rr ...*permissions.Rule) error
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

	_ = rr.Walk(func(rule *permissions.Rule) error {
		// Setting access to "inherit" will make Grant remove the rule
		rule.Access = permissions.Inherit
		return nil
	})

	return resputil.OK(), ctrl.ac.Grant(ctx, rr...)
}

func (ctrl Permissions) Update(ctx context.Context, r *request.PermissionsUpdate) (interface{}, error) {
	rr := r.Rules
	_ = rr.Walk(func(rule *permissions.Rule) error {
		// Make sure everything is properly set
		rule.RoleID = r.RoleID
		return nil
	})

	return resputil.OK(), ctrl.ac.Grant(ctx, rr...)
}
