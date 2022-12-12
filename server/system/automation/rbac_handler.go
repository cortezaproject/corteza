package automation

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	rbacService interface {
		Can(ses rbac.Session, op string, res rbac.Resource) bool
		Grant(ctx context.Context, rules ...*rbac.Rule) (err error)
	}

	rbacUserService interface {
		FindByAny(ctx context.Context, identifier interface{}) (*types.User, error)
	}

	rbacHandler struct {
		reg  rbacHandlerRegistry
		rbac rbacService
		user rbacUserService
		role roleService
	}

	auxResource struct{ res string }
)

var _ rbac.Resource = &auxResource{}

func (r auxResource) RbacResource() string {
	return r.res
}

func RbacHandler(reg rbacHandlerRegistry, rbac rbacService, user rbacUserService, role roleService) *rbacHandler {
	h := &rbacHandler{
		reg:  reg,
		rbac: rbac,
		user: user,
		role: role,
	}

	h.register()
	return h
}

func (h rbacHandler) allow(ctx context.Context, args *rbacAllowArgs) (err error) {
	r, err := lookupRole(ctx, h.role, rolesLookupArgs{
		hasLookup:    args.hasRole,
		Lookup:       args.Role,
		lookupID:     args.roleID,
		lookupHandle: args.roleHandle,
		lookupRes:    args.roleRes,
	})
	if err != nil {
		return
	}

	return h.grant(ctx, rbac.AllowRule(r.ID, args.Resource.RbacResource(), args.Operation))
}

func (h rbacHandler) deny(ctx context.Context, args *rbacDenyArgs) (err error) {
	r, err := lookupRole(ctx, h.role, rolesLookupArgs{
		hasLookup:    args.hasRole,
		Lookup:       args.Role,
		lookupID:     args.roleID,
		lookupHandle: args.roleHandle,
		lookupRes:    args.roleRes,
	})
	if err != nil {
		return
	}

	return h.grant(ctx, rbac.DenyRule(r.ID, args.Resource.RbacResource(), args.Operation))
}

func (h rbacHandler) inherit(ctx context.Context, args *rbacInheritArgs) (err error) {
	r, err := lookupRole(ctx, h.role, rolesLookupArgs{
		hasLookup:    args.hasRole,
		Lookup:       args.Role,
		lookupID:     args.roleID,
		lookupHandle: args.roleHandle,
		lookupRes:    args.roleRes,
	})
	if err != nil {
		return
	}

	return h.grant(ctx, rbac.InheritRule(r.ID, args.Resource.RbacResource(), args.Operation))
}

// verifies grant op (granter needs to be allowed to grant on the component

func (h rbacHandler) grant(ctx context.Context, r *rbac.Rule) (err error) {
	// we can safely create a fake component resource since there are
	// no extra metadata hidden under it,
	cmpRes := &auxResource{rbac.ResourceComponent(r.Resource)}

	if !h.rbac.Can(rbac.ContextToSession(ctx), "grant", cmpRes) {
		return fmt.Errorf("not allowed to grant %s", r.String())
	}

	return h.rbac.Grant(ctx, r)
}

func (h rbacHandler) check(ctx context.Context, args *rbacCheckArgs) (b *rbacCheckResults, err error) {
	var secses rbac.Session

	b = &rbacCheckResults{}

	if args.hasUser {
		// user was explicitly set to rbac checking function
		//
		// load all roles that this user is member of
		// and create new security session out of it
		var (
			mm  types.RoleMemberSet
			ids []uint64
		)

		mm, err = h.role.Membership(ctx, args.User.ID)
		_ = mm.Walk(func(m *types.RoleMember) error {
			ids = append(ids, m.RoleID)
			return nil
		})

		args.User.SetRoles(ids...)

		secses = rbac.NewSession(ctx, args.User)
	} else {
		secses = rbac.ContextToSession(ctx)
	}

	b.Can = h.rbac.Can(secses, args.Operation, args.Resource)
	return
}
