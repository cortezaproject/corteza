package automation

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	rbacService interface {
		Can(ses rbac.Session, op string, res rbac.Resource) bool
		Grant(ctx context.Context, rules ...*rbac.Rule) (err error)
	}

	rbacUserService interface {
		FindByAny(ctx context.Context, identifier interface{}) (*types.User, error)
	}

	rbacRoleService interface {
		Membership(ctx context.Context, userID uint64) (types.RoleMemberSet, error)
	}

	rbacHandler struct {
		reg  rbacHandlerRegistry
		rbac rbacService
		user rbacUserService
		role rbacRoleService
	}

	auxResource struct{ res string }
)

var _ rbac.Resource = &auxResource{}

func (r auxResource) RbacResource() string {
	return r.res
}

func RbacHandler(reg rbacHandlerRegistry, rbac rbacService, user rbacUserService, role rbacRoleService) *rbacHandler {
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
	return h.grant(ctx, rbac.AllowRule(args.roleID, args.Operation, args.Resource.RbacResource()))
}

func (h rbacHandler) deny(ctx context.Context, args *rbacDenyArgs) (err error) {
	return h.grant(ctx, rbac.DenyRule(args.roleID, args.Operation, args.Resource.RbacResource()))
}

func (h rbacHandler) inherit(ctx context.Context, args *rbacInheritArgs) (err error) {
	return h.grant(ctx, rbac.InheritRule(args.roleID, args.Operation, args.Resource.RbacResource()))
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
